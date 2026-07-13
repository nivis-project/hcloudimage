package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = (*imageResource)(nil)
	_ resource.ResourceWithConfigure        = (*imageResource)(nil)
	_ resource.ResourceWithConfigValidators = (*imageResource)(nil)
	_ resource.ResourceWithValidateConfig   = (*imageResource)(nil)
)

// imageResourceModel maps the hcloudimage_image schema (BRIEFING.md §3.2) to Go.
type imageResourceModel struct {
	ImageURL        types.String   `tfsdk:"image_url"`
	ImagePath       types.String   `tfsdk:"image_path"`
	ImageSHA256     types.String   `tfsdk:"image_sha256"`
	Architecture    types.String   `tfsdk:"architecture"`
	Compression     types.String   `tfsdk:"compression"`
	Format          types.String   `tfsdk:"format"`
	ServerType      types.String   `tfsdk:"server_type"`
	Location        types.String   `tfsdk:"location"`
	ImageSize       types.Int64    `tfsdk:"image_size"`
	Description     types.String   `tfsdk:"description"`
	Labels          types.Map      `tfsdk:"labels"`
	ID              types.Int64    `tfsdk:"id"`
	EffectiveLabels types.Map      `tfsdk:"effective_labels"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}

type imageResource struct {
	uploader Uploader
}

// NewImageResource is the resource factory registered on the provider.
func NewImageResource() resource.Resource {
	return &imageResource{}
}

func (r *imageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

func (r *imageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(providerData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("expected providerData, got %T", req.ProviderData))
		return
	}
	r.uploader = data.Uploader
}

func (r *imageResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	forceNewString := []planmodifier.String{stringplanmodifier.RequiresReplace()}
	// Computed defaults (compression/format/location) must keep their prior
	// state value when the config omits them, or every refresh would show them
	// going to "known after apply" and spuriously force replacement.
	forceNewComputedString := []planmodifier.String{
		stringplanmodifier.UseStateForUnknown(),
		stringplanmodifier.RequiresReplace(),
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Uploads a raw disk image into a Hetzner Cloud project and snapshots it via the rescue-server upload trick.\n\n" +
			"Provide the image with either `image_url` (the rescue server pulls it — fast, off your uplink) or `image_path` " +
			"(streamed from the apply host over SSH — bounded by your upload bandwidth). Use `filesha256(var.image_path)` for " +
			"`image_sha256` so local-file changes trigger a new snapshot.",
		Attributes: map[string]schema.Attribute{
			"image_url": schema.StringAttribute{
				MarkdownDescription: "Public `https://` URL of the image; the rescue server pulls it directly (fast, off your uplink). Mutually exclusive with `image_path`. Changing it forces a new snapshot.",
				Optional:            true,
				PlanModifiers:       forceNewString,
			},
			"image_path": schema.StringAttribute{
				MarkdownDescription: "Local file path on the apply host; streamed over SSH (bounded by your upload bandwidth). Mutually exclusive with `image_url`. Requires `image_sha256`. Changing it forces a new snapshot.",
				Optional:            true,
				PlanModifiers:       forceNewString,
			},
			"image_sha256": schema.StringAttribute{
				MarkdownDescription: "SHA-256 of the local image file, required when `image_path` is set. This is the ForceNew trigger for local files — set it with `filesha256(var.image_path)`.",
				Optional:            true,
				PlanModifiers:       forceNewString,
			},
			"architecture": schema.StringAttribute{
				MarkdownDescription: "Guest architecture: `x86` or `arm`. Changing it forces a new snapshot.",
				Required:            true,
				PlanModifiers:       forceNewString,
				Validators:          []validator.String{stringvalidator.OneOf(ValidArchitectures...)},
			},
			"compression": schema.StringAttribute{
				MarkdownDescription: "Image compression: `none` (default), `bz2`, `xz`, or `zstd`. Changing it forces a new snapshot.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       forceNewComputedString,
				Validators:          []validator.String{stringvalidator.OneOf(ValidCompressions...)},
			},
			"format": schema.StringAttribute{
				MarkdownDescription: "Image format: `raw` (default) or `qcow2`. Changing it forces a new snapshot.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       forceNewComputedString,
				Validators:          []validator.String{stringvalidator.OneOf(ValidFormats...)},
			},
			"server_type": schema.StringAttribute{
				MarkdownDescription: "Override the temporary rescue server type. Defaults per architecture (`x86 → cx22`, `arm → cax11`). Changing it forces a new snapshot.",
				Optional:            true,
				PlanModifiers:       forceNewString,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Temporary server location. Defaults to `fsn1`. Changing it forces a new snapshot.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       forceNewComputedString,
			},
			"image_size": schema.Int64Attribute{
				MarkdownDescription: "Optional pre-write size validation, passed through to the upload library. Changing it forces a new snapshot.",
				Optional:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Snapshot description. Updated in place without a re-upload. If unset, " +
					"Hetzner assigns a default (`snapshot <timestamp>`), which is reflected here.",
				Optional: true,
				// Computed too: the hcloud API auto-generates a description when the
				// user omits one, so Terraform must accept that server value without
				// treating it as drift.
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"labels": schema.MapAttribute{
				MarkdownDescription: "User labels merged onto the library defaults; managed in place without a re-upload. Values must not contain `/` (Hetzner rule).",
				ElementType:         types.StringType,
				Optional:            true,
				Validators:          []validator.Map{labelValuesValidator{}},
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Snapshot image ID.",
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"effective_labels": schema.MapAttribute{
				MarkdownDescription: "Final label set on the snapshot (user labels + library defaults such as `apricote.de/created-by`).",
				ElementType:         types.StringType,
				Computed:            true,
				PlanModifiers:       []planmodifier.Map{effectiveLabelsPlan()},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
		},
	}
}

func (r *imageResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("image_url"),
			path.MatchRoot("image_path"),
		),
		sha256RequiredWithPathValidator{},
	}
}

func (r *imageResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	sha256RequiredWithPathValidator{}.ValidateResource(ctx, req, resp)
}

// toUploadRequest maps the resource model to the library-agnostic UploadRequest
// (BRIEFING.md §4.2). Applies the architecture-based server-type and location
// defaults from the Appendix.
func (m imageResourceModel) toUploadRequest(ctx context.Context) (UploadRequest, error) {
	labels := map[string]string{}
	if !m.Labels.IsNull() && !m.Labels.IsUnknown() {
		diags := m.Labels.ElementsAs(ctx, &labels, false)
		if diags.HasError() {
			return UploadRequest{}, fmt.Errorf("invalid labels")
		}
	}

	compression := m.Compression.ValueString()
	if m.Compression.IsNull() || m.Compression.IsUnknown() {
		compression = "none"
	}
	format := m.Format.ValueString()
	if m.Format.IsNull() || m.Format.IsUnknown() {
		format = "raw"
	}
	location := m.Location.ValueString()
	if m.Location.IsNull() || m.Location.IsUnknown() {
		location = "fsn1"
	}

	return UploadRequest{
		ImageURL:     m.ImageURL.ValueString(),
		ImagePath:    m.ImagePath.ValueString(),
		ImageSHA256:  m.ImageSHA256.ValueString(),
		Architecture: m.Architecture.ValueString(),
		Compression:  compression,
		Format:       format,
		ServerType:   m.ServerType.ValueString(),
		Location:     location,
		ImageSize:    m.ImageSize.ValueInt64(),
		Description:  m.Description.ValueString(),
		Labels:       labels,
	}, nil
}

func (r *imageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan imageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uploadReq, err := plan.toUploadRequest(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	id, effective, err := r.uploader.Upload(ctx, uploadReq)
	if err != nil {
		resp.Diagnostics.AddError("Upload failed", err.Error())
		return
	}

	plan.ID = types.Int64Value(id)
	plan.Compression = types.StringValue(uploadReq.Compression)
	plan.Format = types.StringValue(uploadReq.Format)
	plan.Location = types.StringValue(uploadReq.Location)
	resp.Diagnostics.Append(setEffectiveLabels(ctx, &plan, effective)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// description is Optional+Computed: when the config omits it, Hetzner assigns
	// a default (e.g. "snapshot <timestamp>"). Read it back so state matches the
	// server and the next plan is clean.
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		info, err := r.uploader.Get(ctx, id)
		if err == nil && info != nil {
			plan.Description = types.StringValue(info.Description)
		} else {
			// Fall back to empty-known so the value is never left null/unknown.
			plan.Description = types.StringValue("")
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *imageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state imageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	info, err := r.uploader.Get(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Read failed", err.Error())
		return
	}
	if info == nil {
		// Snapshot no longer exists — drop from state without error (§3.2).
		resp.State.RemoveResource(ctx)
		return
	}

	// description is Computed, so reflect the server's value directly (the server
	// assigns a default when the user omits one).
	state.Description = types.StringValue(info.Description)
	resp.Diagnostics.Append(setEffectiveLabels(ctx, &state, info.Labels)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *imageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state imageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only description/labels reach Update; everything else is ForceNew.
	labels := map[string]string{}
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	effective, err := r.uploader.UpdateMetadata(ctx, state.ID.ValueInt64(), plan.Description.ValueString(), labels)
	if err != nil {
		resp.Diagnostics.AddError("In-place update failed", err.Error())
		return
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(setEffectiveLabels(ctx, &plan, effective)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *imageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state imageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.uploader.Delete(ctx, state.ID.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}

// setEffectiveLabels writes the merged label map into the model, returning any
// conversion diagnostics.
func setEffectiveLabels(ctx context.Context, m *imageResourceModel, labels map[string]string) diag.Diagnostics {
	value, d := types.MapValueFrom(ctx, types.StringType, labels)
	if d.HasError() {
		return d
	}
	m.EffectiveLabels = value
	return d
}
