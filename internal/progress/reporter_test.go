package progress

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestReporter_QuietProducesNoOutput(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, true)
	r.Start(10 * time.Millisecond)
	r.Add(100, 2048)
	time.Sleep(30 * time.Millisecond)
	r.Stop()

	if buf.Len() != 0 {
		t.Errorf("expected no output in quiet mode, got %q", buf.String())
	}
}

func TestReporter_Add_Accumulates(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, true)
	r.Add(10, 100)
	r.Add(5, 50)

	if got := r.Lines(); got != 15 {
		t.Errorf("Lines() = %d, want 15", got)
	}
	if got := r.Bytes(); got != 150 {
		t.Errorf("Bytes() = %d, want 150", got)
	}
}

func TestReporter_PrintsProgress(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, false)
	r.Add(42, 1024)
	r.Stop()

	out := buf.String()
	if !strings.Contains(out, "42 lines") {
		t.Errorf("expected line count in output, got %q", out)
	}
	if !strings.Contains(out, "1.0 KiB") {
		t.Errorf("expected byte size in output, got %q", out)
	}
}

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		input int64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KiB"},
		{1536, "1.5 KiB"},
		{1048576, "1.0 MiB"},
		{1073741824, "1.0 GiB"},
	}
	for _, tc := range cases {
		got := formatBytes(tc.input)
		if got != tc.want {
			t.Errorf("formatBytes(%d) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestReporter_TickerFires(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, false)
	r.Start(10 * time.Millisecond)
	r.Add(7, 700)
	time.Sleep(35 * time.Millisecond)
	r.Stop()

	out := buf.String()
	if !strings.Contains(out, "7 lines") {
		t.Errorf("expected ticker output to contain line count, got %q", out)
	}
}

// TestReporter_Add_ZeroValues verifies that adding zero lines and bytes
// does not change the reporter's accumulated totals.
func TestReporter_Add_ZeroValues(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, true)
	r.Add(10, 100)
	r.Add(0, 0)

	if got := r.Lines(); got != 10 {
		t.Errorf("Lines() = %d, want 10", got)
	}
	if got := r.Bytes(); got != 100 {
		t.Errorf("Bytes() = %d, want 100", got)
	}
}
