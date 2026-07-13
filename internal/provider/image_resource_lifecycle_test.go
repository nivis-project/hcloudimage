package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// configURL builds a minimal hcloudimage_image config from a URL source.
func configURL(url, description string, labels map[string]string) string {
	labelHCL := ""
	for k, v := range labels {
		labelHCL += fmt.Sprintf("    %s = %q\n", k, v)
	}
	return fmt.Sprintf(`
resource "hcloudimage_image" "test" {
  image_url    = %q
  architecture = "x86"
  compression  = "xz"
  description  = %q
  labels = {
%s  }
}
`, url, description, labelHCL)
}

// configPath builds a config from a local path source (requires image_sha256).
func configPath(path, sha string) string {
	return fmt.Sprintf(`
resource "hcloudimage_image" "test" {
  image_path   = %q
  image_sha256 = %q
  architecture = "x86"
}
`, path, sha)
}

func TestAccImageResource_CreateAndDestroy(t *testing.T) {
	fake := NewFakeUploader()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories(fake),
		CheckDestroy: func(*terraform.State) error {
			if len(fake.DeleteCalls) == 0 {
				return fmt.Errorf("destroy did not call Delete")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: configURL("https://example.com/img.raw.xz", "first", map[string]string{"os": "nixos"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hcloudimage_image.test", "id"),
					resource.TestCheckResourceAttr("hcloudimage_image.test", "effective_labels.%", "2"),
					resource.TestCheckResourceAttr("hcloudimage_image.test", "effective_labels.os", "nixos"),
					checkEffectiveLabel("hcloudimage_image.test", CreatedByLabelKey, CreatedByLabelValue),
					func(*terraform.State) error {
						if len(fake.UploadCalls) != 1 {
							return fmt.Errorf("expected 1 upload, got %d", len(fake.UploadCalls))
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccImageResource_InPlaceUpdate(t *testing.T) {
	fake := NewFakeUploader()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories(fake),
		Steps: []resource.TestStep{
			{
				Config: configURL("https://example.com/img.raw.xz", "first", map[string]string{"env": "a"}),
			},
			{
				Config: configURL("https://example.com/img.raw.xz", "second", map[string]string{"env": "b"}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hcloudimage_image.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hcloudimage_image.test", "description", "second"),
					resource.TestCheckResourceAttr("hcloudimage_image.test", "labels.env", "b"),
					func(*terraform.State) error {
						if len(fake.UploadCalls) != 1 {
							return fmt.Errorf("in-place update should not re-upload; uploads = %d", len(fake.UploadCalls))
						}
						if len(fake.UpdateCalls) == 0 {
							return fmt.Errorf("in-place update did not call UpdateMetadata")
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccImageResource_ForceNewReplace(t *testing.T) {
	fake := NewFakeUploader()

	dir := t.TempDir()
	imgPath := dir + "/img.raw"
	if err := os.WriteFile(imgPath, []byte("disk image bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories(fake),
		Steps: []resource.TestStep{
			{
				Config: configPath(imgPath, "sha-aaa"),
			},
			{
				Config: configPath(imgPath, "sha-bbb"), // ForceNew trigger
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hcloudimage_image.test", plancheck.ResourceActionReplace),
					},
				},
				Check: func(*terraform.State) error {
					if len(fake.UploadCalls) != 2 {
						return fmt.Errorf("ForceNew should re-upload; uploads = %d, want 2", len(fake.UploadCalls))
					}
					if len(fake.DeleteCalls) == 0 {
						return fmt.Errorf("ForceNew replace should delete the old snapshot")
					}
					return nil
				},
			},
		},
	})
}

func TestAccImageResource_OutOfBandDeletion(t *testing.T) {
	fake := NewFakeUploader()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories(fake),
		Steps: []resource.TestStep{
			{
				Config: configURL("https://example.com/img.raw.xz", "first", map[string]string{"k": "v"}),
			},
			{
				// Delete the snapshot out of band, then refresh: Read finds the
				// snapshot gone and removes the resource from state, so the refresh
				// plan is non-empty (Terraform would recreate it on the next apply).
				PreConfig: func() {
					for _, id := range lastUploadedIDs(fake) {
						fake.MarkDeleted(id)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				Check: func(s *terraform.State) error {
					if _, ok := s.RootModule().Resources["hcloudimage_image.test"]; ok {
						return fmt.Errorf("resource should have been removed from state after out-of-band deletion")
					}
					return nil
				},
			},
		},
	})
}

// lastUploadedIDs returns the ids currently live in the fake.
func lastUploadedIDs(f *FakeUploader) []int64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	var ids []int64
	for id, snap := range f.store {
		if !snap.deleted {
			ids = append(ids, id)
		}
	}
	return ids
}

// checkEffectiveLabel asserts a single effective_labels entry by inspecting the
// state directly, avoiding the dot-delimited attribute-path escaping problem for
// keys that themselves contain '.' and '/' (e.g. apricote.de/created-by).
func checkEffectiveLabel(resourceName, key, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		got := rs.Primary.Attributes["effective_labels."+key]
		if got != want {
			return fmt.Errorf("effective_labels[%q] = %q, want %q", key, got, want)
		}
		return nil
	}
}
