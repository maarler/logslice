package rotate

import (
	"os"
	"path/filepath"
	"testing"
)

func writeContent(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeContent: %v", err)
	}
}

func TestNewDetector_MissingFile(t *testing.T) {
	_, err := NewDetector("/nonexistent/path/logfile.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestDetector_NoRotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	writeContent(t, path, "line1\nline2\n")

	d, err := NewDetector(path)
	if err != nil {
		t.Fatalf("NewDetector: %v", err)
	}

	// Append more data — not a rotation.
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	_, _ = f.WriteString("line3\n")
	f.Close()

	rotated, err := d.Check()
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if rotated {
		t.Error("expected rotated=false after append, got true")
	}
}

func TestDetector_Truncation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	writeContent(t, path, "lots of log data here\n")

	d, err := NewDetector(path)
	if err != nil {
		t.Fatalf("NewDetector: %v", err)
	}

	// Truncate the file.
	writeContent(t, path, "")

	rotated, err := d.Check()
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !rotated {
		t.Error("expected rotated=true after truncation, got false")
	}
}

func TestDetector_FileReplaced(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	writeContent(t, path, "original content\n")

	d, err := NewDetector(path)
	if err != nil {
		t.Fatalf("NewDetector: %v", err)
	}

	// Replace file (new inode).
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove: %v", err)
	}
	writeContent(t, path, "brand new file\n")

	rotated, err := d.Check()
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !rotated {
		t.Error("expected rotated=true after file replacement, got false")
	}
}

func TestDetector_Reset(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	writeContent(t, path, "initial\n")

	d, err := NewDetector(path)
	if err != nil {
		t.Fatalf("NewDetector: %v", err)
	}

	writeContent(t, path, "")
	if err := d.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}

	rotated, err := d.Check()
	if err != nil {
		t.Fatalf("Check after Reset: %v", err)
	}
	if rotated {
		t.Error("expected rotated=false after Reset, got true")
	}
}
