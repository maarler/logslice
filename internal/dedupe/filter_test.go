package dedupe

import (
	"bytes"
	"strings"
	"testing"
)

func TestFilter_ApplyLines_RemovesDupes(t *testing.T) {
	f := NewFilter(DefaultOptions())
	input := []string{"a", "b", "a", "c", "b"}
	got := f.ApplyLines(input)
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: want %q got %q", i, want[i], got[i])
		}
	}
}

func TestFilter_ApplyLines_Empty(t *testing.T) {
	f := NewFilter(DefaultOptions())
	got := f.ApplyLines(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestFilter_Apply_Reader(t *testing.T) {
	f := NewFilter(DefaultOptions())
	input := "foo\nbar\nfoo\nbaz\nbar\n"
	r := strings.NewReader(input)
	var w bytes.Buffer
	suppressed, err := f.Apply(r, &w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if suppressed != 2 {
		t.Fatalf("expected 2 suppressed, got %d", suppressed)
	}
	want := "foo\nbar\nbaz\n"
	if w.String() != want {
		t.Fatalf("expected %q, got %q", want, w.String())
	}
}

func TestFilter_Apply_NoTrailingNewline(t *testing.T) {
	f := NewFilter(DefaultOptions())
	r := strings.NewReader("hello\nhello")
	var w bytes.Buffer
	suppressed, err := f.Apply(r, &w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if suppressed != 1 {
		t.Fatalf("expected 1 suppressed, got %d", suppressed)
	}
	if !strings.Contains(w.String(), "hello") {
		t.Fatalf("expected hello in output, got %q", w.String())
	}
}

func TestFilter_Apply_ConsecutiveMode(t *testing.T) {
	f := NewFilter(Options{WindowSize: 8, Consecutive: true})
	input := "x\nx\ny\nx\n"
	r := strings.NewReader(input)
	var w bytes.Buffer
	suppressed, err := f.Apply(r, &w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// second "x" is consecutive dupe; third "x" is not consecutive.
	if suppressed != 1 {
		t.Fatalf("expected 1 suppressed, got %d", suppressed)
	}
	lines := strings.Split(strings.TrimSpace(w.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 output lines, got %d: %v", len(lines), lines)
	}
}
