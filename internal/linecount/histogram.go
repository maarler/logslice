package linecount

import (
	"fmt"
	"io"
	"math"
	"sort"
	"time"
)

// HistogramBucket holds a time bucket label and its line count.
type HistogramBucket struct {
	Label string
	Count int64
	Bytes int64
}

// TimeHistogram aggregates log line counts into fixed-width time buckets.
type TimeHistogram struct {
	buckets  map[string]*HistogramBucket
	ordered  []string
	window   time.Duration
	format   string
}

// NewTimeHistogram creates a histogram with the given bucket window.
func NewTimeHistogram(window time.Duration) *TimeHistogram {
	format := bucketFormat(window)
	return &TimeHistogram{
		buckets: make(map[string]*HistogramBucket),
		window:  window,
		format:  format,
	}
}

func bucketFormat(d time.Duration) string {
	switch {
	case d >= 24*time.Hour:
		return "2006-01-02"
	case d >= time.Hour:
		return "2006-01-02 15h"
	default:
		return "2006-01-02 15:04"
	}
}

// Add records a line with its timestamp and byte size into the appropriate bucket.
func (h *TimeHistogram) Add(t time.Time, bytes int64) {
	truncated := t.Truncate(h.window)
	label := truncated.Format(h.format)
	if _, ok := h.buckets[label]; !ok {
		h.buckets[label] = &HistogramBucket{Label: label}
		h.ordered = append(h.ordered, label)
	}
	h.buckets[label].Count++
	h.buckets[label].Bytes += bytes
}

// Buckets returns buckets sorted chronologically.
func (h *TimeHistogram) Buckets() []HistogramBucket {
	sorted := make([]string, len(h.ordered))
	copy(sorted, h.ordered)
	sort.Strings(sorted)
	out := make([]HistogramBucket, 0, len(sorted))
	for _, k := range sorted {
		out = append(out, *h.buckets[k])
	}
	return out
}

// WriteHistogramReport writes an ASCII histogram to w.
func WriteHistogramReport(w io.Writer, h *TimeHistogram, barWidth int) {
	buckets := h.Buckets()
	if len(buckets) == 0 {
		fmt.Fprintln(w, "(no data)")
		return
	}
	var maxCount int64
	for _, b := range buckets {
		if b.Count > maxCount {
			maxCount = b.Count
		}
	}
	fmt.Fprintf(w, "%-20s %8s  %s\n", "bucket", "lines", "distribution")
	fmt.Fprintf(w, "%-20s %8s  %s\n", "------", "-----", "------------")
	for _, b := range buckets {
		width := 0
		if maxCount > 0 {
			width = int(math.Round(float64(b.Count) / float64(maxCount) * float64(barWidth)))
		}
		bar := ""
		for i := 0; i < width; i++ {
			bar += "█"
		}
		fmt.Fprintf(w, "%-20s %8d  %s\n", b.Label, b.Count, bar)
	}
}
