package linecount

import (
	"strings"
	"testing"
)

func TestNewFieldCounter_Empty(t *testing.T) {
	fc := NewFieldCounter("level")
	if got := fc.Counts(); len(got) != 0 {
		t.Fatalf("expected empty counts, got %v", got)
	}
}

func TestFieldCounter_Add_SingleField(t *testing.T) {
	fc := NewFieldCounter("level")
	fc.Add(`ts=2024-01-01 level=info msg="hello world"`)
	fc.Add(`ts=2024-01-01 level=info msg="second"`)
	fc.Add(`ts=2024-01-01 level=error msg="boom"`)

	counts := fc.Counts()
	if counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", counts["info"])
	}
	if counts["error"] != 1 {
		t.Errorf("expected error=1, got %d", counts["error"])
	}
}

func TestFieldCounter_Add_QuotedValue(t *testing.T) {
	fc := NewFieldCounter("user")
	fc.Add(`level=info user="alice" action=login`)

	if got := fc.Counts()["alice"]; got != 1 {
		t.Errorf("expected alice=1, got %d", got)
	}
}

func TestFieldCounter_Add_MissingField_Skipped(t *testing.T) {
	fc := NewFieldCounter("level")
	fc.Add("no fields here")
	fc.Add("other=value")

	if got := fc.Counts(); len(got) != 0 {
		t.Errorf("expected empty counts, got %v", got)
	}
}

func TestSortedFieldEntries_Order(t *testing.T) {
	counts := map[string]int{"warn": 1, "info": 5, "error": 3}
	entries := SortedFieldEntries(counts)

	if entries[0].Value != "info" || entries[0].Count != 5 {
		t.Errorf("expected first entry info=5, got %s=%d", entries[0].Value, entries[0].Count)
	}
	if entries[1].Value != "error" || entries[1].Count != 3 {
		t.Errorf("expected second entry error=3, got %s=%d", entries[1].Value, entries[1].Count)
	}
}

func TestSortedFieldEntries_TieBreakByValue(t *testing.T) {
	counts := map[string]int{"warn": 2, "info": 2}
	entries := SortedFieldEntries(counts)

	if entries[0].Value != "info" {
		t.Errorf("expected info first (tie-break), got %s", entries[0].Value)
	}
}

func TestCountFieldReader(t *testing.T) {
	input := strings.NewReader(
		"level=info msg=a\n" +
			"level=warn msg=b\n" +
			"level=info msg=c\n" +
			"msg=no-level\n",
	)
	entries, err := CountFieldReader(input, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Value != "info" || entries[0].Count != 2 {
		t.Errorf("expected info=2, got %s=%d", entries[0].Value, entries[0].Count)
	}
}
