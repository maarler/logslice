// Package highlight provides ANSI colour highlighting for matched patterns in log lines.
package highlight

import (
	"regexp"
	"strings"
)

// ANSI colour codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Rule pairs a compiled regexp with the ANSI colour to apply on a match.
type Rule struct {
	re    *regexp.Regexp
	color string
}

// NewRule compiles pattern and returns a Rule that highlights matches with color.
// Returns an error if pattern is not a valid regular expression.
func NewRule(pattern, color string) (Rule, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return Rule{}, err
	}
	return Rule{re: re, color: color}, nil
}

// Highlighter applies a set of Rules to log lines.
type Highlighter struct {
	rules []Rule
}

// New returns a Highlighter that applies the given rules in order.
func New(rules []Rule) *Highlighter {
	return &Highlighter{rules: rules}
}

// Line applies all rules to line and returns the highlighted result.
// When no rules are defined the original line is returned unchanged.
func (h *Highlighter) Line(line string) string {
	if len(h.rules) == 0 {
		return line
	}
	for _, r := range h.rules {
		line = r.re.ReplaceAllStringFunc(line, func(match string) string {
			return r.color + match + Reset
		})
	}
	return line
}

// Lines applies highlighting to every element of lines and returns a new slice.
func (h *Highlighter) Lines(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = h.Line(l)
	}
	return out
}

// StripANSI removes all ANSI escape sequences from s.
func StripANSI(s string) string {
	var b strings.Builder
	inEsc := false
	for _, ch := range s {
		switch {
		case ch == '\033':
			inEsc = true
		case inEsc && ch == 'm':
			inEsc = false
		case !inEsc:
			b.WriteRune(ch)
		}
	}
	return b.String()
}
