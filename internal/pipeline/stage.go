package pipeline

import "context"

// Stage represents a single processing step in a log pipeline.
// It consumes lines from in and emits transformed lines to out.
type Stage func(ctx context.Context, in <-chan string, out chan<- string) error

// Processor is a function that transforms a single log line.
// Returning an empty string drops the line from the stream.
type Processor func(line string) string

// FromProcessor wraps a Processor into a Stage.
func FromProcessor(p Processor) Stage {
	return func(ctx context.Context, in <-chan string, out chan<- string) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case line, ok := <-in:
				if !ok {
					return nil
				}
				if result := p(line); result != "" {
					select {
					case out <- result:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			}
		}
	}
}
