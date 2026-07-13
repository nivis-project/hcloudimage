package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the provider satisfies the framework interface.
var _ provider.Provider = (*hcloudimageProvider)(nil)

// hcloudimageProvider is the provider implementation.
type hcloudimageProvider struct {
	// version is set at build time and surfaced to Terraform via Metadata.
	version string
}

// providerModel maps the provider configuration schema (BRIEFING.md §3.1) to Go values.
type providerModel struct {
	Token        types.String `tfsdk:"token"`
	Endpoint     types.String `tfsdk:"endpoint"`
	PollInterval types.String `tfsdk:"poll_interval"`
}

// providerConfig is the resolved configuration handed to resources and data sources.
//
// The scaffold plumbs the configuration through only; the real hcloud client is wired
// behind the Uploader interface in a later milestone (BRIEFING.md §4).
type providerConfig struct {
	Token        string
	Endpoint     string
	PollInterval string
}

// New returns a factory for the provider, capturing the build version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hcloudimageProvider{version: version}
	}
}

func (p *hcloudimageProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hcloudimage"
	resp.Version = p.version
}

func (p *hcloudimageProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Uploads a raw disk image into a Hetzner Cloud project and turns it into a reusable snapshot, using the rescue-server upload trick.",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "Hetzner Cloud API token. Falls back to the `HCLOUD_TOKEN` environment variable when unset.",
				Optional:            true,
				Sensitive:           true,
			},
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Override the hcloud API endpoint (for testing or a mock). Defaults to the SDK default.",
				Optional:            true,
			},
			"poll_interval": schema.StringAttribute{
				MarkdownDescription: "Optional passthrough for action polling, expressed as a Go duration string (e.g. `500ms`, `2s`).",
				Optional:            true,
			},
		},
	}
}

func (p *hcloudimageProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// token falls back to HCLOUD_TOKEN when the attribute is unset (BRIEFING.md §3.1).
	token := config.Token.ValueString()
	if config.Token.IsNull() || config.Token.IsUnknown() {
		token = os.Getenv("HCLOUD_TOKEN")
	}

	cfg := providerConfig{
		Token:        token,
		Endpoint:     config.Endpoint.ValueString(),
		PollInterval: config.PollInterval.ValueString(),
	}

	// Resources and data sources read the resolved configuration from here.
	resp.ResourceData = cfg
	resp.DataSourceData = cfg
}

func (p *hcloudimageProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewImageResource,
	}
}

func (p *hcloudimageProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
