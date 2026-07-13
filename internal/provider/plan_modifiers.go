package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// effectiveLabelsModifier computes the planned effective_labels from the planned
// user labels merged with the library's created-by default. Because every input
// is known at plan time, this yields an exact plan (no "known after apply") and
// stays consistent with what Create/Update actually write — avoiding both
// spurious diffs and "inconsistent result after apply" errors.
type effectiveLabelsModifier struct{}

func effectiveLabelsPlan() planmodifier.Map { return effectiveLabelsModifier{} }

func (m effectiveLabelsModifier) Description(_ context.Context) string {
	return "computes effective_labels as user labels merged with the library created-by label"
}

func (m effectiveLabelsModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m effectiveLabelsModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	// Derive effective_labels from the planned user labels attribute.
	var userLabels types.Map
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, req.Path.ParentPath().AtName("labels"), &userLabels)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if userLabels.IsUnknown() {
		resp.PlanValue = types.MapUnknown(types.StringType)
		return
	}

	raw := map[string]string{}
	if !userLabels.IsNull() {
		resp.Diagnostics.Append(userLabels.ElementsAs(ctx, &raw, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	value, d := types.MapValueFrom(ctx, types.StringType, mergeEffectiveLabels(raw))
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.PlanValue = value
}
