package linecount

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// HourlyEntry holds the line count for a single hour bucket.
type HourlyEntry struct {
	Hour  time.Time
	Count int64
}

// HourlyCounter accumulates line counts bucketed by calendar hour.
type HourlyCounter struct {
	buckets map[string]*HourlyEntry
}

// NewHourlyCounter returns an initialised HourlyCounter.
func NewHourlyCounter() *HourlyCounter {
	return &HourlyCounter{buckets: make(map[string]*HourlyEntry)}
}

func hourKey(t time.Time) string {
	return t.UTC().Format("2006-01-02T15")
}

// Add records one line observed at time t.
func (c *HourlyCounter) Add(t time.Time) {
	k := hourKey(t)
	if _, ok := c.buckets[k]; !ok {
		c.buckets[k] = &HourlyEntry{
			Hour: t.UTC().Truncate(time.Hour),
		}
	}
	c.buckets[k].Count++
}

// Entries returns all buckets sorted chronologically.
func (c *HourlyCounter) Entries() []HourlyEntry {
	out := make([]HourlyEntry, 0, len(c.buckets))
	for _, e := range c.buckets {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Hour.Before(out[j].Hour)
	})
	return out
}

// WriteHourlyReport writes a human-readable hourly breakdown to w.
func WriteHourlyReport(w io.Writer, entries []HourlyEntry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No data.")
		return err
	}

	var max int64
	for _, e := range entries {
		if e.Count > max {
			max = e.Count
		}
	}

	_, err := fmt.Fprintf(w, "%-20s %10s  %s\n", "Hour (UTC)", "Lines", "Bar")
	if err != nil {
		return err
	}
	for _, e := range entries {
		bar := buildBar(e.Count, max, 30)
		_, err = fmt.Fprintf(w, "%-20s %10d  %s\n",
			e.Hour.Format("2006-01-02 15:00"), e.Count, bar)
		if err != nil {
			return err
		}
	}
	return nil
}
