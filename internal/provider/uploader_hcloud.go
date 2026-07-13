package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sort"
	"time"

	"github.com/apricote/hcloud-upload-image/hcloudimages/v2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// hcloudUploader is the production Uploader, wrapping hcloudimages/v2 for the
// upload-to-snapshot flow and the hcloud Image API for lifecycle operations
// (BRIEFING.md §4.2, Appendix). Verified against hcloudimages/v2 v2.0.1.
type hcloudUploader struct {
	client *hcloud.Client
	images *hcloudimages.Client
}

var _ Uploader = (*hcloudUploader)(nil)

// newHcloudUploader builds the real uploader from resolved provider config.
func newHcloudUploader(cfg providerConfig) (*hcloudUploader, error) {
	opts := []hcloud.ClientOption{hcloud.WithToken(cfg.Token)}
	if cfg.Endpoint != "" {
		opts = append(opts, hcloud.WithEndpoint(cfg.Endpoint))
	}
	if cfg.PollInterval != "" {
		d, err := time.ParseDuration(cfg.PollInterval)
		if err != nil {
			return nil, fmt.Errorf("invalid poll_interval %q: %w", cfg.PollInterval, err)
		}
		opts = append(opts, hcloud.WithPollOpts(hcloud.PollOpts{
			BackoffFunc: hcloud.ConstantBackoff(d),
		}))
	}

	client := hcloud.NewClient(opts...)
	return &hcloudUploader{client: client, images: hcloudimages.NewClient(client)}, nil
}

// compression/format/architecture mapping to the library's typed constants.
func libCompression(s string) hcloudimages.Compression {
	switch s {
	case "bz2":
		return hcloudimages.CompressionBZ2
	case "xz":
		return hcloudimages.CompressionXZ
	case "zstd":
		return hcloudimages.CompressionZSTD
	default:
		return hcloudimages.CompressionNone
	}
}

func libFormat(s string) hcloudimages.Format {
	if s == "qcow2" {
		return hcloudimages.FormatQCOW2
	}
	return hcloudimages.FormatRaw
}

func libArchitecture(s string) hcloud.Architecture {
	if s == "arm" {
		return hcloud.ArchitectureARM
	}
	return hcloud.ArchitectureX86
}

func (u *hcloudUploader) Upload(ctx context.Context, req UploadRequest) (int64, map[string]string, error) {
	write := hcloudimages.WriteOptions{
		ImageCompression: libCompression(req.Compression),
		ImageFormat:      libFormat(req.Format),
		ImageSize:        req.ImageSize,
	}

	switch {
	case req.ImageURL != "":
		parsed, err := url.Parse(req.ImageURL)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid image_url: %w", err)
		}
		write.ImageURL = parsed
	case req.ImagePath != "":
		f, err := os.Open(req.ImagePath)
		if err != nil {
			return 0, nil, fmt.Errorf("opening image_path: %w", err)
		}
		defer func() { _ = f.Close() }()
		write.ImageReader = f
	default:
		return 0, nil, fmt.Errorf("neither image_url nor image_path set")
	}

	opts := hcloudimages.UploadOptions{
		WriteOptions: write,
		Architecture: libArchitecture(req.Architecture),
		Labels:       req.Labels,
	}
	if req.Description != "" {
		opts.Description = &req.Description
	}
	// Only override server_type/location when the user set them; otherwise the
	// library applies its per-architecture defaults.
	if req.ServerType != "" {
		opts.ServerType = &hcloud.ServerType{Name: req.ServerType}
	}
	if req.Location != "" && req.Location != "fsn1" {
		opts.Location = &hcloud.Location{Name: req.Location}
	}

	// Debug-only cleanup skip, env-gated — never a public schema attribute.
	if os.Getenv("HCLOUDIMAGE_DEBUG_SKIP_CLEANUP") == "1" {
		opts.DebugSkipResourceCleanup = true
	}

	image, err := u.images.Upload(ctx, opts)
	if err != nil {
		return 0, nil, err
	}
	return image.ID, image.Labels, nil
}

func (u *hcloudUploader) Delete(ctx context.Context, imageID int64) error {
	_, err := u.client.Image.Delete(ctx, &hcloud.Image{ID: imageID})
	return err
}

func (u *hcloudUploader) Get(ctx context.Context, imageID int64) (*SnapshotInfo, error) {
	image, _, err := u.client.Image.GetByID(ctx, imageID)
	if err != nil {
		return nil, err
	}
	if image == nil {
		return nil, nil
	}
	return toSnapshotInfo(image), nil
}

func (u *hcloudUploader) UpdateMetadata(ctx context.Context, imageID int64, description string, labels map[string]string) (map[string]string, error) {
	opts := hcloud.ImageUpdateOpts{Labels: labels}
	if description != "" {
		opts.Description = &description
	}
	image, _, err := u.client.Image.Update(ctx, &hcloud.Image{ID: imageID}, opts)
	if err != nil {
		return nil, err
	}
	return image.Labels, nil
}

func (u *hcloudUploader) Find(ctx context.Context, byID int64, selector string, mostRecent bool) (*SnapshotInfo, error) {
	if byID != 0 {
		return u.Get(ctx, byID)
	}

	images, err := u.client.Image.AllWithOpts(ctx, hcloud.ImageListOpts{
		ListOpts: hcloud.ListOpts{LabelSelector: selector},
	})
	if err != nil {
		return nil, err
	}
	switch {
	case len(images) == 0:
		return nil, fmt.Errorf("no snapshot matches selector %q", selector)
	case len(images) == 1:
		return toSnapshotInfo(images[0]), nil
	default:
		if !mostRecent {
			return nil, fmt.Errorf("selector %q matched %d snapshots; set most_recent = true to disambiguate", selector, len(images))
		}
		sort.Slice(images, func(i, j int) bool { return images[i].Created.After(images[j].Created) })
		return toSnapshotInfo(images[0]), nil
	}
}

func toSnapshotInfo(image *hcloud.Image) *SnapshotInfo {
	arch := "x86"
	if image.Architecture == hcloud.ArchitectureARM {
		arch = "arm"
	}
	return &SnapshotInfo{
		ID:           image.ID,
		Name:         image.Name,
		Description:  image.Description,
		Architecture: arch,
		Labels:       image.Labels,
		Created:      image.Created.Format(time.RFC3339),
	}
}
