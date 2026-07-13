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

	// newUploader builds the Uploader from resolved config. Defaults to the fake
	// (milestone 02); milestone 04 swaps in the real hcloudimages/v2-backed one.
	// Tests override this to inject a specific fake instance.
	newUploader func(providerConfig) (Uploader, error)
}

// providerModel maps the provider configuration schema (BRIEFING.md §3.1) to Go values.
type providerModel struct {
	Token        types.String `tfsdk:"token"`
	Endpoint     types.String `tfsdk:"endpoint"`
	PollInterval types.String `tfsdk:"poll_interval"`
}

// providerConfig is the resolved provider configuration.
type providerConfig struct {
	Token        string
	Endpoint     string
	PollInterval string
}

// providerData is handed to resources and data sources via configure data. It
// carries the resolved config and the Uploader implementation to use. Tests can
// inject a fake uploader by overriding newUploader.
type providerData struct {
	Config   providerConfig
	Uploader Uploader
}

// New returns a factory for the provider, capturing the build version. The
// default uploader is the in-memory fake (milestone 02).
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hcloudimageProvider{
			version: version,
			newUploader: func(providerConfig) (Uploader, error) {
				return NewFakeUploader(), nil
			},
		}
	}
}

// NewWithUploader is a test seam: it builds the provider with a fixed Uploader,
// so lifecycle tests can inspect the same fake the resource uses.
func NewWithUploader(version string, uploader Uploader) func() provider.Provider {
	return func() provider.Provider {
		return &hcloudimageProvider{
			version:     version,
			newUploader: func(providerConfig) (Uploader, error) { return uploader, nil },
		}
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

	uploader, err := p.newUploader(cfg)
	if err != nil {
		resp.Diagnostics.AddError("Failed to construct uploader", err.Error())
		return
	}

	data := providerData{Config: cfg, Uploader: uploader}

	// Resources and data sources read config + uploader from here.
	resp.ResourceData = data
	resp.DataSourceData = data
}

func (p *hcloudimageProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewImageResource,
	}
}

func (p *hcloudimageProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSnapshotDataSource,
	}
}
