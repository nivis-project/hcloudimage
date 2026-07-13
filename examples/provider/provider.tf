terraform {
  required_providers {
    hcloudimage = {
      source  = "nivis-project/hcloudimage"
      version = "~> 0.1"
    }
  }
}

# Reads HCLOUD_TOKEN from the environment when token is unset.
provider "hcloudimage" {}
