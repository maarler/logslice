package linecount

import (
	"strings"
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestWindowCounter_Empty(t *testing.T) {
	wc := NewWindowCounter()
	if got := wc.Stats(); len(got) != 0 {
		t.Fatalf("expected empty stats, got %d entries", len(got))
	}
}

func TestWindowCounter_Add_SingleWindow(t *testing.T) {
	wc := NewWindowCounter()
	start := epoch
	end := epoch.Add(time.Minute)
	wc.Add("2024-01-01T00:00", start, end, "hello world")
	wc.Add("2024-01-01T00:00", start, end, "foo bar")

	stats := wc.Stats()
	if len(stats) != 1 {
		t.Fatalf("expected 1 window, got %d", len(stats))
	}
	if stats[0].Lines != 2 {
		t.Errorf("expected 2 lines, got %d", stats[0].Lines)
	}
	expectedBytes := int64(len("hello world") + len("foo bar"))
	if stats[0].Bytes != expectedBytes {
		t.Errorf("expected %d bytes, got %d", expectedBytes, stats[0].Bytes)
	}
}

func TestWindowCounter_Add_MultipleWindows_PreservesOrder(t *testing.T) {
	wc := NewWindowCounter()
	keys := []string{"w1", "w2", "w3"}
	for _, k := range keys {
		wc.Add(k, epoch, epoch.Add(time.Minute), "line")
	}
	stats := wc.Stats()
	if len(stats) != 3 {
		t.Fatalf("expected 3 windows, got %d", len(stats))
	}
	for i, s := range stats {
		if s.Key != keys[i] {
			t.Errorf("position %d: expected key %q, got %q", i, keys[i], s.Key)
		}
	}
}

func TestCountWindowReader_GroupsByKey(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T00:00:01 alpha",
		"2024-01-01T00:00:45 beta",
		"2024-01-01T00:01:10 gamma",
	}, "\n")

	// key function: use first 16 chars as minute bucket
	keyFn := func(line string) (string, time.Time, time.Time) {
		if len(line) < 16 {
			return "", time.Time{}, time.Time{}
		}
		key := line[:16]
		return key, epoch, epoch.Add(time.Minute)
	}

	wc := CountWindowReader(strings.NewReader(input), keyFn)
	stats := wc.Stats()
	if len(stats) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(stats))
	}
	if stats[0].Lines != 2 {
		t.Errorf("first window: expected 2 lines, got %d", stats[0].Lines)
	}
	if stats[1].Lines != 1 {
		t.Errorf("second window: expected 1 line, got %d", stats[1].Lines)
	}
}

func TestCountWindowReader_SkipsEmptyKey(t *testing.T) {
	input := "no-timestamp line\n2024-01-01T00:00:01 valid"
	keyFn := func(line string) (string, time.Time, time.Time) {
		if len(line) < 20 {
			return "", time.Time{}, time.Time{}
		}
		return line[:16], epoch, epoch.Add(time.Minute)
	}
	wc := CountWindowReader(strings.NewReader(input), keyFn)
	if len(wc.Stats()) != 1 {
		t.Errorf("expected 1 window, got %d", len(wc.Stats()))
	}
}
