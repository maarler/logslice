// Package mask provides log line redaction for sensitive data patterns.
package mask

import (
	"regexp"
	"strings"
)

// Rule describes a single masking rule: a compiled pattern and its replacement.
type Rule struct {
	pattern     *regexp.Regexp
	replacement string
	label       string
}

// NewRule compiles a masking rule from a raw regex pattern.
// replacement is the string substituted for each match (e.g. "[REDACTED]").
func NewRule(label, pattern, replacement string) (*Rule, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Rule{pattern: re, replacement: replacement, label: label}, nil
}

// String returns a human-readable description of the rule.
func (r *Rule) String() string {
	return r.label + ":" + r.pattern.String()
}

// Masker applies a set of Rules to log lines.
type Masker struct {
	rules []*Rule
}

// New creates a Masker with the given rules.
func New(rules ...*Rule) *Masker {
	return &Masker{rules: rules}
}

// Line applies all masking rules to a single log line and returns the result.
func (m *Masker) Line(line string) string {
	for _, r := range m.rules {
		line = r.pattern.ReplaceAllString(line, r.replacement)
	}
	return line
}

// Lines applies masking to each element of lines in-place and returns the slice.
func (m *Masker) Lines(lines []string) []string {
	for i, l := range lines {
		lines[i] = m.Line(l)
	}
	return lines
}

// Apply reads all content from src, masks every line, and returns the result.
func (m *Masker) Apply(src string) string {
	if len(m.rules) == 0 {
		return src
	}
	lines := strings.Split(src, "\n")
	for i, l := range lines {
		lines[i] = m.Line(l)
	}
	return strings.Join(lines, "\n")
}

// Preset returns a Masker pre-loaded with common sensitive-data rules.
func Preset() (*Masker, error) {
	defs := []struct{ label, pattern, repl string }{
		{"ipv4", `\b(?:\d{1,3}\.){3}\d{1,3}\b`, "[IP]"},
		{"email", `[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`, "[EMAIL]"},
		{"token", `(?i)(token|key|secret|password)=[^\s&]+`, "$1=[REDACTED]"},
	}
	var rules []*Rule
	for _, d := range defs {
		r, err := NewRule(d.label, d.pattern, d.repl)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return New(rules...), nil
}
