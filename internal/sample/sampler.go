// Package sample provides line sampling strategies for large log files.
package sample

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
)

// Strategy defines how lines are selected during sampling.
type Strategy int

const (
	// StrategyRandom selects lines with a fixed probability.
	StrategyRandom Strategy = iota
	// StrategyNth selects every Nth line.
	StrategyNth
)

// Options controls sampler behaviour.
type Options struct {
	Strategy Strategy
	// Rate is the fraction of lines to keep (0.0–1.0) for StrategyRandom.
	Rate float64
	// N is the step size for StrategyNth.
	N int
	// Seed is used to initialise the random source; 0 means use default.
	Seed int64
}

// Sampler filters lines from a reader according to the chosen strategy.
type Sampler struct {
	opts Options
	rng  *rand.Rand
}

// New returns a Sampler configured with opts.
// It returns an error if the options are invalid.
func New(opts Options) (*Sampler, error) {
	if opts.Strategy == StrategyRandom {
		if opts.Rate <= 0 || opts.Rate > 1 {
			return nil, fmt.Errorf("sample: rate must be in (0, 1], got %f", opts.Rate)
		}
	}
	if opts.Strategy == StrategyNth {
		if opts.N < 1 {
			return nil, fmt.Errorf("sample: N must be >= 1, got %d", opts.N)
		}
	}
	src := rand.NewSource(opts.Seed)
	return &Sampler{opts: opts, rng: rand.New(src)}, nil
}

// Sample reads lines from r and writes selected lines to w.
// It returns the number of lines read and the number written.
func (s *Sampler) Sample(r io.Reader, w io.Writer) (read, written int, err error) {
	scanner := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)
	for scanner.Scan() {
		read++
		if s.keep(read) {
			if _, werr := bw.WriteString(scanner.Text() + "\n"); werr != nil {
				return read, written, werr
			}
			written++
		}
	}
	if serr := scanner.Err(); serr != nil {
		return read, written, serr
	}
	return read, written, bw.Flush()
}

func (s *Sampler) keep(lineNum int) bool {
	switch s.opts.Strategy {
	case StrategyNth:
		return lineNum%s.opts.N == 0
	default: // StrategyRandom
		return s.rng.Float64() < s.opts.Rate
	}
}
