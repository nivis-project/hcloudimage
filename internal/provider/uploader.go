package provider

import "context"

// CreatedByLabelKey / CreatedByLabelValue is the label the upstream
// hcloud-upload-image library always adds to snapshots it creates
// (BRIEFING.md Appendix). Both the real uploader and the fake surface it in
// effective_labels so state matches reality.
const (
	CreatedByLabelKey   = "apricote.de/created-by"
	CreatedByLabelValue = "hcloud-upload-image"
)

// Allowed value sets for the enumerated schema attributes. Kept here so the
// validators, the schema, and the real uploader's library mapping all agree on
// one source of truth.
var (
	ValidArchitectures = []string{"x86", "arm"}
	ValidCompressions  = []string{"none", "bz2", "xz", "zstd"}
	ValidFormats       = []string{"raw", "qcow2"}
)

// UploadRequest is the config-shaped input to an upload, decoupled from the
// hcloud library's option types (BRIEFING.md §4.2).
type UploadRequest struct {
	ImageURL     string
	ImagePath    string
	ImageSHA256  string
	Architecture string
	Compression  string
	Format       string
	ServerType   string
	Location     string
	ImageSize    int64
	Description  string
	Labels       map[string]string
}

// SnapshotInfo describes a snapshot resolved from the cloud (or the fake).
type SnapshotInfo struct {
	ID           int64
	Name         string
	Description  string
	Architecture string
	Labels       map[string]string
	Created      string
}

// Uploader is the seam the resource and data source depend on. The provider
// injects a concrete implementation: the in-memory fake for unit/hermetic tests
// (this milestone), the real hcloudimages/v2-backed one in milestone 04.
type Uploader interface {
	// Upload creates a snapshot from the request and returns its ID and the
	// effective label set (user labels merged with library defaults).
	Upload(ctx context.Context, req UploadRequest) (imageID int64, effectiveLabels map[string]string, err error)

	// Delete removes the snapshot with the given ID.
	Delete(ctx context.Context, imageID int64) error

	// Get returns the snapshot by ID, or (nil, nil) if it no longer exists.
	Get(ctx context.Context, imageID int64) (*SnapshotInfo, error)

	// UpdateMetadata updates description and labels in place (no re-upload),
	// returning the recomputed effective labels.
	UpdateMetadata(ctx context.Context, imageID int64, description string, labels map[string]string) (effectiveLabels map[string]string, err error)

	// Find resolves a data-source lookup by ID or by label selector.
	// When the selector matches more than one snapshot, it errors unless
	// mostRecent is true, in which case it returns the newest.
	Find(ctx context.Context, byID int64, selector string, mostRecent bool) (*SnapshotInfo, error)
}

// mergeEffectiveLabels returns the user labels with the library's created-by
// label merged in, without mutating the input.
func mergeEffectiveLabels(userLabels map[string]string) map[string]string {
	out := make(map[string]string, len(userLabels)+1)
	for k, v := range userLabels {
		out[k] = v
	}
	out[CreatedByLabelKey] = CreatedByLabelValue
	return out
}
