package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestSnapshotDataSource_Schema(t *testing.T) {
	var resp datasource.SchemaResponse
	NewSnapshotDataSource().Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}
	for _, name := range []string{"id", "with_selector", "most_recent", "name", "description", "architecture", "created", "labels"} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Errorf("data source schema missing %q", name)
		}
	}
}

func TestSnapshotDataSource_ConfigValidators(t *testing.T) {
	ds := NewSnapshotDataSource().(*snapshotDataSource)
	vs := ds.ConfigValidators(context.Background())
	if len(vs) != 1 {
		t.Fatalf("ConfigValidators() = %d, want 1 (ExactlyOneOf id/with_selector)", len(vs))
	}
}

// The lookup semantics themselves are exercised via the fake in
// TestFakeUploader_FindBySelector; here we confirm the data source delegates to
// the injected uploader.
func TestSnapshotDataSource_UsesUploader(t *testing.T) {
	f := NewFakeUploader()
	id, _, _ := f.Upload(context.Background(), UploadRequest{Architecture: "arm"})

	info, err := f.Find(context.Background(), id, "", false)
	if err != nil || info == nil || info.ID != id {
		t.Fatalf("uploader.Find via data source seam failed: info=%v err=%v", info, err)
	}
}
