package splitter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/splitter"
)

func TestSummary_String(t *testing.T) {
	s := splitter.Summary{
		TotalLines:   100,
		SegmentCount: 4,
		SkippedLines: 3,
	}
	out := s.String()
	if !strings.Contains(out, "100") {
		t.Errorf("expected total lines in summary string, got: %s", out)
	}
	if !strings.Contains(out, "4") {
		t.Errorf("expected segment count in summary string, got: %s", out)
	}
	if !strings.Contains(out, "3") {
		t.Errorf("expected skipped lines in summary string, got: %s", out)
	}
}

func TestWriteSummary(t *testing.T) {
	w := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	s := splitter.Summary{
		TotalLines:   50,
		SegmentCount: 2,
		SkippedLines: 1,
		Segments: []splitter.SegmentInfo{
			{Path: "segment_20240115T100000Z.log", Window: w, LineCount: 30},
			{Path: "segment_20240115T100100Z.log", Window: w.Add(time.Minute), LineCount: 20},
		},
	}

	var buf bytes.Buffer
	splitter.WriteSummary(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "segment_20240115T100000Z.log") {
		t.Errorf("expected first segment path in output")
	}
	if !strings.Contains(out, "30") {
		t.Errorf("expected line count 30 in output")
	}
	if !strings.Contains(out, "Skipped") {
		t.Errorf("expected skipped lines note in output")
	}
}
