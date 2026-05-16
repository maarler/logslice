package linecount

import (
	"strings"
	"testing"
)

func TestErrorCounter_Empty(t *testing.T) {
	c := NewErrorCounter()
	if c.Total() != 0 {
		t.Fatalf("expected 0 total, got %d", c.Total())
	}
	if len(c.Counts()) != 0 {
		t.Fatalf("expected empty counts map")
	}
}

func TestErrorCounter_Add_SingleLevel(t *testing.T) {
	c := NewErrorCounter()
	c.Add("2024-01-01 ERROR something went wrong")
	if c.Total() != 1 {
		t.Fatalf("expected total 1, got %d", c.Total())
	}
	if c.Counts()[LevelError] != 1 {
		t.Fatalf("expected ERROR=1")
	}
}

func TestErrorCounter_Add_NormalisesWarning(t *testing.T) {
	c := NewErrorCounter()
	c.Add("WARNING: disk almost full")
	if c.Counts()[LevelWarn] != 1 {
		t.Fatalf("expected WARN=1 after WARNING line")
	}
}

func TestErrorCounter_Add_NormalisesCritical(t *testing.T) {
	c := NewErrorCounter()
	c.Add("CRITICAL system failure")
	if c.Counts()[LevelFatal] != 1 {
		t.Fatalf("expected FATAL=1 after CRITICAL line")
	}
}

func TestErrorCounter_Add_IgnoresUnknown(t *testing.T) {
	c := NewErrorCounter()
	c.Add("plain log line with no level")
	if c.Total() != 0 {
		t.Fatalf("expected total 0 for unrecognised line")
	}
}

func TestErrorCounter_Add_MultipleLevels(t *testing.T) {
	c := NewErrorCounter()
	lines := []string{
		"INFO starting up",
		"DEBUG connecting",
		"ERROR failed",
		"ERROR retry",
		"FATAL abort",
	}
	for _, l := range lines {
		c.Add(l)
	}
	if c.Total() != 5 {
		t.Fatalf("expected total 5, got %d", c.Total())
	}
	if c.Counts()[LevelError] != 2 {
		t.Fatalf("expected ERROR=2")
	}
}

func TestSortedErrorEntries_ByCountDesc(t *testing.T) {
	counts := map[ErrorLevel]int{
		LevelInfo:  10,
		LevelError: 5,
		LevelWarn:  7,
	}
	entries := SortedErrorEntries(counts)
	if entries[0].Level != LevelInfo {
		t.Fatalf("expected INFO first, got %s", entries[0].Level)
	}
	if entries[1].Level != LevelWarn {
		t.Fatalf("expected WARN second, got %s", entries[1].Level)
	}
}

func TestCountErrorReader(t *testing.T) {
	input := strings.NewReader("ERROR line one\nINFO line two\nERROR line three\n")
	c := CountErrorReader(input)
	if c.Total() != 3 {
		t.Fatalf("expected total 3, got %d", c.Total())
	}
	if c.Counts()[LevelError] != 2 {
		t.Fatalf("expected ERROR=2")
	}
}
