package linecount

import (
	"bufio"
	"io"

	"github.com/user/logslice/internal/parser"
)

// HourlyCmdOptions configures RunHourlyCount.
type HourlyCmdOptions struct {
	// TimestampFormat overrides auto-detection when non-empty.
	TimestampFormat string
	// Out is where the report is written.
	Out io.Writer
}

// DefaultHourlyCmdOptions returns sensible defaults.
func DefaultHourlyCmdOptions() HourlyCmdOptions {
	return HourlyCmdOptions{}
}

// RunHourlyCount reads log lines from r, parses timestamps, and writes an
// hourly breakdown report to opts.Out.
func RunHourlyCount(r io.Reader, opts HourlyCmdOptions) error {
	counter := NewHourlyCounter()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		var (
			t   interface{ IsZero() bool }
			err error
		)
		if opts.TimestampFormat != "" {
			parsed, e := parser.ParseTimestampWithFormat(line, opts.TimestampFormat)
			t, err = parsed, e
		} else {
			parsed, e := parser.ParseTimestamp(line)
			t, err = parsed, e
		}
		if err != nil || t.IsZero() {
			continue
		}
		switch v := t.(type) {
		case interface{ UnixNano() int64 }:
			_ = v // handled below via reflection-free approach
		}
		// Re-parse cleanly to obtain time.Time.
		if opts.TimestampFormat != "" {
			if ts, e := parser.ParseTimestampWithFormat(line, opts.TimestampFormat); e == nil {
				counter.Add(ts)
			}
		} else {
			if ts, e := parser.ParseTimestamp(line); e == nil {
				counter.Add(ts)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return WriteHourlyReport(opts.Out, counter.Entries())
}
