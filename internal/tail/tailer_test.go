package tail_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/tail"
)

func writeFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestTailer_ReadsAllLines(t *testing.T) {
	path := writeFile(t, "line1\nline2\nline3\n")
	tr := tail.New(path, tail.DefaultOptions())

	ctx := context.Background()
	lines, errCh := tr.Lines(ctx)

	var got []string
	for l := range lines {
		got = append(got, l)
	}
	if err := <-errCh; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(got), got)
	}
	if got[0] != "line1" || got[1] != "line2" || got[2] != "line3" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestTailer_EmptyFile(t *testing.T) {
	path := writeFile(t, "")
	tr := tail.New(path, tail.DefaultOptions())

	ctx := context.Background()
	lines, errCh := tr.Lines(ctx)

	var got []string
	for l := range lines {
		got = append(got, l)
	}
	if err := <-errCh; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no lines, got %v", got)
	}
}

func TestTailer_MissingFile(t *testing.T) {
	tr := tail.New("/nonexistent/path/file.log", tail.DefaultOptions())
	ctx := context.Background()
	lines, errCh := tr.Lines(ctx)
	for range lines {
	}
	if err := <-errCh; err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestTailer_ContextCancel_StopsFollow(t *testing.T) {
	path := writeFile(t, "hello\n")
	opts := tail.Options{Follow: true, PollInterval: 20 * time.Millisecond}
	tr := tail.New(path, opts)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	lines, _ := tr.Lines(ctx)
	for range lines {
	}
	// Test passes if it completes within the deadline (context cancels the follow loop).
}
