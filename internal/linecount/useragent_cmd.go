package linecount

import (
	"bufio"
	"io"
	"os"
)

// UserAgentCmdOptions controls RunUserAgentCount behaviour.
type UserAgentCmdOptions struct {
	// Pattern is the regex used to extract the user-agent string.
	// Leave empty to use the default quoted-string heuristic.
	Pattern string
	// TopN limits output to the N most frequent agents. 0 means unlimited.
	TopN int
	// Output is where the report is written. Defaults to os.Stdout.
	Output io.Writer
}

// DefaultUserAgentCmdOptions returns sensible defaults.
func DefaultUserAgentCmdOptions() UserAgentCmdOptions {
	return UserAgentCmdOptions{
		TopN:   20,
		Output: os.Stdout,
	}
}

// RunUserAgentCount reads lines from r, counts user-agents and writes a report.
func RunUserAgentCount(r io.Reader, opts UserAgentCmdOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	counter, err := NewUserAgentCounter(opts.Pattern)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		counter.Add(sc.Text())
	}
	if err := sc.Err(); err != nil {
		return err
	}

	entries := SortedUserAgentEntries(counter.Counts())
	if opts.TopN > 0 && len(entries) > opts.TopN {
		entries = entries[:opts.TopN]
	}

	WriteUserAgentReport(opts.Output, entries, counter.Total())
	return nil
}
