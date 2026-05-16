package linecount

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
)

// LatencyStats holds aggregated latency measurements.
type LatencyStats struct {
	Count  int64
	Sum    float64
	Min    float64
	Max    float64
	values []float64
}

// LatencyCounter extracts and aggregates latency values from log lines.
type LatencyCounter struct {
	field string
	sep   string
	kvSep string
	stats LatencyStats
}

// NewLatencyCounter creates a LatencyCounter that reads the given field name.
func NewLatencyCounter(field, sep, kvSep string) *LatencyCounter {
	if sep == "" {
		sep = " "
	}
	if kvSep == "" {
		kvSep = "="
	}
	return &LatencyCounter{
		field: field,
		sep:   sep,
		kvSep: kvSep,
		stats: LatencyStats{Min: math.MaxFloat64},
	}
}

// Add parses a single log line and records the latency value if found.
func (c *LatencyCounter) Add(line string) {
	for _, part := range strings.Split(line, c.sep) {
		kv := strings.SplitN(part, c.kvSep, 2)
		if len(kv) != 2 || strings.TrimSpace(kv[0]) != c.field {
			continue
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(kv[1]), 64)
		if err != nil {
			continue
		}
		c.stats.Count++
		c.stats.Sum += v
		c.stats.values = append(c.stats.values, v)
		if v < c.stats.Min {
			c.stats.Min = v
		}
		if v > c.stats.Max {
			c.stats.Max = v
		}
		return
	}
}

// Percentile returns the p-th percentile (0–100) of observed values.
func (c *LatencyCounter) Percentile(p float64) float64 {
	if len(c.stats.values) == 0 {
		return 0
	}
	sorted := make([]float64, len(c.stats.values))
	copy(sorted, c.stats.values)
	sort.Float64s(sorted)
	idx := int(math.Ceil(p/100.0*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	return sorted[idx]
}

// Stats returns a copy of the current statistics.
func (c *LatencyCounter) Stats() LatencyStats {
	s := c.stats
	if s.Count == 0 {
		s.Min = 0
	}
	return s
}

// WriteLatencyReport writes a latency summary table to w.
func WriteLatencyReport(w io.Writer, c *LatencyCounter) {
	s := c.Stats()
	if s.Count == 0 {
		fmt.Fprintln(w, "no latency data found")
		return
	}
	avg := s.Sum / float64(s.Count)
	fmt.Fprintf(w, "%-10s %10s\n", "metric", "value (ms)")
	fmt.Fprintf(w, "%-10s %10s\n", strings.Repeat("-", 10), strings.Repeat("-", 10))
	fmt.Fprintf(w, "%-10s %10d\n", "count", s.Count)
	fmt.Fprintf(w, "%-10s %10.3f\n", "min", s.Min)
	fmt.Fprintf(w, "%-10s %10.3f\n", "max", s.Max)
	fmt.Fprintf(w, "%-10s %10.3f\n", "avg", avg)
	fmt.Fprintf(w, "%-10s %10.3f\n", "p50", c.Percentile(50))
	fmt.Fprintf(w, "%-10s %10.3f\n", "p95", c.Percentile(95))
	fmt.Fprintf(w, "%-10s %10.3f\n", "p99", c.Percentile(99))
}
