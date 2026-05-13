package highlight

import (
	"strings"
	"testing"
)

func TestNewRule_Valid(t *testing.T) {
	r, err := NewRule(`\d+`, Yellow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.re == nil {
		t.Fatal("expected compiled regexp, got nil")
	}
}

func TestNewRule_Invalid(t *testing.T) {
	_, err := NewRule(`[invalid`, Red)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestHighlighter_NoRules_Unchanged(t *testing.T) {
	h := New(nil)
	line := "2024-01-01 ERROR something went wrong"
	if got := h.Line(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestHighlighter_Line_MatchesPattern(t *testing.T) {
	r, _ := NewRule("ERROR", Red)
	h := New([]Rule{r})

	line := "2024-01-01 ERROR something went wrong"
	got := h.Line(line)

	if !strings.Contains(got, Red+"ERROR"+Reset) {
		t.Errorf("expected highlighted ERROR in %q", got)
	}
	if !strings.Contains(got, "something went wrong") {
		t.Errorf("expected rest of line preserved in %q", got)
	}
}

func TestHighlighter_Line_MultipleRules(t *testing.T) {
	r1, _ := NewRule("ERROR", Red)
	r2, _ := NewRule(`\d{4}-\d{2}-\d{2}`, Cyan)
	h := New([]Rule{r1, r2})

	got := h.Line("2024-01-01 ERROR msg")

	if !strings.Contains(got, Red+"ERROR"+Reset) {
		t.Errorf("ERROR not highlighted in %q", got)
	}
	if !strings.Contains(got, Cyan) {
		t.Errorf("date not highlighted in %q", got)
	}
}

func TestHighlighter_Lines(t *testing.T) {
	r, _ := NewRule("WARN", Yellow)
	h := New([]Rule{r})

	input := []string{"INFO ok", "WARN slow", "INFO done"}
	out := h.Lines(input)

	if len(out) != len(input) {
		t.Fatalf("expected %d lines, got %d", len(input), len(out))
	}
	if strings.Contains(out[0], Yellow) {
		t.Errorf("line 0 should not be highlighted")
	}
	if !strings.Contains(out[1], Yellow+"WARN"+Reset) {
		t.Errorf("line 1 should be highlighted, got %q", out[1])
	}
}

func TestStripANSI(t *testing.T) {
	input := Red + "ERROR" + Reset + " plain text"
	got := StripANSI(input)
	want := "ERROR plain text"
	if got != want {
		t.Errorf("StripANSI: want %q, got %q", want, got)
	}
}

func TestStripANSI_NoEscapes(t *testing.T) {
	s := "plain log line"
	if got := StripANSI(s); got != s {
		t.Errorf("expected unchanged string, got %q", got)
	}
}
