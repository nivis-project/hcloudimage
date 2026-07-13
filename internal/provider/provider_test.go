package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProvider_Metadata(t *testing.T) {
	p := New("test1.2.3")()

	var resp provider.MetadataResponse
	p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)

	if resp.TypeName != "hcloudimage" {
		t.Errorf("TypeName = %q, want %q", resp.TypeName, "hcloudimage")
	}
	if resp.Version != "test1.2.3" {
		t.Errorf("Version = %q, want %q", resp.Version, "test1.2.3")
	}
}

func TestProvider_Schema(t *testing.T) {
	p := New("test")()

	var resp provider.SchemaResponse
	p.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	for _, name := range []string{"token", "endpoint", "poll_interval"} {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("provider schema missing attribute %q", name)
			continue
		}
		if !attr.IsOptional() {
			t.Errorf("attribute %q should be optional", name)
		}
	}

	if token := resp.Schema.Attributes["token"]; token != nil && !token.IsSensitive() {
		t.Errorf("attribute %q should be sensitive", "token")
	}
}

func TestProvider_Resources(t *testing.T) {
	prov := New("test")()
	resources := prov.Resources(context.Background())
	if len(resources) != 1 {
		t.Fatalf("Resources() returned %d resources, want 1", len(resources))
	}
}

func TestProvider_DataSources(t *testing.T) {
	prov := New("test")()
	dataSources := prov.DataSources(context.Background())
	if len(dataSources) != 1 {
		t.Fatalf("DataSources() returned %d data sources, want 1", len(dataSources))
	}
}
