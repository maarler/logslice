package linecount

import (
	"strings"
	"testing"
)

func TestWriteErrorReport_Empty(t *testing.T) {
	var buf strings.Builder
	err := WriteErrorReport(&buf, nil, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no matching lines") {
		t.Fatalf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteErrorReport_ContainsHeaders(t *testing.T) {
	var buf strings.Builder
	entries := []ErrorEntry{{Level: LevelError, Count: 3}}
	if err := WriteErrorReport(&buf, entries, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Level") {
		t.Fatalf("expected 'Level' header, got: %s", out)
	}
	if !strings.Contains(out, "Count") {
		t.Fatalf("expected 'Count' header, got: %s", out)
	}
	if !strings.Contains(out, "Percent") {
		t.Fatalf("expected 'Percent' header, got: %s", out)
	}
}

func TestWriteErrorReport_PercentageCalculation(t *testing.T) {
	var buf strings.Builder
	entries := []ErrorEntry{
		{Level: LevelError, Count: 1},
		{Level: LevelInfo, Count: 3},
	}
	if err := WriteErrorReport(&buf, entries, 4); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "25.0%") {
		t.Fatalf("expected 25.0%% for ERROR, got: %s", out)
	}
	if !strings.Contains(out, "75.0%") {
		t.Fatalf("expected 75.0%% for INFO, got: %s", out)
	}
}

func TestWriteErrorReport_ZeroTotal_NoPanic(t *testing.T) {
	var buf strings.Builder
	entries := []ErrorEntry{{Level: LevelWarn, Count: 2}}
	if err := WriteErrorReport(&buf, entries, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "0.0%") {
		t.Fatalf("expected 0.0%% when total is zero")
	}
}

func TestWriteErrorReport_TotalLine(t *testing.T) {
	var buf strings.Builder
	entries := []ErrorEntry{{Level: LevelFatal, Count: 7}}
	if err := WriteErrorReport(&buf, entries, 7); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Total matched: 7") {
		t.Fatalf("expected total line, got: %s", buf.String())
	}
}
