package linecount

import (
	"bytes"
	"strings"
	"testing"
)

func buildCounter(t *testing.T, entries []TopNEntry) *WindowCounter {
	t.Helper()
	wc := NewWindowCounter()
	for _, e := range entries {
		for i := int64(0); i < e.Lines; i++ {
			wc.Add(e.Key, e.Bytes/e.Lines)
		}
	}
	return wc
}

func TestTopN_Empty(t *testing.T) {
	wc := NewWindowCounter()
	result := TopN(wc, 3)
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(result))
	}
}

func TestTopN_SortedDescending(t *testing.T) {
	wc := NewWindowCounter()
	wc.Add("2024-01-01 00:00", 10)
	wc.Add("2024-01-01 00:00", 10)
	wc.Add("2024-01-01 01:00", 20)
	wc.Add("2024-01-01 02:00", 30)
	wc.Add("2024-01-01 02:00", 30)
	wc.Add("2024-01-01 02:00", 30)

	result := TopN(wc, 0)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Lines < result[1].Lines || result[1].Lines < result[2].Lines {
		t.Errorf("entries not sorted descending: %+v", result)
	}
}

func TestTopN_LimitN(t *testing.T) {
	wc := NewWindowCounter()
	for _, key := range []string{"a", "b", "c", "d", "e"} {
		wc.Add(key, 5)
	}
	result := TopN(wc, 3)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestTopN_NLargerThanEntries(t *testing.T) {
	wc := NewWindowCounter()
	wc.Add("only", 100)
	result := TopN(wc, 10)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
}

func TestWriteTopNReport_Empty(t *testing.T) {
	wc := NewWindowCounter()
	var buf bytes.Buffer
	if err := WriteTopNReport(&buf, wc, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data' message, got: %q", buf.String())
	}
}

func TestWriteTopNReport_ContainsHeaders(t *testing.T) {
	wc := NewWindowCounter()
	wc.Add("2024-01-01 00:00", 512)
	wc.Add("2024-01-01 00:00", 512)

	var buf bytes.Buffer
	if err := WriteTopNReport(&buf, wc, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"RANK", "WINDOW", "LINES", "BYTES"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output:\n%s", hdr, out)
		}
	}
}

func TestWriteTopNReport_RankOrder(t *testing.T) {
	wc := NewWindowCounter()
	wc.Add("low", 10)
	wc.Add("high", 20)
	wc.Add("high", 20)
	wc.Add("mid", 15)

	var buf bytes.Buffer
	if err := WriteTopNReport(&buf, wc, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	highPos := strings.Index(out, "high")
	midPos := strings.Index(out, "mid")
	lowPos := strings.Index(out, "low")
	if highPos > midPos || midPos > lowPos {
		t.Errorf("unexpected rank order in output:\n%s", out)
	}
}
