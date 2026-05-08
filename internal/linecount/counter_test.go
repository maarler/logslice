package linecount_test

import (
	"os"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/linecount"
)

func TestCountReader_Empty(t *testing.T) {
	res, err := linecount.CountReader(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Lines != 0 {
		t.Errorf("expected 0 lines, got %d", res.Lines)
	}
}

func TestCountReader_Lines(t *testing.T) {
	input := "line one\nline two\nline three\n"
	res, err := linecount.CountReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Lines != 3 {
		t.Errorf("expected 3 lines, got %d", res.Lines)
	}
	if res.Bytes != int64(len(input)) {
		t.Errorf("expected %d bytes, got %d", len(input), res.Bytes)
	}
}

func TestCountFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	content := "alpha\nbeta\ngamma\n"
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write: %v", err)
	}
	f.Close()

	res, err := linecount.CountFile(f.Name())
	if err != nil {
		t.Fatalf("CountFile: %v", err)
	}
	if res.Lines != 3 {
		t.Errorf("expected 3 lines, got %d", res.Lines)
	}
}

func TestCountFile_Missing(t *testing.T) {
	_, err := linecount.CountFile("/nonexistent/path/file.log")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestFraction(t *testing.T) {
	cases := []struct {
		done, total int64
		want        float64
	}{
		{0, 0, 0},
		{0, 100, 0},
		{50, 100, 0.5},
		{100, 100, 1.0},
		{200, 100, 1.0}, // clamped
	}
	for _, tc := range cases {
		got := linecount.Fraction(tc.done, tc.total)
		if got != tc.want {
			t.Errorf("Fraction(%d,%d) = %v, want %v", tc.done, tc.total, got, tc.want)
		}
	}
}
