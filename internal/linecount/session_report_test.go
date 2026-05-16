package linecount

import (
	"strings"
	"testing"
	"time"
)

func TestWriteSessionReport_Empty(t *testing.T) {
	var sb strings.Builder
	if err := WriteSessionReport(&sb, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "no sessions") {
		t.Errorf("expected 'no sessions' message, got: %q", sb.String())
	}
}

func TestWriteSessionReport_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	sessions := []SessionEntry{
		{Start: sessionTime(0), End: sessionTime(5 * time.Minute), LineCount: 10},
	}
	if err := WriteSessionReport(&sb, sessions); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, hdr := range []string{"Start", "End", "Duration", "Lines"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output: %s", hdr, out)
		}
	}
}

func TestWriteSessionReport_ContainsTotals(t *testing.T) {
	var sb strings.Builder
	sessions := []SessionEntry{
		{Start: sessionTime(0), End: sessionTime(5 * time.Minute), LineCount: 7},
		{Start: sessionTime(60 * time.Minute), End: sessionTime(65 * time.Minute), LineCount: 3},
	}
	if err := WriteSessionReport(&sb, sessions); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "Total sessions: 2") {
		t.Errorf("expected total sessions in output: %s", out)
	}
	if !strings.Contains(out, "Total lines: 10") {
		t.Errorf("expected total lines in output: %s", out)
	}
}

func TestFormatSessionDuration_Seconds(t *testing.T) {
	got := formatSessionDuration(45 * time.Second)
	if got != "45s" {
		t.Errorf("expected '45s', got %q", got)
	}
}

func TestFormatSessionDuration_Minutes(t *testing.T) {
	got := formatSessionDuration(5*time.Minute + 30*time.Second)
	if got != "5m30s" {
		t.Errorf("expected '5m30s', got %q", got)
	}
}

func TestFormatSessionDuration_Hours(t *testing.T) {
	got := formatSessionDuration(2*time.Hour + 15*time.Minute)
	if got != "2h15m" {
		t.Errorf("expected '2h15m', got %q", got)
	}
}
