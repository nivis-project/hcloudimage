package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/nivis-project/terraform-provider-hcloudimage/internal/provider"
)

// These are set by goreleaser at build time via ldflags.
var (
	version = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// Registry address used by Terraform/OpenTofu to source the provider.
		Address: "registry.terraform.io/nivis-project/hcloudimage",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
