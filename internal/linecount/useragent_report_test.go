package linecount

import (
	"strings"
	"testing"
)

func TestWriteUserAgentReport_Empty(t *testing.T) {
	var sb strings.Builder
	WriteUserAgentReport(&sb, nil, 0)
	if !strings.Contains(sb.String(), "No user-agent") {
		t.Errorf("expected empty message, got: %s", sb.String())
	}
}

func TestWriteUserAgentReport_ContainsHeaders(t *testing.T) {
	entries := []UserAgentEntry{{Agent: "Mozilla/5.0", Count: 10}}
	var sb strings.Builder
	WriteUserAgentReport(&sb, entries, 10)
	out := sb.String()
	if !strings.Contains(out, "User-Agent") {
		t.Errorf("expected header 'User-Agent', got: %s", out)
	}
	if !strings.Contains(out, "Count") {
		t.Errorf("expected header 'Count', got: %s", out)
	}
	if !strings.Contains(out, "Percent") {
		t.Errorf("expected header 'Percent', got: %s", out)
	}
}

func TestWriteUserAgentReport_PercentageCalculation(t *testing.T) {
	entries := []UserAgentEntry{
		{Agent: "Mozilla/5.0", Count: 75},
		{Agent: "curl/7.68", Count: 25},
	}
	var sb strings.Builder
	WriteUserAgentReport(&sb, entries, 100)
	out := sb.String()
	if !strings.Contains(out, "75.0%") {
		t.Errorf("expected 75.0%% in output, got: %s", out)
	}
	if !strings.Contains(out, "25.0%") {
		t.Errorf("expected 25.0%% in output, got: %s", out)
	}
}

func TestWriteUserAgentReport_TotalLine(t *testing.T) {
	entries := []UserAgentEntry{{Agent: "wget/1.21", Count: 5}}
	var sb strings.Builder
	WriteUserAgentReport(&sb, entries, 5)
	out := sb.String()
	if !strings.Contains(out, "Total") {
		t.Errorf("expected Total line, got: %s", out)
	}
}

func TestWriteUserAgentReport_ZeroTotal_NoPanic(t *testing.T) {
	entries := []UserAgentEntry{{Agent: "bot", Count: 3}}
	var sb strings.Builder
	WriteUserAgentReport(&sb, entries, 0)
	// should not panic; percentage should be 0
	if !strings.Contains(sb.String(), "0.0%") {
		t.Errorf("expected 0.0%% when total is zero, got: %s", sb.String())
	}
}
