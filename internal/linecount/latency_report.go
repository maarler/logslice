package linecount

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteLatencyReport writes a formatted latency statistics table to w.
func WriteLatencyReport(w io.Writer, c *LatencyCounter) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "Latency Report")
	fmt.Fprintln(tw, repeatChar('-', 50))

	if c.Count() == 0 {
		fmt.Fprintln(tw, "No latency samples recorded.")
		return tw.Flush()
	}

	fmt.Fprintf(tw, "Samples:\t%d\n", c.Count())
	fmt.Fprintf(tw, "min:\t%s\n", formatLatencyDuration(c.Min()))
	fmt.Fprintf(tw, "max:\t%s\n", formatLatencyDuration(c.Max()))
	fmt.Fprintf(tw, "mean:\t%s\n", formatLatencyDuration(c.Mean()))

	fmt.Fprintln(tw, "")
	fmt.Fprintln(tw, "Percentiles")
	fmt.Fprintln(tw, repeatChar('-', 30))

	percentiles := []float64{50, 75, 90, 95, 99, 99.9}
	for _, p := range percentiles {
		val := c.Percentile(p)
		if p == 99.9 {
			fmt.Fprintf(tw, "p99.9:\t%s\n", formatLatencyDuration(val))
		} else {
			fmt.Fprintf(tw, "p%.0f:\t%s\n", p, formatLatencyDuration(val))
		}
	}

	return tw.Flush()
}

// formatLatencyDuration formats a float64 millisecond value into a human-readable string.
func formatLatencyDuration(ms float64) string {
	switch {
	case ms >= 60_000:
		return fmt.Sprintf("%.2fm", ms/60_000)
	case ms >= 1_000:
		return fmt.Sprintf("%.2fs", ms/1_000)
	case ms >= 1:
		return fmt.Sprintf("%.2fms", ms)
	default:
		return fmt.Sprintf("%.2fµs", ms*1_000)
	}
}
