package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the resource satisfies the framework interface.
var _ resource.Resource = (*imageResource)(nil)

// imageResource manages an uploaded Hetzner Cloud snapshot (hcloudimage_image).
//
// Scaffold placeholder: the schema and CRUD behaviour of BRIEFING.md §3.2 are added in
// milestone 02. This registers the resource so the provider server starts.
type imageResource struct{}

// NewImageResource is the resource factory registered on the provider.
func NewImageResource() resource.Resource {
	return &imageResource{}
}

func (r *imageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

func (r *imageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Uploads a raw disk image to Hetzner Cloud and snapshots it. (Schema completed in milestone 02.)",
		Attributes:          map[string]schema.Attribute{},
	}
}

func (r *imageResource) Create(_ context.Context, _ resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError("Not implemented", "hcloudimage_image is not yet implemented (milestone 02).")
}

func (r *imageResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *imageResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Not implemented", "hcloudimage_image is not yet implemented (milestone 02).")
}

func (r *imageResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
