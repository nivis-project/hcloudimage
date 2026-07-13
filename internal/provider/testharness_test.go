package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// protoV6Factories wires the provider (backed by the given fake) into the
// terraform-plugin-testing harness, so a test can drive Terraform and inspect
// the same fake instance the provider uses.
func protoV6Factories(fake Uploader) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"hcloudimage": providerserver.NewProtocol6WithError(NewWithUploader("test", fake)()),
	}
}

// protoV6RealFactories wires the provider with its normal uploader selection
// (real when HCLOUD_TOKEN is set). Used by the billable acceptance tests.
func protoV6RealFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"hcloudimage": providerserver.NewProtocol6WithError(New("acc")()),
	}
}
