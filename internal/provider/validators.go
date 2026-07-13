package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// sha256RequiredWithPathValidator enforces that image_sha256 is set if and only
// if image_path is set (BRIEFING.md §3.2). image_url uses the URL string itself
// as the ForceNew trigger, so it needs no sha.
type sha256RequiredWithPathValidator struct{}

func (v sha256RequiredWithPathValidator) Description(_ context.Context) string {
	return "image_sha256 must be set when image_path is used, and must be omitted otherwise"
}

func (v sha256RequiredWithPathValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v sha256RequiredWithPathValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data imageResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathSet := !data.ImagePath.IsNull() && !data.ImagePath.IsUnknown()
	shaSet := !data.ImageSHA256.IsNull() && !data.ImageSHA256.IsUnknown()

	if pathSet && !shaSet {
		resp.Diagnostics.AddAttributeError(
			path.Root("image_sha256"),
			"Missing image_sha256",
			"image_sha256 is required when image_path is set. Use filesha256(var.image_path).",
		)
	}
	if !pathSet && shaSet {
		resp.Diagnostics.AddAttributeError(
			path.Root("image_sha256"),
			"Unexpected image_sha256",
			"image_sha256 only applies to image_path; omit it when using image_url.",
		)
	}
}

// labelValuesValidator rejects label values containing '/', which Hetzner
// disallows (BRIEFING.md §3.2).
type labelValuesValidator struct{}

func (v labelValuesValidator) Description(_ context.Context) string {
	return "label values must not contain '/'"
}

func (v labelValuesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v labelValuesValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var labels map[string]string
	resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &labels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	for k, val := range labels {
		if strings.Contains(val, "/") {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid label value",
				fmt.Sprintf("label %q has value %q; Hetzner label values must not contain '/'.", k, val),
			)
		}
	}
}

// ensure interface satisfaction
var (
	_ resource.ConfigValidator = sha256RequiredWithPathValidator{}
	_ validator.Map            = labelValuesValidator{}
)
