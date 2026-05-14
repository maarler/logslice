// Package pipeline chains processing stages for log lines.
package pipeline

import (
	"context"
	"fmt"
)

// Pipeline connects a series of Stage functions, piping the output of each
// into the input of the next.
type Pipeline struct {
	stages []Stage
	bufSize int
}

// New creates a Pipeline with the given stages.
// bufSize controls the channel buffer between stages (default 64).
func New(bufSize int, stages ...Stage) *Pipeline {
	if bufSize <= 0 {
		bufSize = 64
	}
	return &Pipeline{stages: stages, bufSize: bufSize}
}

// Run executes the pipeline, reading from src and sending final output to dst.
// It blocks until all stages complete or ctx is cancelled.
func (p *Pipeline) Run(ctx context.Context, src <-chan string, dst chan<- string) error {
	if len(p.stages) == 0 {
		return drainTo(ctx, src, dst)
	}

	channels := make([]chan string, len(p.stages)+1)
	channels[0] = make(chan string, p.bufSize)
	for i := 1; i <= len(p.stages); i++ {
		channels[i] = make(chan string, p.bufSize)
	}

	// Forward src into the first internal channel.
	go func() {
		defer close(channels[0])
		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-src:
				if !ok {
					return
				}
				channels[0] <- line
			}
		}
	}()

	errCh := make(chan error, len(p.stages))
	for i, stage := range p.stages {
		i, stage := i, stage
		go func() {
			defer close(channels[i+1])
			errCh <- stage(ctx, channels[i], channels[i+1])
		}()
	}

	// Drain the final stage output into dst.
	go func() {
		for line := range channels[len(p.stages)] {
			dst <- line
		}
	}()

	for range p.stages {
		if err := <-errCh; err != nil {
			return fmt.Errorf("pipeline stage error: %w", err)
		}
	}
	return nil
}

func drainTo(ctx context.Context, src <-chan string, dst chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-src:
			if !ok {
				return nil
			}
			dst <- line
		}
	}
}
