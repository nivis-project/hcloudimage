package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                     = (*snapshotDataSource)(nil)
	_ datasource.DataSourceWithConfigure        = (*snapshotDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*snapshotDataSource)(nil)
)

// snapshotDataSourceModel maps the hcloudimage_snapshot schema (BRIEFING.md §3.3).
type snapshotDataSourceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	WithSelector types.String `tfsdk:"with_selector"`
	MostRecent   types.Bool   `tfsdk:"most_recent"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Architecture types.String `tfsdk:"architecture"`
	Created      types.String `tfsdk:"created"`
	Labels       types.Map    `tfsdk:"labels"`
}

type snapshotDataSource struct {
	uploader Uploader
}

// NewSnapshotDataSource is the data-source factory registered on the provider.
func NewSnapshotDataSource() datasource.DataSource {
	return &snapshotDataSource{}
}

func (d *snapshotDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

func (d *snapshotDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(providerData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("expected providerData, got %T", req.ProviderData))
		return
	}
	d.uploader = data.Uploader
}

func (d *snapshotDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Looks up an existing Hetzner Cloud snapshot by ID or label selector, so you can reference images not created in this state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Look up by image ID. Mutually exclusive with `with_selector`.",
				Optional:            true,
				Computed:            true,
			},
			"with_selector": schema.StringAttribute{
				MarkdownDescription: "Hetzner label selector (`key=value`). Must resolve to exactly one snapshot unless `most_recent` is set.",
				Optional:            true,
			},
			"most_recent": schema.BoolAttribute{
				MarkdownDescription: "When a selector matches multiple snapshots, pick the newest instead of erroring.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Snapshot name.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Snapshot description.",
				Computed:            true,
			},
			"architecture": schema.StringAttribute{
				MarkdownDescription: "Guest architecture (`x86` or `arm`).",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "Creation timestamp.",
				Computed:            true,
			},
			"labels": schema.MapAttribute{
				MarkdownDescription: "Snapshot labels.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *snapshotDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("with_selector"),
		),
	}
}

func (d *snapshotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config snapshotDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	info, err := d.uploader.Find(ctx, config.ID.ValueInt64(), config.WithSelector.ValueString(), config.MostRecent.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Snapshot lookup failed", err.Error())
		return
	}

	config.ID = types.Int64Value(info.ID)
	config.Name = types.StringValue(info.Name)
	config.Description = types.StringValue(info.Description)
	config.Architecture = types.StringValue(info.Architecture)
	config.Created = types.StringValue(info.Created)

	labels, diags := types.MapValueFrom(ctx, types.StringType, info.Labels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Labels = labels

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
