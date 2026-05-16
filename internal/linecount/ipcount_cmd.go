package linecount

import (
	"io"
)

// IPCmdOptions configures RunIPCount.
type IPCmdOptions struct {
	// TopN limits output to the N most frequent IPs. 0 means no limit.
	TopN int
}

// DefaultIPCmdOptions returns sensible defaults.
func DefaultIPCmdOptions() IPCmdOptions {
	return IPCmdOptions{TopN: 20}
}

// RunIPCount reads log lines from r, counts IP addresses, and writes a
// report to w. Options control how many top entries to display.
func RunIPCount(r io.Reader, w io.Writer, opts IPCmdOptions) error {
	counter := CountIPReader(r)
	entries := SortedIPEntries(counter.Counts())

	if opts.TopN > 0 && len(entries) > opts.TopN {
		entries = entries[:opts.TopN]
	}

	return WriteIPReport(w, entries, counter.Total())
}
