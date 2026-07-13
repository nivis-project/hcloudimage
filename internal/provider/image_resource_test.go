package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newImageModel(mut func(*imageResourceModel)) imageResourceModel {
	labels, _ := types.MapValue(types.StringType, map[string]attr.Value{})
	m := imageResourceModel{
		ImageURL:     types.StringNull(),
		ImagePath:    types.StringNull(),
		ImageSHA256:  types.StringNull(),
		Architecture: types.StringValue("x86"),
		Compression:  types.StringNull(),
		Format:       types.StringNull(),
		ServerType:   types.StringNull(),
		Location:     types.StringNull(),
		ImageSize:    types.Int64Null(),
		Description:  types.StringNull(),
		Labels:       labels,
	}
	if mut != nil {
		mut(&m)
	}
	return m
}

func TestToUploadRequest_Defaults(t *testing.T) {
	m := newImageModel(func(m *imageResourceModel) {
		m.ImageURL = types.StringValue("https://example.com/img.raw")
	})
	req, err := m.toUploadRequest(context.Background())
	if err != nil {
		t.Fatalf("toUploadRequest: %v", err)
	}
	if req.Compression != "none" {
		t.Errorf("default compression = %q, want none", req.Compression)
	}
	if req.Format != "raw" {
		t.Errorf("default format = %q, want raw", req.Format)
	}
	if req.Location != "fsn1" {
		t.Errorf("default location = %q, want fsn1", req.Location)
	}
}

func TestToUploadRequest_AllCompressions(t *testing.T) {
	for _, c := range ValidCompressions {
		m := newImageModel(func(m *imageResourceModel) {
			m.ImageURL = types.StringValue("https://example.com/img.raw")
			m.Compression = types.StringValue(c)
		})
		req, err := m.toUploadRequest(context.Background())
		if err != nil {
			t.Fatalf("compression %q: %v", c, err)
		}
		if req.Compression != c {
			t.Errorf("compression = %q, want %q", req.Compression, c)
		}
	}
}

func TestToUploadRequest_BothSourcesAndArches(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*imageResourceModel)
		want func(UploadRequest) bool
	}{
		{
			name: "url source arm",
			mut: func(m *imageResourceModel) {
				m.ImageURL = types.StringValue("https://example.com/a.raw")
				m.Architecture = types.StringValue("arm")
			},
			want: func(r UploadRequest) bool { return r.ImageURL != "" && r.Architecture == "arm" },
		},
		{
			name: "path source x86",
			mut: func(m *imageResourceModel) {
				m.ImagePath = types.StringValue("/tmp/a.raw")
				m.ImageSHA256 = types.StringValue("deadbeef")
				m.Architecture = types.StringValue("x86")
			},
			want: func(r UploadRequest) bool {
				return r.ImagePath == "/tmp/a.raw" && r.ImageSHA256 == "deadbeef" && r.Architecture == "x86"
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := newImageModel(tc.mut)
			req, err := m.toUploadRequest(context.Background())
			if err != nil {
				t.Fatalf("toUploadRequest: %v", err)
			}
			if !tc.want(req) {
				t.Errorf("mapping mismatch: %+v", req)
			}
		})
	}
}

// TestImageSchema_ForceNewVsInPlace inspects the schema to confirm the
// ForceNew attributes carry a plan modifier and description/labels do not.
func TestImageSchema_ForceNewVsInPlace(t *testing.T) {
	var resp resource.SchemaResponse
	NewImageResource().Schema(context.Background(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}

	forceNew := []string{"image_url", "image_path", "image_sha256", "architecture",
		"compression", "format", "server_type", "location", "image_size"}
	for _, name := range forceNew {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("missing attribute %q", name)
			continue
		}
		if countPlanModifiers(attr) == 0 {
			t.Errorf("attribute %q should have a ForceNew plan modifier", name)
		}
	}

	for _, name := range []string{"description", "labels"} {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("missing attribute %q", name)
			continue
		}
		if countPlanModifiers(attr) != 0 {
			t.Errorf("attribute %q should update in place (no plan modifier)", name)
		}
	}
}

// countPlanModifiers reports how many plan modifiers an attribute declares,
// across the string/int64/map kinds this resource uses.
func countPlanModifiers(a schema.Attribute) int {
	switch v := a.(type) {
	case schema.StringAttribute:
		return len(v.PlanModifiers)
	case schema.Int64Attribute:
		return len(v.PlanModifiers)
	case schema.MapAttribute:
		return len(v.PlanModifiers)
	default:
		return 0
	}
}
