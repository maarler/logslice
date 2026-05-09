package truncate

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_Defaults(t *testing.T) {
	tr := New(Options{})
	if tr.maxBytes != DefaultMaxBytes {
		t.Errorf("expected maxBytes %d, got %d", DefaultMaxBytes, tr.maxBytes)
	}
	if string(tr.suffix) != DefaultSuffix {
		t.Errorf("expected suffix %q, got %q", DefaultSuffix, string(tr.suffix))
	}
}

func TestLine_ShortLine_Unchanged(t *testing.T) {
	tr := New(Options{MaxBytes: 20, Suffix: "..."})
	input := []byte("short line")
	got := tr.Line(input)
	if string(got) != "short line" {
		t.Errorf("expected %q, got %q", "short line", got)
	}
}

func TestLine_ExactLength_Unchanged(t *testing.T) {
	tr := New(Options{MaxBytes: 10, Suffix: "..."})
	input := []byte("1234567890") // exactly 10 bytes
	got := tr.Line(input)
	if string(got) != "1234567890" {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestLine_LongLine_Truncated(t *testing.T) {
	tr := New(Options{MaxBytes: 10, Suffix: "..."})
	input := []byte("this is a very long line")
	got := tr.Line(input)
	if len(got) > 10 {
		t.Errorf("expected len <= 10, got %d", len(got))
	}
	if !strings.HasSuffix(string(got), "...") {
		t.Errorf("expected suffix '...', got %q", got)
	}
}

func TestLine_DoesNotShareMemory(t *testing.T) {
	tr := New(Options{MaxBytes: 20, Suffix: "..."})
	input := []byte("hello world")
	got := tr.Line(input)
	input[0] = 'X'
	if got[0] == 'X' {
		t.Error("returned slice shares memory with input")
	}
}

func TestApplyLines(t *testing.T) {
	tr := New(Options{MaxBytes: 10, Suffix: "..."})
	lines := []string{
		"short",
		"this is definitely longer than ten bytes",
		"ok",
	}
	result := tr.ApplyLines(lines)
	if len(result) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result))
	}
	if result[0] != "short" {
		t.Errorf("line 0: expected %q, got %q", "short", result[0])
	}
	if len(result[1]) > 10 {
		t.Errorf("line 1: expected len <= 10, got %d", len(result[1]))
	}
	if result[2] != "ok" {
		t.Errorf("line 2: expected %q, got %q", "ok", result[2])
	}
}

func TestApply_Reader(t *testing.T) {
	tr := New(Options{MaxBytes: 10, Suffix: "..."})
	input := readerFromLines([]string{"hello", "this line is way too long for the limit", "bye"})
	var out bytes.Buffer
	if err := tr.Apply(input, &out); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	outLines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(outLines) != 3 {
		t.Fatalf("expected 3 output lines, got %d", len(outLines))
	}
	if len(outLines[1]) > 10 {
		t.Errorf("line 1 not truncated: len=%d, value=%q", len(outLines[1]), outLines[1])
	}
}

func TestCountTruncated(t *testing.T) {
	tr := New(Options{MaxBytes: 10, Suffix: "..."})
	lines := []string{"short", "this is longer than ten bytes", "hi", "another long line here"}
	got := tr.CountTruncated(lines)
	if got != 2 {
		t.Errorf("expected 2 truncated, got %d", got)
	}
}
