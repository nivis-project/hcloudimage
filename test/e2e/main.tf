terraform {
  required_providers {
    hcloudimage = {
      source  = "nivis-project/hcloudimage"
      version = "0.1.0"
    }
  }
}

provider "hcloudimage" {
  # The hermetic test runs with HCLOUDIMAGE_FAKE=1, so no token is needed.
}

variable "image_sha256" {
  type = string
}

variable "env_label" {
  type = string
}

# Local-path source so image_sha256 is the ForceNew trigger (BRIEFING.md §8.2).
resource "hcloudimage_image" "test" {
  image_path   = "${path.module}/image.raw"
  image_sha256 = var.image_sha256
  architecture = "x86"
  compression  = "xz"

  labels = {
    env = var.env_label
  }
}

output "snapshot_id" {
  value = hcloudimage_image.test.id
}
