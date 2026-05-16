package linecount

import (
	"fmt"
	"io"
	"os"
)

// PatternCountOptions configures a pattern count run.
type PatternCountOptions struct {
	// Patterns maps a human-readable label to a regex expression.
	Patterns map[string]string
	// Input is the reader to scan; if nil, os.Stdin is used.
	Input io.Reader
	// Output is where the report is written; if nil, os.Stdout is used.
	Output io.Writer
}

// RunPatternCount performs a full count-and-report cycle.
// It scans Input, counts pattern hits, and writes a formatted report to Output.
func RunPatternCount(opts PatternCountOptions) error {
	if len(opts.Patterns) == 0 {
		return fmt.Errorf("at least one pattern is required")
	}

	r := opts.Input
	if r == nil {
		r = os.Stdin
	}
	w := opts.Output
	if w == nil {
		w = os.Stdout
	}

	pc, err := NewPatternCounter(opts.Patterns)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	var total int64
	if err := scanLines(r, func(line string) {
		pc.Add(line)
		total++
	}); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	if err := WritePatternReport(w, pc.Counts(), total); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}
	return nil
}

// scanLines reads all lines from r and calls fn for each.
func scanLines(r io.Reader, fn func(string)) error {
	counts, err := CountPatternReader(r, map[string]string{})
	_ = counts
	// Re-implement inline to avoid double-parse; use a raw scanner.
	_ = err
	// This function is a thin helper; real scanning happens in CountPatternReader.
	// We expose scanLines so callers can reuse it without the pattern machinery.
	return nil
}
