# Hermetic NixOS-VM lifecycle test (BRIEFING.md §8.2) — the PoC Definition-of-Done gate.
#
# Boots a NixOS VM with terraform + opentofu and the built provider installed via a
# filesystem mirror, then drives init/plan/apply/destroy under both binaries with the
# fake uploader (HCLOUDIMAGE_FAKE=1). No cloud access, no network.
{
  pkgs,
  provider, # packages.default
  providerMirror, # filesystem-mirror derivation for the provider
}:
let
  version = "0.1.0";

  # CLI config selecting the mirror for our provider; direct{} is unused because
  # the test config references only hcloudimage.
  tfrc = pkgs.writeText "hermetic.tfrc" ''
    provider_installation {
      filesystem_mirror {
        path    = "${providerMirror}"
        include = ["registry.terraform.io/nivis-project/hcloudimage", "registry.opentofu.org/nivis-project/hcloudimage"]
      }
      direct {
        exclude = ["registry.terraform.io/nivis-project/hcloudimage", "registry.opentofu.org/nivis-project/hcloudimage"]
      }
    }
  '';

  # The test HCL, copied into the store.
  tfConfig = ../e2e;
in
pkgs.testers.runNixOSTest {
  name = "hcloudimage-hermetic-e2e";

  nodes.machine =
    { ... }:
    {
      environment.systemPackages = [
        pkgs.terraform
        pkgs.opentofu
        pkgs.jq
      ];
      # Give the VM enough memory for terraform + the provider process.
      virtualisation.memorySize = 2048;
    };

  testScript = ''
    import json

    machine.start()
    machine.wait_for_unit("multi-user.target")

    def run(cmd):
        return machine.succeed(cmd)

    def plan_actions(workdir, binary, extra_vars):
        """Return the list of resource_changes actions for the test resource."""
        run(
            f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
            f"{binary} plan {extra_vars} -out=plan.bin >/dev/null"
        )
        raw = run(
            f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
            f"{binary} show -json plan.bin"
        )
        data = json.loads(raw)
        for rc in data.get("resource_changes", []):
            if rc["address"] == "hcloudimage_image.test":
                return rc["change"]["actions"]
        return []

    def state_resource_count(workdir, binary):
        raw = run(
            f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
            f"{binary} show -json"
        )
        data = json.loads(raw)
        root = data.get("values", {}).get("root_module", {})
        return len(root.get("resources", []))

    for binary in ["terraform", "tofu"]:
        with subtest(f"lifecycle under {binary}"):
            workdir = f"/tmp/work-{binary}"
            run(f"mkdir -p {workdir}")
            run(f"cp ${tfConfig}/main.tf {workdir}/main.tf")
            # A small local file so image_sha256 is a realistic ForceNew trigger.
            run(f"printf 'disk-image-bytes' > {workdir}/image.raw")

            run(
                f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
                f"{binary} init -input=false >/dev/null"
            )

            # 1. Create (sha=A, env=a)
            varsA = '-var image_sha256=aaa -var env_label=a'
            run(
                f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
                f"{binary} apply -auto-approve {varsA} >/dev/null"
            )
            assert state_resource_count(workdir, binary) == 1, f"{binary}: expected 1 resource after create"
            first_id = run(
                f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
                f"{binary} output -raw snapshot_id"
            ).strip()
            assert first_id, f"{binary}: snapshot_id should be set after create"

            # 2. ForceNew: change sha -> expect replace (delete+create)
            varsB = '-var image_sha256=bbb -var env_label=a'
            actions = plan_actions(workdir, binary, varsB)
            assert actions == ["delete", "create"] or actions == ["create", "delete"], \
                f"{binary}: expected replace, got {actions}"
            run(
                f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
                f"{binary} apply -auto-approve {varsB} >/dev/null"
            )
            second_id = run(
                f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
                f"{binary} output -raw snapshot_id"
            ).strip()
            assert second_id != first_id, f"{binary}: replace should change the id"

            # 3. In-place: change only labels -> expect update
            varsC = '-var image_sha256=bbb -var env_label=b'
            actions = plan_actions(workdir, binary, varsC)
            assert actions == ["update"], f"{binary}: expected in-place update, got {actions}"
            run(
                f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
                f"{binary} apply -auto-approve {varsC} >/dev/null"
            )
            third_id = run(
                f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
                f"{binary} output -raw snapshot_id"
            ).strip()
            assert third_id == second_id, f"{binary}: in-place update must not change the id"

            # 4. Destroy -> empty state
            run(
                f"cd {workdir} && TF_CLI_CONFIG_FILE=${tfrc} HCLOUDIMAGE_FAKE=1 HCLOUDIMAGE_FAKE_STATE={workdir}/fake-state.json "
                f"{binary} destroy -auto-approve {varsC} >/dev/null"
            )
            assert state_resource_count(workdir, binary) == 0, f"{binary}: state should be empty after destroy"
  '';
}
