package linecount

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// LatencyCmdOptions configures RunLatencyCount.
type LatencyCmdOptions struct {
	// Field is the key to extract latency values from (e.g. "duration_ms").
	Field string
	// Sep is the token separator within a log line (default: space).
	Sep string
	// KVSep is the key-value separator (default: "=").
	KVSep string
	// Output is where the report is written; defaults to os.Stdout.
	Output io.Writer
}

// DefaultLatencyCmdOptions returns sensible defaults.
func DefaultLatencyCmdOptions() LatencyCmdOptions {
	return LatencyCmdOptions{
		Field: "duration_ms",
		Sep:   " ",
		KVSep: "=",
		Output: os.Stdout,
	}
}

// RunLatencyCount reads lines from r, extracts latency values, and writes a
// percentile report to opts.Output.
func RunLatencyCount(r io.Reader, opts LatencyCmdOptions) error {
	if opts.Field == "" {
		return fmt.Errorf("latency: field name must not be empty")
	}
	if opts.Output == nil {
		opts.Output = os.Stdout
	}
	if opts.Sep == "" {
		opts.Sep = " "
	}
	if opts.KVSep == "" {
		opts.KVSep = "="
	}

	counter := NewLatencyCounter(opts.Field, opts.Sep, opts.KVSep)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		counter.Add(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("latency: reading input: %w", err)
	}

	WriteLatencyReport(opts.Output, counter)
	return nil
}
