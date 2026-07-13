package provider

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// mustAbs turns a path from the environment into an absolute path so terraform's
// module-relative file functions resolve it correctly, and fails loudly if it
// doesn't exist. Note: `go test` runs with the working directory set to the
// package (internal/provider), so a relative value resolves against THAT, not the
// repo root — always pass an absolute path (e.g. "$PWD/$(ls result/*.raw.xz)").
func mustAbs(t *testing.T, name, p string) string {
	t.Helper()
	if p == "" {
		t.Fatalf("%s is empty", name)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		t.Fatalf("%s: resolving %q to an absolute path: %v", name, p, err)
	}
	if _, err := os.Stat(abs); err != nil {
		t.Fatalf("%s=%q does not exist (resolved to %q). Pass an ABSOLUTE path — "+
			"go test's working dir is the package, not the repo root. "+
			"Example: export %s=\"$PWD/$(ls result/*.raw.xz)\"", name, p, abs, name)
	}
	return abs
}

// Billable acceptance tests (BRIEFING.md §8.3). These run only when TF_ACC=1 and
// HCLOUD_TOKEN are set, against a real, isolated, budget-limited Hetzner project.
// They use the REAL uploader and the official hcloud provider, boot a server from
// the produced snapshot, and prove guest reachability by SSHing in with the baked
// throwaway key — not merely that the server reports "running".
//
// Local `go test ./...` skips these (no token), so the default suite needs no
// credentials.

// accPreCheck skips unless the acceptance environment is fully configured.
func accPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("HCLOUD_TOKEN") == "" {
		t.Skip("HCLOUD_TOKEN not set; skipping billable acceptance test")
	}
	for _, v := range []string{"HCLOUDIMAGE_ACC_IMAGE_PATH", "HCLOUDIMAGE_ACC_SSH_KEY"} {
		if os.Getenv(v) == "" {
			t.Skipf("%s not set; skipping acceptance (fixture path and throwaway private key are required)", v)
		}
	}
}

// accCase describes one architecture's acceptance run.
type accCase struct {
	name       string
	arch       string
	serverType string
}

func TestAccImage_RealHetzner_x86(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}
	accPreCheck(t)
	runAcceptance(t, accCase{name: "x86", arch: "x86", serverType: "cx22"})
}

func TestAccImage_RealHetzner_arm(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}
	if os.Getenv("HCLOUDIMAGE_ACC_RUN_ARM") != "1" {
		t.Skip("HCLOUDIMAGE_ACC_RUN_ARM != 1; arm acceptance is toggle-gated (cost)")
	}
	accPreCheck(t)
	runAcceptance(t, accCase{name: "arm", arch: "arm", serverType: "cax11"})
}

func runAcceptance(t *testing.T, c accCase) {
	t.Helper()

	// Resolve to absolute paths: the HCL uses filesha256(image_path) and file(key),
	// which terraform resolves relative to the module dir (a temp dir), not the
	// caller's cwd. A relative env value (e.g. "result/x.raw.xz") would otherwise
	// fail deep inside the plan with a confusing filesha256 error.
	imagePath := mustAbs(t, "HCLOUDIMAGE_ACC_IMAGE_PATH", os.Getenv("HCLOUDIMAGE_ACC_IMAGE_PATH"))
	sshKey := mustAbs(t, "HCLOUDIMAGE_ACC_SSH_KEY", os.Getenv("HCLOUDIMAGE_ACC_SSH_KEY"))

	// The real uploader is selected automatically because HCLOUD_TOKEN is set and
	// HCLOUDIMAGE_FAKE is unset.
	config := fmt.Sprintf(`
terraform {
  required_providers {
    # hcloudimage is injected in-process by ProtoV6ProviderFactories, so it must
    # NOT be listed here — declaring it makes terraform try to resolve it from the
    # registry during init and write an unsatisfiable lock entry.
    hcloud = { source = "hetznercloud/hcloud", version = "~> 1.48" }
  }
}

provider "hcloud" {}

resource "hcloudimage_image" "test" {
  image_path   = %q
  image_sha256 = filesha256(%q)
  architecture = %q
  compression  = "xz"
  location     = "nbg1"
  labels = { test = "hcloudimage-acc" }
}

resource "hcloud_ssh_key" "test" {
  name       = "hcloudimage-acc-%s"
  public_key = file(%q)
}

resource "hcloud_server" "test" {
  name        = "hcloudimage-acc-%s"
  image       = hcloudimage_image.test.id
  server_type = %q
  location    = "nbg1"
  ssh_keys    = [hcloud_ssh_key.test.id]
}

output "server_ipv4" {
  value = hcloud_server.test.ipv4_address
}
`, imagePath, imagePath, c.arch, c.name, sshKey+".pub", c.name, c.serverType)

	resource.Test(t, resource.TestCase{
		// The acceptance run pulls both providers from their registries via the
		// standard test harness; no ProtoV6ProviderFactories override here because
		// we need the real provider binary + the official hcloud provider.
		ExternalProviders: map[string]resource.ExternalProvider{
			"hcloud": {Source: "hetznercloud/hcloud", VersionConstraint: "~> 1.48"},
		},
		// Real provider: with HCLOUD_TOKEN set and HCLOUDIMAGE_FAKE unset, New()
		// selects the hcloudimages/v2-backed uploader.
		ProtoV6ProviderFactories: protoV6RealFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hcloudimage_image.test", "id"),
					resource.TestCheckResourceAttr("hcloudimage_image.test", "architecture", c.arch),
					checkGuestReachable("hcloud_server.test", sshKey),
				),
			},
		},
	})
}

// checkGuestReachable SSHes from the runner into the booted server with the
// throwaway key and reads /etc/os-release, proving the rootfs actually booted —
// "running" from the hypervisor is not sufficient (BRIEFING.md §8.3).
func checkGuestReachable(resourceName, privateKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not in state", resourceName)
		}
		ip := rs.Primary.Attributes["ipv4_address"]
		if ip == "" {
			return fmt.Errorf("server has no ipv4_address")
		}

		// Retry: the guest needs a moment to bring up sshd after boot.
		var lastErr error
		for attempt := 0; attempt < 30; attempt++ {
			out, err := sshCommand(privateKey, ip, "cat /etc/os-release")
			if err == nil {
				if !strings.Contains(strings.ToLower(out), "alpine") &&
					!strings.Contains(strings.ToLower(out), "nixos") {
					return fmt.Errorf("unexpected /etc/os-release from guest: %s", out)
				}
				return nil
			}
			lastErr = err
			time.Sleep(10 * time.Second)
		}
		return fmt.Errorf("guest never became reachable over SSH: %w", lastErr)
	}
}

func sshCommand(privateKey, host, cmd string) (string, error) {
	args := []string{
		"-i", privateKey,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=10",
		"-o", "BatchMode=yes",
		fmt.Sprintf("root@%s", host),
		cmd,
	}
	out, err := exec.Command("ssh", args...).CombinedOutput()
	return string(out), err
}
