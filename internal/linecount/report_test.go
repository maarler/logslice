package linecount

import (
	"strings"
	"testing"
	"time"
)

func TestWriteWindowReport_Empty(t *testing.T) {
	var buf strings.Builder
	if err := WriteWindowReport(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no windows") {
		t.Errorf("expected '(no windows)' message, got: %q", buf.String())
	}
}

func TestWriteWindowReport_ContainsHeaders(t *testing.T) {
	stats := []*WindowStats{
		{Key: "2024-01-01T00:00", Start: epoch, End: epoch.Add(time.Minute), Lines: 5, Bytes: 100},
	}
	var buf strings.Builder
	if err := WriteWindowReport(&buf, stats); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"WINDOW", "LINES", "BYTES", "DISTRIBUTION"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output:\n%s", hdr, out)
		}
	}
}

func TestWriteWindowReport_ContainsKey(t *testing.T) {
	stats := []*WindowStats{
		{Key: "bucket-A", Lines: 10, Bytes: 200},
		{Key: "bucket-B", Lines: 5, Bytes: 80},
	}
	var buf strings.Builder
	_ = WriteWindowReport(&buf, stats)
	out := buf.String()
	if !strings.Contains(out, "bucket-A") {
		t.Errorf("missing bucket-A in output")
	}
	if !strings.Contains(out, "bucket-B") {
		t.Errorf("missing bucket-B in output")
	}
}

func TestBuildBar_MaxValue(t *testing.T) {
	bar := buildBar(10, 10, 20)
	if len(bar) != 20 {
		t.Errorf("expected bar length 20, got %d", len(bar))
	}
}

func TestBuildBar_ZeroMax(t *testing.T) {
	bar := buildBar(0, 0, 20)
	if bar != "" {
		t.Errorf("expected empty bar for zero max, got %q", bar)
	}
}

func TestBuildBar_MinWidth(t *testing.T) {
	// value is 1, max is 1000 — should still render at least MinBarWidth
	bar := buildBar(1, 1000, 40)
	if len(bar) < MinBarWidth {
		t.Errorf("bar too short: %d", len(bar))
	}
}

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		input int64
		want  string
	}{
		{0, "0B"},
		{512, "512B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{1048576, "1.0MB"},
	}
	for _, tc := range cases {
		got := formatBytes(tc.input)
		if got != tc.want {
			t.Errorf("formatBytes(%d) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
