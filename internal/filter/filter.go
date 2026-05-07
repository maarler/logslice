package filter

import (
	"bufio"
	"io"
)

// Result holds the lines that passed the filter and a count of those dropped.
type Result struct {
	Lines   []string
	Dropped int
}

// Apply reads all lines from r and returns those that satisfy the chain.
func Apply(r io.Reader, chain *Chain) (*Result, error) {
	result := &Result{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if chain.Match(line) {
			result.Lines = append(result.Lines, line)
		} else {
			result.Dropped++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// ApplyLines filters an already-split slice of lines.
func ApplyLines(lines []string, chain *Chain) *Result {
	result := &Result{}
	for _, line := range lines {
		if chain.Match(line) {
			result.Lines = append(result.Lines, line)
		} else {
			result.Dropped++
		}
	}
	return result
}
