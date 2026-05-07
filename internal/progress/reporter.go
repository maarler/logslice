package progress

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// Reporter tracks and displays progress of log processing.
type Reporter struct {
	out        io.Writer
	linesRead  atomic.Int64
	bytesRead  atomic.Int64
	start      time.Time
	ticker     *time.Ticker
	done       chan struct{}
	quiet      bool
}

// New creates a new Reporter that writes progress to out.
// If quiet is true, no output is produced.
func New(out io.Writer, quiet bool) *Reporter {
	return &Reporter{
		out:   out,
		start: time.Now(),
		done:  make(chan struct{}),
		quiet: quiet,
	}
}

// Start begins periodic progress reporting at the given interval.
func (r *Reporter) Start(interval time.Duration) {
	if r.quiet {
		return
	}
	r.ticker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-r.ticker.C:
				r.print()
			case <-r.done:
				return
			}
		}
	}()
}

// Add records that n lines and b bytes have been processed.
func (r *Reporter) Add(lines int64, bytes int64) {
	r.linesRead.Add(lines)
	r.bytesRead.Add(bytes)
}

// Stop halts periodic reporting and prints a final summary.
func (r *Reporter) Stop() {
	if r.ticker != nil {
		r.ticker.Stop()
		close(r.done)
	}
	if !r.quiet {
		r.print()
		fmt.Fprintln(r.out)
	}
}

// Lines returns the total number of lines processed so far.
func (r *Reporter) Lines() int64 {
	return r.linesRead.Load()
}

// Bytes returns the total number of bytes processed so far.
func (r *Reporter) Bytes() int64 {
	return r.bytesRead.Load()
}

func (r *Reporter) print() {
	elapsed := time.Since(r.start).Truncate(time.Millisecond)
	lines := r.linesRead.Load()
	bytes := r.bytesRead.Load()
	fmt.Fprintf(r.out, "\r[%s] %d lines, %s processed",
		elapsed, lines, formatBytes(bytes))
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
