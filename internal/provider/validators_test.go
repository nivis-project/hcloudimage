package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestLabelValuesValidator(t *testing.T) {
	cases := map[string]struct {
		labels    map[string]string
		wantError bool
	}{
		"clean":       {map[string]string{"os": "nixos", "creator": "nivis"}, false},
		"slash value": {map[string]string{"path": "a/b"}, true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			elems := make(map[string]attr.Value, len(tc.labels))
			for k, v := range tc.labels {
				elems[k] = types.StringValue(v)
			}
			m, d := types.MapValue(types.StringType, elems)
			if d.HasError() {
				t.Fatalf("MapValue: %v", d)
			}
			var resp validator.MapResponse
			labelValuesValidator{}.ValidateMap(context.Background(), validator.MapRequest{
				Path:        path.Root("labels"),
				ConfigValue: m,
			}, &resp)

			if got := resp.Diagnostics.HasError(); got != tc.wantError {
				t.Errorf("HasError = %v, want %v (%v)", got, tc.wantError, resp.Diagnostics)
			}
		})
	}
}
