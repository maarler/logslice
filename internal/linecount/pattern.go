package linecount

import (
	"bufio"
	"io"
	"regexp"
	"sort"
)

// PatternCounter counts line occurrences matching named regex patterns.
type PatternCounter struct {
	rules  []*patternRule
	counts map[string]int64
}

type patternRule struct {
	name string
	re   *regexp.Regexp
}

// NewPatternCounter creates a PatternCounter from a map of name→regex pairs.
func NewPatternCounter(patterns map[string]string) (*PatternCounter, error) {
	pc := &PatternCounter{
		counts: make(map[string]int64, len(patterns)),
	}
	for name, expr := range patterns {
		re, err := regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
		pc.rules = append(pc.rules, &patternRule{name: name, re: re})
		pc.counts[name] = 0
	}
	// stable order for deterministic output
	sort.Slice(pc.rules, func(i, j int) bool {
		return pc.rules[i].name < pc.rules[j].name
	})
	return pc, nil
}

// Add tests line against every pattern and increments matched counters.
func (pc *PatternCounter) Add(line string) {
	for _, r := range pc.rules {
		if r.re.MatchString(line) {
			pc.counts[r.name]++
		}
	}
}

// Counts returns a copy of the current counts map.
func (pc *PatternCounter) Counts() map[string]int64 {
	out := make(map[string]int64, len(pc.counts))
	for k, v := range pc.counts {
		out[k] = v
	}
	return out
}

// CountPatternReader scans r and returns per-pattern hit counts.
func CountPatternReader(r io.Reader, patterns map[string]string) (map[string]int64, error) {
	pc, err := NewPatternCounter(patterns)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		pc.Add(sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return pc.Counts(), nil
}
