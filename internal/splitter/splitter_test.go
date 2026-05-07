package splitter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/splitter"
)

const sampleLogs = `2024-01-15 10:00:01 INFO starting service
2024-01-15 10:00:45 DEBUG connection established
2024-01-15 10:01:10 INFO request received
2024-01-15 10:02:05 WARN slow query detected
2024-01-15 10:03:30 ERROR timeout occurred
`

func TestSplit_ByMinute(t *testing.T) {
	dir := t.TempDir()
	s := splitter.New(splitter.Options{
		WindowSize: time.Minute,
		OutputDir:  dir,
		Prefix:     "test",
	})

	r := strings.NewReader(sampleLogs)
	files, err := s.Split(r)
	if err != nil {
		t.Fatalf("Split() error = %v", err)
	}

	if len(files) != 4 {
		t.Errorf("expected 4 segment files, got %d", len(files))
	}
	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", f)
		}
	}
}

func TestSplit_EmptyInput(t *testing.T) {
	dir := t.TempDir()
	s := splitter.New(splitter.Options{
		WindowSize: time.Minute,
		OutputDir:  dir,
	})

	files, err := s.Split(strings.NewReader(""))
	if err != nil {
		t.Fatalf("Split() error = %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files for empty input, got %d", len(files))
	}
}

func TestSplit_DefaultOptions(t *testing.T) {
	dir := t.TempDir()
	s := splitter.New(splitter.Options{
		WindowSize: 5 * time.Minute,
		OutputDir:  dir,
	})

	r := strings.NewReader(sampleLogs)
	files, err := s.Split(r)
	if err != nil {
		t.Fatalf("Split() error = %v", err)
	}

	for _, f := range files {
		base := filepath.Base(f)
		if !strings.HasPrefix(base, "segment_") {
			t.Errorf("expected default prefix 'segment_', got %s", base)
		}
	}
}
