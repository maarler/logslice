package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/checkpoint"
)

func tempStore(t *testing.T) *checkpoint.Store {
	t.Helper()
	dir := t.TempDir()
	return checkpoint.NewStore(filepath.Join(dir, "sub", "state.json"))
}

func TestLoad_MissingFile_ReturnsZero(t *testing.T) {
	st, err := tempStore(t).Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st.Offset != 0 || st.Path != "" {
		t.Fatalf("expected zero state, got %+v", st)
	}
}

func TestSave_And_Load_RoundTrip(t *testing.T) {
	store := tempStore(t)
	want := checkpoint.State{Path: "/var/log/app.log", Offset: 4096}

	if err := store.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Path != want.Path {
		t.Errorf("Path: got %q, want %q", got.Path, want.Path)
	}
	if got.Offset != want.Offset {
		t.Errorf("Offset: got %d, want %d", got.Offset, want.Offset)
	}
	if got.SavedAt.IsZero() {
		t.Error("SavedAt should be set after Save")
	}
}

func TestSave_CreatesParentDirectories(t *testing.T) {
	dir := t.TempDir()
	store := checkpoint.NewStore(filepath.Join(dir, "a", "b", "c", "state.json"))

	if err := store.Save(checkpoint.State{Path: "x", Offset: 1}); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestDelete_RemovesFile(t *testing.T) {
	store := tempStore(t)
	_ = store.Save(checkpoint.State{Path: "x", Offset: 10})

	if err := store.Delete(); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// Second delete should be a no-op, not an error.
	if err := store.Delete(); err != nil {
		t.Fatalf("second Delete: %v", err)
	}
}

func TestLoad_CorruptFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "state.json")
	_ = os.WriteFile(file, []byte("not-json{"), 0o644)

	store := checkpoint.NewStore(file)
	_, err := store.Load()
	if err == nil {
		t.Fatal("expected error for corrupt JSON, got nil")
	}
}
