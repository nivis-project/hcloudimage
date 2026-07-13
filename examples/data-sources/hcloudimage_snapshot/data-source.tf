terraform {
  required_providers {
    hcloudimage = {
      source  = "nivis-project/hcloudimage"
      version = "~> 0.1"
    }
  }
}

provider "hcloudimage" {}

# Look up a snapshot by label selector, choosing the newest match.
data "hcloudimage_snapshot" "base" {
  with_selector = "os=nixos"
  most_recent   = true
}

# Or look up a specific snapshot by ID.
data "hcloudimage_snapshot" "by_id" {
  id = 12345678
}

output "base_snapshot_architecture" {
  value = data.hcloudimage_snapshot.base.architecture
}
