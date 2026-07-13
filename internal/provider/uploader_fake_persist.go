package provider

import (
	"encoding/json"
	"os"
)

// The fake uploader normally lives only in memory. But the hermetic NixOS-VM
// test (BRIEFING.md §8.2) runs each terraform/tofu command as a fresh provider
// process, so the fake must persist its store between invocations for the
// lifecycle (create → replace → update → destroy) to be observable. When
// HCLOUDIMAGE_FAKE_STATE names a file, the fake loads it on construction and
// rewrites it after every mutation.

// persistState is the on-disk shape of the fake's store.
type persistState struct {
	NextID    int64                     `json:"next_id"`
	Snapshots map[int64]persistSnapshot `json:"snapshots"`
	Order     []int64                   `json:"order"`
}

type persistSnapshot struct {
	Info    SnapshotInfo `json:"info"`
	Deleted bool         `json:"deleted"`
}

// statePath returns the persistence file path, or "" for pure in-memory mode.
func statePath() string { return os.Getenv("HCLOUDIMAGE_FAKE_STATE") }

// load reads persisted state into the fake if HCLOUDIMAGE_FAKE_STATE is set and
// the file exists. Caller holds the lock (or is constructing).
func (f *FakeUploader) load() {
	path := statePath()
	if path == "" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return // no file yet — start empty
	}
	var ps persistState
	if json.Unmarshal(data, &ps) != nil {
		return
	}
	if ps.NextID != 0 {
		f.nextID = ps.NextID
	}
	f.store = make(map[int64]*fakeSnapshot, len(ps.Snapshots))
	for id, s := range ps.Snapshots {
		snap := s
		f.store[id] = &fakeSnapshot{info: snap.Info, deleted: snap.Deleted}
	}
	f.createdOrder = ps.Order
}

// save writes the fake's store to disk when HCLOUDIMAGE_FAKE_STATE is set.
// Caller holds the lock.
func (f *FakeUploader) save() {
	path := statePath()
	if path == "" {
		return
	}
	ps := persistState{
		NextID:    f.nextID,
		Snapshots: make(map[int64]persistSnapshot, len(f.store)),
		Order:     f.createdOrder,
	}
	for id, snap := range f.store {
		ps.Snapshots[id] = persistSnapshot{Info: snap.info, Deleted: snap.deleted}
	}
	data, err := json.Marshal(ps)
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o600)
}
