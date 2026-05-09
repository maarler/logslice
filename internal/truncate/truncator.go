// Package truncate provides utilities for truncating log lines that exceed
// a maximum byte length, optionally appending a suffix to indicate truncation.
package truncate

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

const (
	DefaultMaxBytes = 8192
	DefaultSuffix   = "...[truncated]"
)

// Options configures the Truncator.
type Options struct {
	// MaxBytes is the maximum number of bytes allowed per line (excluding newline).
	// Lines longer than this will be truncated. Zero means use DefaultMaxBytes.
	MaxBytes int
	// Suffix is appended to truncated lines. Defaults to DefaultSuffix.
	Suffix string
}

// Truncator truncates log lines that exceed a maximum byte length.
type Truncator struct {
	maxBytes int
	suffix   []byte
}

// New creates a new Truncator with the given options.
func New(opts Options) *Truncator {
	if opts.MaxBytes <= 0 {
		opts.MaxBytes = DefaultMaxBytes
	}
	if opts.Suffix == "" {
		opts.Suffix = DefaultSuffix
	}
	return &Truncator{
		maxBytes: opts.MaxBytes,
		suffix:   []byte(opts.Suffix),
	}
}

// Line truncates a single line if it exceeds MaxBytes.
// The returned slice shares no memory with the input.
func (t *Truncator) Line(line []byte) []byte {
	if len(line) <= t.maxBytes {
		out := make([]byte, len(line))
		copy(out, line)
		return out
	}
	cutAt := t.maxBytes - len(t.suffix)
	if cutAt < 0 {
		cutAt = 0
	}
	out := make([]byte, 0, t.maxBytes)
	out = append(out, line[:cutAt]...)
	out = append(out, t.suffix...)
	return out
}

// Apply reads lines from r, truncates each one, and writes the result to w.
func (t *Truncator) Apply(r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 64*1024)
	for sc.Scan() {
		truncated := t.Line(sc.Bytes())
		if _, err := w.Write(append(truncated, '\n')); err != nil {
			return err
		}
	}
	return sc.Err()
}

// ApplyLines truncates each string in lines and returns a new slice.
func (t *Truncator) ApplyLines(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = string(t.Line([]byte(l)))
	}
	return out
}

// CountTruncated returns how many lines in the input would be truncated.
func (t *Truncator) CountTruncated(lines []string) int {
	n := 0
	for _, l := range lines {
		if len(l) > t.maxBytes {
			n++
		}
	}
	return n
}

// lineFromReader is a helper used in tests.
func lineFromReader(r io.Reader) string {
	var sb strings.Builder
	io.Copy(&sb, r) //nolint:errcheck
	return strings.TrimRight(sb.String(), "\n")
}

// readerFromLines builds an io.Reader from a slice of lines.
func readerFromLines(lines []string) io.Reader {
	return bytes.NewBufferString(strings.Join(lines, "\n") + "\n")
}
