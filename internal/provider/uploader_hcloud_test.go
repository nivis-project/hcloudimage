package provider

import (
	"testing"

	"github.com/apricote/hcloud-upload-image/hcloudimages/v2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLibCompressionMapping(t *testing.T) {
	cases := map[string]hcloudimages.Compression{
		"none": hcloudimages.CompressionNone,
		"bz2":  hcloudimages.CompressionBZ2,
		"xz":   hcloudimages.CompressionXZ,
		"zstd": hcloudimages.CompressionZSTD,
	}
	for in, want := range cases {
		if got := libCompression(in); got != want {
			t.Errorf("libCompression(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLibFormatMapping(t *testing.T) {
	if got := libFormat("raw"); got != hcloudimages.FormatRaw {
		t.Errorf("libFormat(raw) = %q, want raw", got)
	}
	if got := libFormat("qcow2"); got != hcloudimages.FormatQCOW2 {
		t.Errorf("libFormat(qcow2) = %q, want qcow2", got)
	}
}

func TestLibArchitectureMapping(t *testing.T) {
	if got := libArchitecture("x86"); got != hcloud.ArchitectureX86 {
		t.Errorf("libArchitecture(x86) = %q, want x86", got)
	}
	if got := libArchitecture("arm"); got != hcloud.ArchitectureARM {
		t.Errorf("libArchitecture(arm) = %q, want arm", got)
	}
}

func TestNewHcloudUploader_InvalidPollInterval(t *testing.T) {
	_, err := newHcloudUploader(providerConfig{Token: "t", PollInterval: "not-a-duration"})
	if err == nil {
		t.Error("expected error for invalid poll_interval")
	}
}

func TestToSnapshotInfo(t *testing.T) {
	info := toSnapshotInfo(&hcloud.Image{
		ID:           42,
		Name:         "snap",
		Description:  "desc",
		Architecture: hcloud.ArchitectureARM,
		Labels:       map[string]string{"k": "v"},
	})
	if info.ID != 42 || info.Architecture != "arm" || info.Labels["k"] != "v" {
		t.Errorf("toSnapshotInfo mapping wrong: %+v", info)
	}
}
