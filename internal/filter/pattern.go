package filter

import (
	"fmt"
	"regexp"
)

// Pattern holds a compiled regular expression used to match log lines.
type Pattern struct {
	re      *regexp.Regexp
	raw     string
	negate  bool
}

// NewPattern compiles a regex pattern. If negate is true, the filter
// matches lines that do NOT match the expression.
func NewPattern(expr string, negate bool) (*Pattern, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("filter: invalid pattern %q: %w", expr, err)
	}
	return &Pattern{re: re, raw: expr, negate: negate}, nil
}

// Match reports whether the line satisfies the pattern filter.
func (p *Pattern) Match(line string) bool {
	matched := p.re.MatchString(line)
	if p.negate {
		return !matched
	}
	return matched
}

// String returns a human-readable description of the pattern.
func (p *Pattern) String() string {
	if p.negate {
		return fmt.Sprintf("NOT(%s)", p.raw)
	}
	return p.raw
}
