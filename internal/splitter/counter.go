package splitter

import (
	"fmt"
	"io"
	"time"
)

// SegmentInfo holds metadata about a single log segment.
type SegmentInfo struct {
	Path      string
	Window    time.Time
	LineCount int
}

// Summary holds aggregate statistics from a split operation.
type Summary struct {
	TotalLines    int
	SegmentCount  int
	SkippedLines  int
	Segments      []SegmentInfo
}

// String returns a human-readable summary of the split operation.
func (s Summary) String() string {
	return fmt.Sprintf(
		"Split complete: %d lines → %d segments (%d lines skipped)",
		s.TotalLines, s.SegmentCount, s.SkippedLines,
	)
}

// WriteSummary writes a formatted summary table to w.
func WriteSummary(w io.Writer, s Summary) {
	fmt.Fprintf(w, "%-40s %10s %15s\n", "Segment", "Lines", "Window")
	fmt.Fprintf(w, "%s\n", repeatChar('-', 70))
	for _, seg := range s.Segments {
		fmt.Fprintf(w, "%-40s %10d %15s\n",
			seg.Path,
			seg.LineCount,
			seg.Window.UTC().Format("2006-01-02 15:04"),
		)
	}
	fmt.Fprintf(w, "%s\n", repeatChar('-', 70))
	fmt.Fprintf(w, "Total: %d lines in %d segments\n", s.TotalLines, s.SegmentCount)
	if s.SkippedLines > 0 {
		fmt.Fprintf(w, "Skipped (no timestamp): %d lines\n", s.SkippedLines)
	}
}

func repeatChar(c rune, n int) string {
	buf := make([]rune, n)
	for i := range buf {
		buf[i] = c
	}
	return string(buf)
}
