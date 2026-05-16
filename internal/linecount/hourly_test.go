package linecount

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func hourTime(y, mo, d, h int) time.Time {
	return time.Date(y, time.Month(mo), d, h, 0, 0, 0, time.UTC)
}

func TestHourlyCounter_Empty(t *testing.T) {
	c := NewHourlyCounter()
	if got := c.Entries(); len(got) != 0 {
		t.Fatalf("expected empty, got %d entries", len(got))
	}
}

func TestHourlyCounter_Add_SingleHour(t *testing.T) {
	c := NewHourlyCounter()
	base := hourTime(2024, 3, 10, 14)
	c.Add(base)
	c.Add(base.Add(30 * time.Minute))

	entries := c.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(entries))
	}
	if entries[0].Count != 2 {
		t.Errorf("expected count 2, got %d", entries[0].Count)
	}
}

func TestHourlyCounter_Add_MultipleHours_SortedChronologically(t *testing.T) {
	c := NewHourlyCounter()
	c.Add(hourTime(2024, 3, 10, 16))
	c.Add(hourTime(2024, 3, 10, 14))
	c.Add(hourTime(2024, 3, 10, 15))

	entries := c.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if !entries[i].Hour.After(entries[i-1].Hour) {
			t.Errorf("entries not sorted at index %d", i)
		}
	}
}

func TestHourlyCounter_Add_AccumulatesSameHour(t *testing.T) {
	c := NewHourlyCounter()
	h := hourTime(2024, 3, 10, 9)
	for i := 0; i < 5; i++ {
		c.Add(h.Add(time.Duration(i) * time.Minute))
	}
	entries := c.Entries()
	if entries[0].Count != 5 {
		t.Errorf("expected 5, got %d", entries[0].Count)
	}
}

func TestWriteHourlyReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteHourlyReport(&buf, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No data") {
		t.Errorf("expected 'No data', got %q", buf.String())
	}
}

func TestWriteHourlyReport_ContainsHour(t *testing.T) {
	c := NewHourlyCounter()
	c.Add(hourTime(2024, 3, 10, 11))
	c.Add(hourTime(2024, 3, 10, 11))

	var buf bytes.Buffer
	if err := WriteHourlyReport(&buf, c.Entries()); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "2024-03-10 11:00") {
		t.Errorf("expected hour label in output, got:\n%s", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected count 2 in output, got:\n%s", out)
	}
}
