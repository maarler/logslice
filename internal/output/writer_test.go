package output

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriter_DryRun(t *testing.T) {
	var buf bytes.Buffer
	w := New("/tmp", "test-")
	w.DryRun = true
	w.Stdout = &buf

	if err := w.Write("2024-01-01", "hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[2024-01-01]") {
		t.Errorf("expected key in output, got: %q", out)
	}
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected line in output, got: %q", out)
	}
}

func TestWriter_WritesFile(t *testing.T) {
	dir := t.TempDir()
	w := New(dir, "seg-")

	lines := []string{"first line", "second line", "third line"}
	for _, l := range lines {
		if err := w.Write("2024-01-02T10", l); err != nil {
			t.Fatalf("Write error: %v", err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close error: %v", err)
	}

	path := filepath.Join(dir, "seg-2024-01-02T10.log")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	for _, l := range lines {
		if !strings.Contains(string(data), l) {
			t.Errorf("expected %q in file content", l)
		}
	}
}

func TestWriter_MultipleKeys(t *testing.T) {
	dir := t.TempDir()
	w := New(dir, "")

	_ = w.Write("key-a", "line a")
	_ = w.Write("key-b", "line b")
	_ = w.Write("key-a", "line a2")

	keys := w.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
	_ = w.Close()
}

func TestWriter_Close_Idempotent(t *testing.T) {
	dir := t.TempDir()
	w := New(dir, "")
	_ = w.Write("k", "line")
	if err := w.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
}
