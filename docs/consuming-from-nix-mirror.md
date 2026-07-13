# Consuming the provider from the Nix mirror (registry-less)

Nivis consumes `terraform-provider-hcloudimage` from a Nix-built filesystem
mirror, with **no public registry** involved. This is the supported offline path
and is exercised by `scripts/verify-mirror.sh` and the hermetic VM test.

## Build the mirror

```sh
nix build .#provider-mirror
# -> ./result/registry.terraform.io/nivis-project/hcloudimage/<version>/<os>_<arch>/...
#    ./result/registry.opentofu.org/...      (so tofu resolves it too)
```

## Point the CLI at it (recommended: filesystem_mirror)

This recipe works for **both** `terraform` and `tofu`, including `init`:

```hcl
# ~/.terraformrc  (or point TF_CLI_CONFIG_FILE at this file)
provider_installation {
  filesystem_mirror {
    path    = "/absolute/path/to/result"
    include = [
      "registry.terraform.io/nivis-project/hcloudimage",
      "registry.opentofu.org/nivis-project/hcloudimage",
    ]
  }
  direct {
    # Everything else (e.g. hetznercloud/hcloud) still installs from its registry.
    exclude = [
      "registry.terraform.io/nivis-project/hcloudimage",
      "registry.opentofu.org/nivis-project/hcloudimage",
    ]
  }
}
```

Your consumer config pins the version normally:

```hcl
terraform {
  required_providers {
    hcloudimage = {
      source  = "nivis-project/hcloudimage"
      version = "0.1.0"
    }
  }
}
```

Then `terraform init && terraform plan` (or `tofu`) resolves the provider from
the mirror. Verify the whole path with:

```sh
make consume-mirror   # runs scripts/verify-mirror.sh under both binaries
```

## Alternative: dev_overrides (tight iteration, terraform)

For fast local iteration you can skip `init` entirely and point directly at a
freshly built binary:

```sh
nix build .#default   # -> ./result/bin/terraform-provider-hcloudimage
```

```hcl
provider_installation {
  dev_overrides {
    "nivis-project/hcloudimage" = "/absolute/path/to/result/bin"
  }
  direct {}
}
```

> **Note:** `dev_overrides` is honoured by `terraform` at `plan`/`apply` time and
> lets you skip `init`. `tofu init` still tries to resolve overridden providers
> from its registry, so for OpenTofu use the `filesystem_mirror` recipe above.
