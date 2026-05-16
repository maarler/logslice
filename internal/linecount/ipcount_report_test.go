package linecount

import (
	"strings"
	"testing"
)

func TestWriteIPReport_Empty(t *testing.T) {
	var sb strings.Builder
	if err := WriteIPReport(&sb, nil, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No IP") {
		t.Fatalf("expected empty message, got: %q", sb.String())
	}
}

func TestWriteIPReport_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	entries := []IPEntry{{IP: "10.0.0.1", Count: 5}}
	if err := WriteIPReport(&sb, entries, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, hdr := range []string{"IP ADDRESS", "COUNT", "PERCENT"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestWriteIPReport_PercentageCalculation(t *testing.T) {
	var sb strings.Builder
	entries := []IPEntry{
		{IP: "192.168.1.1", Count: 1},
		{IP: "192.168.1.2", Count: 3},
	}
	if err := WriteIPReport(&sb, entries, 4); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "25.00%") {
		t.Errorf("expected 25.00%% in output, got: %q", out)
	}
	if !strings.Contains(out, "75.00%") {
		t.Errorf("expected 75.00%% in output, got: %q", out)
	}
}

func TestWriteIPReport_TotalLine(t *testing.T) {
	var sb strings.Builder
	entries := []IPEntry{{IP: "10.1.1.1", Count: 42}}
	_ = WriteIPReport(&sb, entries, 42)
	if !strings.Contains(sb.String(), "Total hits: 42") {
		t.Errorf("expected total line in output")
	}
}

func TestWriteIPReport_ZeroTotal_NoPanic(t *testing.T) {
	var sb strings.Builder
	entries := []IPEntry{{IP: "10.0.0.1", Count: 0}}
	if err := WriteIPReport(&sb, entries, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "0.00%") {
		t.Errorf("expected 0.00%% when total is zero")
	}
}
