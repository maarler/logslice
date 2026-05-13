package linecount

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountReader_Empty(t *testing.T) {
	n, err := CountReader(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Fatalf("expected 0, got %d", n)
	}
}

func TestCountReader_Lines(t *testing.T) {
	input := "line1\nline2\nline3\n"
	n, err := CountReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3, got %d", n)
	}
}

func TestCountFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	if err := os.WriteFile(path, []byte("a\nb\nc\n"), 0644); err != nil {
		t.Fatal(err)
	}
	n, err := CountFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3, got %d", n)
	}
}

func TestCountFile_Missing(t *testing.T) {
	_, err := CountFile("/no/such/file.log")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFraction(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	lines := "l1\nl2\nl3\nl4\nl5\nl6\nl7\nl8\nl9\nl10\n"
	if err := os.WriteFile(path, []byte(lines), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		frac float64
		want int64
	}{
		{0.0, 0},
		{1.0, 10},
		{0.5, 5},
		{0.1, 1},
	}
	for _, tc := range tests {
		got, err := Fraction(path, tc.frac)
		if err != nil {
			t.Fatalf("frac=%.1f: unexpected error: %v", tc.frac, err)
		}
		if got != tc.want {
			t.Errorf("frac=%.1f: expected %d, got %d", tc.frac, tc.want, got)
		}
	}
}
