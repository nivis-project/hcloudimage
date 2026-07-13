package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	// Read only the two attributes we care about. Decoding the whole model can
	// fail on unknown map values during plan, and we don't need the rest here.
	var imagePath, imageSHA256 types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("image_path"), &imagePath)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("image_sha256"), &imageSHA256)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// An unknown value (e.g. filesha256(var.x) or var.x) is resolved at apply
	// time; we can't judge presence yet, so don't flag it. Only a definitively
	// null attribute counts as "absent".
	if imagePath.IsUnknown() || imageSHA256.IsUnknown() {
		return
	}

	pathSet := !imagePath.IsNull()
	shaSet := !imageSHA256.IsNull()

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
	// Iterate elements directly: a map can be known while individual values are
	// unknown (e.g. `labels = { env = var.x }` during the initial plan walk).
	// Decoding such a map into map[string]string errors, so skip unknown/null
	// element values instead.
	for k, elem := range req.ConfigValue.Elements() {
		strVal, ok := elem.(types.String)
		if !ok || strVal.IsUnknown() || strVal.IsNull() {
			continue
		}
		if strings.Contains(strVal.ValueString(), "/") {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid label value",
				fmt.Sprintf("label %q has value %q; Hetzner label values must not contain '/'.", k, strVal.ValueString()),
			)
		}
	}
}

// ensure interface satisfaction
var (
	_ resource.ConfigValidator = sha256RequiredWithPathValidator{}
	_ validator.Map            = labelValuesValidator{}
)
