package provider

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// fakeSnapshot is a snapshot stored in the fake uploader.
type fakeSnapshot struct {
	info    SnapshotInfo
	deleted bool // marked deleted "out of band" (Get/Find behave as not-found)
}

// FakeUploader is an in-memory Uploader for unit and hermetic lifecycle tests
// (BRIEFING.md §4.1). It records calls, hands out synthetic incrementing IDs,
// mirrors the library's created-by label merge, and can simulate a snapshot
// being deleted out of band.
type FakeUploader struct {
	mu     sync.Mutex
	nextID int64
	store  map[int64]*fakeSnapshot

	// Call recorders, for assertions in tests.
	UploadCalls  []UploadRequest
	DeleteCalls  []int64
	UpdateCalls  []int64
	createdOrder []int64 // insertion order, so most_recent is deterministic
}

// NewFakeUploader returns a ready-to-use fake.
func NewFakeUploader() *FakeUploader {
	return &FakeUploader{
		nextID: 1000,
		store:  make(map[int64]*fakeSnapshot),
	}
}

var _ Uploader = (*FakeUploader)(nil)

func (f *FakeUploader) Upload(_ context.Context, req UploadRequest) (int64, map[string]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.UploadCalls = append(f.UploadCalls, req)

	id := f.nextID
	f.nextID++

	effective := mergeEffectiveLabels(req.Labels)
	f.store[id] = &fakeSnapshot{
		info: SnapshotInfo{
			ID:           id,
			Name:         fmt.Sprintf("snapshot-%d", id),
			Description:  req.Description,
			Architecture: req.Architecture,
			Labels:       effective,
			Created:      fmt.Sprintf("2026-01-01T00:00:%02dZ", int(id%60)),
		},
	}
	f.createdOrder = append(f.createdOrder, id)
	return id, effective, nil
}

func (f *FakeUploader) Delete(_ context.Context, imageID int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.DeleteCalls = append(f.DeleteCalls, imageID)
	delete(f.store, imageID)
	return nil
}

func (f *FakeUploader) Get(_ context.Context, imageID int64) (*SnapshotInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	snap, ok := f.store[imageID]
	if !ok || snap.deleted {
		return nil, nil
	}
	info := snap.info
	return &info, nil
}

func (f *FakeUploader) UpdateMetadata(_ context.Context, imageID int64, description string, labels map[string]string) (map[string]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.UpdateCalls = append(f.UpdateCalls, imageID)

	snap, ok := f.store[imageID]
	if !ok || snap.deleted {
		return nil, fmt.Errorf("snapshot %d not found", imageID)
	}
	effective := mergeEffectiveLabels(labels)
	snap.info.Description = description
	snap.info.Labels = effective
	return effective, nil
}

func (f *FakeUploader) Find(_ context.Context, byID int64, selector string, mostRecent bool) (*SnapshotInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if byID != 0 {
		snap, ok := f.store[byID]
		if !ok || snap.deleted {
			return nil, fmt.Errorf("snapshot %d not found", byID)
		}
		info := snap.info
		return &info, nil
	}

	// Selector match: a label "k=v" is present in the snapshot's labels.
	matches := f.matchSelector(selector)
	switch {
	case len(matches) == 0:
		return nil, fmt.Errorf("no snapshot matches selector %q", selector)
	case len(matches) == 1:
		info := f.store[matches[0]].info
		return &info, nil
	default:
		if !mostRecent {
			return nil, fmt.Errorf("selector %q matched %d snapshots; set most_recent = true to disambiguate", selector, len(matches))
		}
		// Newest = highest insertion order among matches.
		sort.Slice(matches, func(i, j int) bool { return matches[i] > matches[j] })
		info := f.store[matches[0]].info
		return &info, nil
	}
}

// matchSelector returns the IDs of live snapshots whose labels satisfy a simple
// "key=value" selector (sufficient for tests).
func (f *FakeUploader) matchSelector(selector string) []int64 {
	parts := strings.SplitN(selector, "=", 2)
	if len(parts) != 2 {
		return nil
	}
	k, v := parts[0], parts[1]
	var out []int64
	for id, snap := range f.store {
		if snap.deleted {
			continue
		}
		if snap.info.Labels[k] == v {
			out = append(out, id)
		}
	}
	return out
}

// MarkDeleted simulates an out-of-band deletion: the snapshot stays in the map
// but Get/Find report it as gone. Test helper.
func (f *FakeUploader) MarkDeleted(imageID int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if snap, ok := f.store[imageID]; ok {
		snap.deleted = true
	}
}
