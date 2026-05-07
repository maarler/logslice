package filter

// Chain holds an ordered list of patterns. A line passes the chain only
// when it satisfies every pattern in the list.
type Chain struct {
	patterns []*Pattern
}

// NewChain creates an empty filter chain.
func NewChain() *Chain {
	return &Chain{}
}

// Add appends a pattern to the chain.
func (c *Chain) Add(p *Pattern) {
	c.patterns = append(c.patterns, p)
}

// AddExpr is a convenience wrapper that compiles expr and adds it.
func (c *Chain) AddExpr(expr string, negate bool) error {
	p, err := NewPattern(expr, negate)
	if err != nil {
		return err
	}
	c.Add(p)
	return nil
}

// Match returns true when the line satisfies all patterns in the chain.
// An empty chain matches every line.
func (c *Chain) Match(line string) bool {
	for _, p := range c.patterns {
		if !p.Match(line) {
			return false
		}
	}
	return true
}

// Len returns the number of patterns in the chain.
func (c *Chain) Len() int {
	return len(c.patterns)
}
