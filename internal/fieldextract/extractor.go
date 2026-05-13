// Package fieldextract provides utilities for extracting named fields
// from structured log lines (e.g. key=value or JSON-style pairs).
package fieldextract

import (
	"strings"
)

// Extractor extracts key=value pairs from a log line.
type Extractor struct {
	delimiter string
	pairSep   string
}

// Options configures an Extractor.
type Options struct {
	// Delimiter separates key from value (default "=").
	Delimiter string
	// PairSep separates key-value pairs (default " ").
	PairSep string
}

// New returns an Extractor with the given options.
// Zero-value fields fall back to sensible defaults.
func New(opts Options) *Extractor {
	if opts.Delimiter == "" {
		opts.Delimiter = "="
	}
	if opts.PairSep == "" {
		opts.PairSep = " "
	}
	return &Extractor{
		delimiter: opts.Delimiter,
		pairSep:   opts.PairSep,
	}
}

// Extract parses line and returns a map of all key=value pairs found.
// Values that are quoted with double-quotes have the quotes stripped.
func (e *Extractor) Extract(line string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(line, e.pairSep)
	for _, pair := range pairs {
		idx := strings.Index(pair, e.delimiter)
		if idx < 1 {
			continue
		}
		key := strings.TrimSpace(pair[:idx])
		val := strings.TrimSpace(pair[idx+len(e.delimiter):])
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		if key != "" {
			result[key] = val
		}
	}
	return result
}

// Get extracts a single named field from line.
// Returns the value and true if found, or empty string and false otherwise.
func (e *Extractor) Get(line, key string) (string, bool) {
	fields := e.Extract(line)
	v, ok := fields[key]
	return v, ok
}
