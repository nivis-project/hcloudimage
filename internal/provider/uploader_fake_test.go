package provider

import (
	"context"
	"testing"
)

func TestFakeUploader_UploadGetDelete(t *testing.T) {
	ctx := context.Background()
	f := NewFakeUploader()

	id, effective, err := f.Upload(ctx, UploadRequest{Architecture: "x86", Labels: map[string]string{"os": "nixos"}})
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if id == 0 {
		t.Fatal("Upload returned zero id")
	}
	if effective[CreatedByLabelKey] != CreatedByLabelValue {
		t.Errorf("effective labels missing created-by: %v", effective)
	}
	if effective["os"] != "nixos" {
		t.Errorf("effective labels missing user label: %v", effective)
	}

	info, err := f.Get(ctx, id)
	if err != nil || info == nil {
		t.Fatalf("Get after upload: info=%v err=%v", info, err)
	}
	if info.Architecture != "x86" {
		t.Errorf("architecture = %q, want x86", info.Architecture)
	}

	if err := f.Delete(ctx, id); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	info, err = f.Get(ctx, id)
	if err != nil {
		t.Fatalf("Get after delete: %v", err)
	}
	if info != nil {
		t.Error("Get after delete should return nil")
	}
}

func TestFakeUploader_OutOfBandDeletion(t *testing.T) {
	ctx := context.Background()
	f := NewFakeUploader()

	id, _, _ := f.Upload(ctx, UploadRequest{Architecture: "arm"})
	f.MarkDeleted(id)

	info, err := f.Get(ctx, id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if info != nil {
		t.Error("Get should report out-of-band-deleted snapshot as gone")
	}
}

func TestFakeUploader_UpdateMetadata(t *testing.T) {
	ctx := context.Background()
	f := NewFakeUploader()

	id, _, _ := f.Upload(ctx, UploadRequest{Architecture: "x86", Description: "old"})
	effective, err := f.UpdateMetadata(ctx, id, "new", map[string]string{"env": "ci"})
	if err != nil {
		t.Fatalf("UpdateMetadata: %v", err)
	}
	if effective["env"] != "ci" || effective[CreatedByLabelKey] != CreatedByLabelValue {
		t.Errorf("effective labels wrong: %v", effective)
	}
	info, _ := f.Get(ctx, id)
	if info.Description != "new" {
		t.Errorf("description = %q, want new", info.Description)
	}
	if len(f.UpdateCalls) != 1 {
		t.Errorf("UpdateCalls = %d, want 1", len(f.UpdateCalls))
	}
}

func TestFakeUploader_FindByID(t *testing.T) {
	ctx := context.Background()
	f := NewFakeUploader()
	id, _, _ := f.Upload(ctx, UploadRequest{Architecture: "x86"})

	info, err := f.Find(ctx, id, "", false)
	if err != nil || info == nil {
		t.Fatalf("Find by id: info=%v err=%v", info, err)
	}
	if info.ID != id {
		t.Errorf("Find returned id %d, want %d", info.ID, id)
	}
}

func TestFakeUploader_FindBySelector(t *testing.T) {
	ctx := context.Background()
	f := NewFakeUploader()
	if _, _, err := f.Upload(ctx, UploadRequest{Architecture: "x86", Labels: map[string]string{"role": "base"}}); err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if _, _, err := f.Upload(ctx, UploadRequest{Architecture: "x86", Labels: map[string]string{"role": "base"}}); err != nil {
		t.Fatalf("Upload: %v", err)
	}

	// Ambiguous without most_recent -> error.
	if _, err := f.Find(ctx, 0, "role=base", false); err == nil {
		t.Error("ambiguous selector without most_recent should error")
	}

	// With most_recent -> newest (highest id).
	info, err := f.Find(ctx, 0, "role=base", true)
	if err != nil {
		t.Fatalf("Find most_recent: %v", err)
	}
	if info.ID != 1001 {
		t.Errorf("most_recent returned id %d, want newest 1001", info.ID)
	}

	// No match -> error.
	if _, err := f.Find(ctx, 0, "role=missing", true); err == nil {
		t.Error("no-match selector should error")
	}
}
