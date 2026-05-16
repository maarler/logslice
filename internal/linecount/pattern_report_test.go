package linecount

import (
	"bytes"
	"strings"
	"testing"
)

func TestSortedPatternEntries_ByCountDesc(t *testing.T) {
	counts := map[string]int64{"warn": 5, "error": 10, "info": 1}
	entries := SortedPatternEntries(counts)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Name != "error" || entries[0].Count != 10 {
		t.Errorf("first entry should be error/10, got %s/%d", entries[0].Name, entries[0].Count)
	}
	if entries[1].Name != "warn" || entries[1].Count != 5 {
		t.Errorf("second entry should be warn/5, got %s/%d", entries[1].Name, entries[1].Count)
	}
}

func TestSortedPatternEntries_TieBreakByName(t *testing.T) {
	counts := map[string]int64{"beta": 3, "alpha": 3}
	entries := SortedPatternEntries(counts)
	if entries[0].Name != "alpha" {
		t.Errorf("tie should break alphabetically, got %s first", entries[0].Name)
	}
}

func TestWritePatternReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WritePatternReport(&buf, map[string]int64{}, 0); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no pattern") {
		t.Error("expected 'no pattern' message for empty counts")
	}
}

func TestWritePatternReport_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	counts := map[string]int64{"error": 7}
	if err := WritePatternReport(&buf, counts, 100); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, hdr := range []string{"PATTERN", "COUNT", "PCT"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestWritePatternReport_PercentageCalculation(t *testing.T) {
	var buf bytes.Buffer
	counts := map[string]int64{"error": 50}
	if err := WritePatternReport(&buf, counts, 200); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "25.0%") {
		t.Errorf("expected 25.0%% in output, got:\n%s", buf.String())
	}
}

func TestWritePatternReport_ZeroTotal(t *testing.T) {
	var buf bytes.Buffer
	counts := map[string]int64{"x": 5}
	if err := WritePatternReport(&buf, counts, 0); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "0.0%") {
		t.Errorf("expected 0.0%% when total is zero, got:\n%s", buf.String())
	}
}
