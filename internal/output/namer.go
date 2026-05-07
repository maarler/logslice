package output

import (
	"fmt"
	"time"
)

// Granularity controls how timestamps are truncated for segment naming.
type Granularity int

const (
	ByMinute Granularity = iota
	ByHour
	ByDay
)

// Namer generates segment keys from timestamps or patterns.
type Namer struct {
	Granularity Granularity
}

// NewNamer returns a Namer with the given granularity.
func NewNamer(g Granularity) *Namer {
	return &Namer{Granularity: g}
}

// KeyFromTime returns a string key for the given time based on granularity.
func (n *Namer) KeyFromTime(t time.Time) string {
	switch n.Granularity {
	case ByMinute:
		return t.UTC().Format("2006-01-02T15-04")
	case ByHour:
		return t.UTC().Format("2006-01-02T15")
	case ByDay:
		return t.UTC().Format("2006-01-02")
	default:
		return t.UTC().Format("2006-01-02T15-04")
	}
}

// KeyFromPattern returns a sanitized key derived from a matched pattern label.
func (n *Namer) KeyFromPattern(label string) string {
	return sanitize(label)
}

// sanitize replaces characters unsafe for filenames with underscores.
func sanitize(s string) string {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '.' {
			out[i] = c
		} else {
			out[i] = '_'
		}
	}
	return fmt.Sprintf("%s", out)
}
