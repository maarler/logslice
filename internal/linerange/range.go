// Package linerange provides line-number based range filtering for log lines.
package linerange

import (
	"fmt"
	"io"
	"strings"
)

// Range represents an inclusive line range [Start, End].
// A zero End means open-ended (read until EOF).
type Range struct {
	Start int64
	End   int64 // 0 = unbounded
}

// Validate returns an error if the range is not logically valid.
func (r Range) Validate() error {
	if r.Start < 1 {
		return fmt.Errorf("linerange: start must be >= 1, got %d", r.Start)
	}
	if r.End != 0 && r.End < r.Start {
		return fmt.Errorf("linerange: end (%d) must be >= start (%d)", r.End, r.Start)
	}
	return nil
}

// Contains reports whether line number n (1-based) falls within the range.
func (r Range) Contains(n int64) bool {
	if n < r.Start {
		return false
	}
	if r.End == 0 {
		return true
	}
	return n <= r.End
}

// String returns a human-readable representation of the range.
func (r Range) String() string {
	if r.End == 0 {
		return fmt.Sprintf("%d-", r.Start)
	}
	return fmt.Sprintf("%d-%d", r.Start, r.End)
}

// Parse parses a range string in the form "N", "N-", or "N-M".
func Parse(s string) (Range, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Range{}, fmt.Errorf("linerange: empty range string")
	}
	var r Range
	if !strings.Contains(s, "-") {
		var n int64
		if _, err := fmt.Sscanf(s, "%d", &n); err != nil {
			return Range{}, fmt.Errorf("linerange: invalid range %q", s)
		}
		r = Range{Start: n, End: n}
		return r, r.Validate()
	}
	parts := strings.SplitN(s, "-", 2)
	if _, err := fmt.Sscanf(parts[0], "%d", &r.Start); err != nil {
		return Range{}, fmt.Errorf("linerange: invalid start in %q", s)
	}
	if parts[1] != "" {
		if _, err := fmt.Sscanf(parts[1], "%d", &r.End); err != nil {
			return Range{}, fmt.Errorf("linerange: invalid end in %q", s)
		}
	}
	return r, r.Validate()
}

// Apply reads lines from r and writes only those within the range to w.
// Lines are 1-indexed. It stops reading early when End is exceeded.
func Apply(r Range, src io.Reader, dst io.Writer) error {
	if err := r.Validate(); err != nil {
		return err
	}
	buf := make([]byte, 0, 4096)
	var lineNum int64
	tmp := make([]byte, 1)
	for {
		n, err := src.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[0])
			if tmp[0] == '\n' {
				lineNum++
				if r.Contains(lineNum) {
					if _, werr := dst.Write(buf); werr != nil {
						return werr
					}
				}
				if r.End != 0 && lineNum >= r.End {
					return nil
				}
				buf = buf[:0]
			}
		}
		if err == io.EOF {
			if len(buf) > 0 {
				lineNum++
				if r.Contains(lineNum) {
					_, werr := dst.Write(buf)
					return werr
				}
			}
			return nil
		}
		if err != nil {
			return err
		}
	}
}
