package linecount

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRunHistogram_EmptyReader(t *testing.T) {
	var buf bytes.Buffer
	err := RunHistogram(strings.NewReader(""), &buf, DefaultHistogramOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data', got: %s", buf.String())
	}
}

func TestRunHistogram_SkipsUnparseable(t *testing.T) {
	input := "not a log line\nalso not a log line\n"
	var buf bytes.Buffer
	err := RunHistogram(strings.NewReader(input), &buf, DefaultHistogramOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "skipped 2") {
		t.Errorf("expected skip notice, got: %s", buf.String())
	}
}

func TestRunHistogram_ParsesTimestamps(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-15T10:05:01Z INFO starting server",
		"2024-01-15T10:05:30Z INFO accepting connections",
		"2024-01-15T10:06:01Z WARN high memory",
	}, "\n") + "\n"

	opts := DefaultHistogramOptions()
	opts.Window = time.Minute
	var buf bytes.Buffer
	err := RunHistogram(strings.NewReader(input), &buf, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "10:05") {
		t.Errorf("expected 10:05 bucket, got:\n%s", out)
	}
	if !strings.Contains(out, "10:06") {
		t.Errorf("expected 10:06 bucket, got:\n%s", out)
	}
}

func TestRunHistogram_DefaultOptions_ZeroWindow(t *testing.T) {
	opts := HistogramOptions{Window: 0, BarWidth: 0}
	input := "2024-01-15T10:05:01Z INFO msg\n"
	var buf bytes.Buffer
	if err := RunHistogram(strings.NewReader(input), &buf, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunHistogram_HourlyBuckets(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-15T09:10:00Z INFO a",
		"2024-01-15T09:50:00Z INFO b",
		"2024-01-15T10:05:00Z INFO c",
	}, "\n") + "\n"

	opts := DefaultHistogramOptions()
	opts.Window = time.Hour
	var buf bytes.Buffer
	if err := RunHistogram(strings.NewReader(input), &buf, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	// hourly format: "2006-01-02 15h"
	if !strings.Contains(out, "09h") {
		t.Errorf("expected 09h bucket, got:\n%s", out)
	}
	if !strings.Contains(out, "10h") {
		t.Errorf("expected 10h bucket, got:\n%s", out)
	}
}
