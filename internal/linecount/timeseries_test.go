package linecount

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseTime(offset int) time.Time {
	return time.Date(2024, 1, 1, 0, offset, 0, 0, time.UTC)
}

func TestTimeSeriesCounter_Empty(t *testing.T) {
	c := NewTimeSeriesCounter()
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
	if entries := c.Entries(); len(entries) != 0 {
		t.Fatalf("expected empty entries, got %d", len(entries))
	}
}

func TestTimeSeriesCounter_Add_SingleBucket(t *testing.T) {
	c := NewTimeSeriesCounter()
	c.Add("2024-01-01T00:00", baseTime(0), 5)
	if got := c.Total(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestTimeSeriesCounter_Add_AccumulatesSameBucket(t *testing.T) {
	c := NewTimeSeriesCounter()
	key := "2024-01-01T00:00"
	c.Add(key, baseTime(0), 3)
	c.Add(key, baseTime(0), 7)
	if got := c.Total(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
	if got := len(c.Entries()); got != 1 {
		t.Fatalf("expected 1 entry, got %d", got)
	}
}

func TestTimeSeriesCounter_Entries_SortedChronologically(t *testing.T) {
	c := NewTimeSeriesCounter()
	c.Add("2024-01-01T00:02", baseTime(2), 1)
	c.Add("2024-01-01T00:00", baseTime(0), 1)
	c.Add("2024-01-01T00:01", baseTime(1), 1)

	entries := c.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Time.Before(entries[i-1].Time) {
			t.Errorf("entries not sorted at index %d", i)
		}
	}
}

func TestWriteTimeSeriesReport_Empty(t *testing.T) {
	c := NewTimeSeriesCounter()
	var buf bytes.Buffer
	if err := WriteTimeSeriesReport(&buf, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data' in output, got: %s", buf.String())
	}
}

func TestWriteTimeSeriesReport_ContainsHeaders(t *testing.T) {
	c := NewTimeSeriesCounter()
	c.Add("2024-01-01T00:00", baseTime(0), 42)

	var buf bytes.Buffer
	if err := WriteTimeSeriesReport(&buf, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Time", "Lines", "TOTAL", "42"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output\n%s", want, out)
		}
	}
}

func TestWriteTimeSeriesReport_TotalMatchesSum(t *testing.T) {
	c := NewTimeSeriesCounter()
	c.Add("2024-01-01T00:00", baseTime(0), 10)
	c.Add("2024-01-01T00:01", baseTime(1), 20)
	c.Add("2024-01-01T00:02", baseTime(2), 30)

	if got := c.Total(); got != 60 {
		t.Fatalf("expected total 60, got %d", got)
	}
}
