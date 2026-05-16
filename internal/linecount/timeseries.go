package linecount

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"time"
)

// TimeSeriesEntry holds the line count for a specific time bucket.
type TimeSeriesEntry struct {
	Time  time.Time
	Key   string
	Count int64
}

// TimeSeriesCounter accumulates line counts keyed by a time bucket string.
type TimeSeriesCounter struct {
	entries map[string]*TimeSeriesEntry
	order   []string
}

// NewTimeSeriesCounter creates an empty TimeSeriesCounter.
func NewTimeSeriesCounter() *TimeSeriesCounter {
	return &TimeSeriesCounter{
		entries: make(map[string]*TimeSeriesEntry),
	}
}

// Add increments the count for the given time bucket key.
func (c *TimeSeriesCounter) Add(key string, t time.Time, n int64) {
	if e, ok := c.entries[key]; ok {
		e.Count += n
		return
	}
	c.entries[key] = &TimeSeriesEntry{Time: t, Key: key, Count: n}
	c.order = append(c.order, key)
}

// Entries returns all entries sorted chronologically.
func (c *TimeSeriesCounter) Entries() []TimeSeriesEntry {
	out := make([]TimeSeriesEntry, 0, len(c.entries))
	for _, k := range c.order {
		out = append(out, *c.entries[k])
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Time.Before(out[j].Time)
	})
	return out
}

// Total returns the sum of all counts.
func (c *TimeSeriesCounter) Total() int64 {
	var total int64
	for _, e := range c.entries {
		total += e.Count
	}
	return total
}

// WriteTimeSeriesReport writes a formatted time-series report to w.
func WriteTimeSeriesReport(w io.Writer, c *TimeSeriesCounter) error {
	entries := c.Entries()
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "(no data)")
		return err
	}

	var maxCount int64
	for _, e := range entries {
		if e.Count > maxCount {
			maxCount = e.Count
		}
	}

	bw := bufio.NewWriter(w)
	fmt.Fprintf(bw, "%-24s  %8s  %s\n", "Time", "Lines", "Bar")
	fmt.Fprintf(bw, "%-24s  %8s  %s\n", "------------------------", "--------", "---")
	for _, e := range entries {
		bar := buildBar(e.Count, maxCount, 30)
		fmt.Fprintf(bw, "%-24s  %8d  %s\n", e.Key, e.Count, bar)
	}
	fmt.Fprintf(bw, "%-24s  %8d\n", "TOTAL", c.Total())
	return bw.Flush()
}
