terraform {
  required_providers {
    hcloudimage = {
      source  = "nivis-project/hcloudimage"
      version = "~> 0.1"
    }
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "~> 1.48"
    }
  }
}

provider "hcloudimage" {} # reads HCLOUD_TOKEN
provider "hcloud" {}

variable "image_path" {
  type    = string
  default = "result/nixos-hetzner.raw.xz"
}

# Upload a locally built NixOS image and snapshot it.
resource "hcloudimage_image" "nixos" {
  image_path   = var.image_path
  image_sha256 = filesha256(var.image_path) # ForceNew trigger for local files
  architecture = "x86"
  compression  = "xz"
  location     = "nbg1"

  labels = {
    os      = "nixos"
    creator = "nivis"
  }
}

# Boot a real server from the resulting snapshot.
resource "hcloud_server" "demo" {
  name        = "nixos-demo"
  image       = hcloudimage_image.nixos.id
  server_type = "cx22"
  location    = "nbg1"
}

output "snapshot_id" {
  value = hcloudimage_image.nixos.id
}
