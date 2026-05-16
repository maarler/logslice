package linecount

import (
	"strings"
	"testing"
)

func TestWriteFieldReport_Empty(t *testing.T) {
	var sb strings.Builder
	err := WriteFieldReport(&sb, "level", nil, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No values found") {
		t.Errorf("expected empty message, got: %q", sb.String())
	}
}

func TestWriteFieldReport_ContainsHeaders(t *testing.T) {
	entries := []FieldEntry{{Value: "info", Count: 10}}
	var sb strings.Builder
	if err := WriteFieldReport(&sb, "level", entries, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, hdr := range []string{"VALUE", "COUNT", "PERCENT", "BAR"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output:\n%s", hdr, out)
		}
	}
}

func TestWriteFieldReport_PercentageCalculation(t *testing.T) {
	entries := []FieldEntry{
		{Value: "info", Count: 3},
		{Value: "error", Count: 1},
	}
	var sb strings.Builder
	if err := WriteFieldReport(&sb, "level", entries, 4); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "75.0%") {
		t.Errorf("expected 75.0%% for info, output:\n%s", out)
	}
	if !strings.Contains(out, "25.0%") {
		t.Errorf("expected 25.0%% for error, output:\n%s", out)
	}
}

func TestWriteFieldReport_ContainsValues(t *testing.T) {
	entries := []FieldEntry{
		{Value: "warn", Count: 7},
		{Value: "debug", Count: 2},
	}
	var sb strings.Builder
	if err := WriteFieldReport(&sb, "level", entries, 9); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "warn") {
		t.Errorf("expected 'warn' in output:\n%s", out)
	}
	if !strings.Contains(out, "debug") {
		t.Errorf("expected 'debug' in output:\n%s", out)
	}
}

func TestWriteFieldReport_ZeroTotal_NoPanic(t *testing.T) {
	entries := []FieldEntry{{Value: "info", Count: 5}}
	var sb strings.Builder
	if err := WriteFieldReport(&sb, "level", entries, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "0.0%") {
		t.Errorf("expected 0.0%% when total is zero, got:\n%s", sb.String())
	}
}
