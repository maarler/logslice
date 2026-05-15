package linerange

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Range represents an inclusive line range [Start, End].
// A zero End means open-ended (no upper bound).
type Range struct {
	Start int
	End   int // 0 = unbounded
}

// Validate checks that the range fields are self-consistent.
func (r Range) Validate() error {
	if r.Start < 1 {
		return fmt.Errorf("start must be >= 1, got %d", r.Start)
	}
	if r.End != 0 && r.End < r.Start {
		return fmt.Errorf("end (%d) must be >= start (%d)", r.End, r.Start)
	}
	return nil
}

// Contains reports whether lineNo (1-based) falls within the range.
func (r Range) Contains(lineNo int) bool {
	if lineNo < r.Start {
		return false
	}
	if r.End != 0 && lineNo > r.End {
		return false
	}
	return true
}

// Parse parses a range spec of the form "start:end" or "start:".
// Both start and end are 1-based line numbers.
func Parse(spec string) (Range, error) {
	parts := strings.SplitN(spec, ":", 2)
	if len(parts) != 2 {
		return Range{}, fmt.Errorf("invalid range spec %q: expected 'start:end' or 'start:'")
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return Range{}, fmt.Errorf("invalid start in range spec %q: %w", spec, err)
	}
	var end int
	if parts[1] != "" {
		end, err = strconv.Atoi(parts[1])
		if err != nil {
			return Range{}, fmt.Errorf("invalid end in range spec %q: %w", spec, err)
		}
	}
	r := Range{Start: start, End: end}
	if err := r.Validate(); err != nil {
		return Range{}, fmt.Errorf("range spec %q: %w", spec, err)
	}
	return r, nil
}

// Apply writes lines from src that fall within r to dst.
// Line numbers are 1-based.
func Apply(r Range, src io.Reader, dst io.Writer) error {
	if err := r.Validate(); err != nil {
		return err
	}
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	for i, line := range lines {
		lineNo := i + 1
		if r.Contains(lineNo) {
			if _, werr := fmt.Fprintln(dst, line); werr != nil {
				return werr
			}
		}
	}
	return nil
}
