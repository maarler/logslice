package linecount_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/linecount"
)

func TestRunLatencyCount_EmptyReader(t *testing.T) {
	var out bytes.Buffer
	err := linecount.RunLatencyCount(&out, strings.NewReader(""), linecount.DefaultLatencyCmdOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "Latency") {
		t.Errorf("expected header in output, got: %q", out.String())
	}
}

func TestRunLatencyCount_SkipsUnparseable(t *testing.T) {
	input := strings.NewReader("no latency here\nanother bad line\n")
	var out bytes.Buffer
	opts := linecount.DefaultLatencyCmdOptions()
	err := linecount.RunLatencyCount(&out, input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out.String(), "NaN") {
		t.Errorf("output should not contain NaN: %q", out.String())
	}
}

func TestRunLatencyCount_ParsesValues(t *testing.T) {
	input := strings.NewReader(
		`duration=12ms status=200
duration=50ms status=200
duration=200ms status=500
`,
	)
	var out bytes.Buffer
	opts := linecount.DefaultLatencyCmdOptions()
	opts.Field = "duration"
	err := linecount.RunLatencyCount(&out, input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := out.String()
	if !strings.Contains(result, "p50") {
		t.Errorf("expected p50 in output, got: %q", result)
	}
	if !strings.Contains(result, "p99") {
		t.Errorf("expected p99 in output, got: %q", result)
	}
}

func TestRunLatencyCount_CustomField(t *testing.T) {
	input := strings.NewReader(
		`elapsed=5ms\nelapsed=15ms\nelapsed=25ms\n`,
	)
	var out bytes.Buffer
	opts := linecount.DefaultLatencyCmdOptions()
	opts.Field = "elapsed"
	err := linecount.RunLatencyCount(&out, input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunLatencyCount_HeaderPresent(t *testing.T) {
	var out bytes.Buffer
	err := linecount.RunLatencyCount(&out, strings.NewReader("duration=10ms\n"), linecount.DefaultLatencyCmdOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := out.String()
	if !strings.Contains(result, "min") {
		t.Errorf("expected 'min' header in output, got: %q", result)
	}
	if !strings.Contains(result, "max") {
		t.Errorf("expected 'max' header in output, got: %q", result)
	}
}
