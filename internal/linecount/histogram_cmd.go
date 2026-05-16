package linecount

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/user/logslice/internal/parser"
)

// HistogramOptions controls RunHistogram behaviour.
type HistogramOptions struct {
	// Window is the bucket width (e.g. time.Minute, time.Hour).
	Window time.Duration
	// BarWidth is the maximum ASCII bar width in characters.
	BarWidth int
	// TimestampFormat is an optional strftime-style format hint.
	TimestampFormat string
}

// DefaultHistogramOptions returns sensible defaults.
func DefaultHistogramOptions() HistogramOptions {
	return HistogramOptions{
		Window:   time.Minute,
		BarWidth: 40,
	}
}

// RunHistogram reads lines from r, parses timestamps, buckets them by the
// configured window, and writes an ASCII histogram to out.
func RunHistogram(r io.Reader, out io.Writer, opts HistogramOptions) error {
	if opts.Window <= 0 {
		opts.Window = time.Minute
	}
	if opts.BarWidth <= 0 {
		opts.BarWidth = 40
	}

	h := NewTimeHistogram(opts.Window)
	scanner := bufio.NewScanner(r)
	var skipped int64

	for scanner.Scan() {
		line := scanner.Text()
		var ts time.Time
		var err error
		if opts.TimestampFormat != "" {
			ts, err = parser.ParseTimestampWithFormat(line, opts.TimestampFormat)
		} else {
			ts, err = parser.ParseTimestamp(line)
		}
		if err != nil {
			skipped++
			continue
		}
		h.Add(ts, int64(len(line)))
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("histogram: scan error: %w", err)
	}

	WriteHistogramReport(out, h, opts.BarWidth)
	if skipped > 0 {
		fmt.Fprintf(out, "\n(skipped %d lines with no parseable timestamp)\n", skipped)
	}
	return nil
}
