// Package linerange provides line-number based slicing of log input.
// It allows callers to extract a contiguous range [First, Last] of lines
// from a stream, discarding everything outside that window.
package linerange

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// Range defines an inclusive line-number window (1-based).
type Range struct {
	First int64
	Last  int64 // 0 means "until EOF"
}

// Validate returns an error if the Range is logically invalid.
func (r Range) Validate() error {
	if r.First < 1 {
		return errors.New("linerange: First must be >= 1")
	}
	if r.Last != 0 && r.Last < r.First {
		return fmt.Errorf("linerange: Last (%d) must be >= First (%d)", r.Last, r.First)
	}
	return nil
}

// Contains reports whether the given 1-based line number falls inside the range.
func (r Range) Contains(n int64) bool {
	if n < r.First {
		return false
	}
	if r.Last != 0 && n > r.Last {
		return false
	}
	return true
}

// Past reports whether n is beyond the end of the range.
func (r Range) Past(n int64) bool {
	return r.Last != 0 && n > r.Last
}

// Apply reads lines from src, writing only those within the range to dst.
// It stops reading early once the range is exhausted.
func Apply(src io.Reader, dst io.Writer, r Range) (kept int64, err error) {
	if err = r.Validate(); err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(src)
	var lineNum int64

	for scanner.Scan() {
		lineNum++

		if r.Past(lineNum) {
			break
		}

		if !r.Contains(lineNum) {
			continue
		}

		if _, werr := fmt.Fprintf(dst, "%s\n", scanner.Bytes()); werr != nil {
			return kept, werr
		}
		kept++
	}

	return kept, scanner.Err()
}
