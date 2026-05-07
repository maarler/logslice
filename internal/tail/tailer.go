package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Options controls tail behaviour.
type Options struct {
	// PollInterval is how often to check for new data when the end of file is reached.
	PollInterval time.Duration
	// Follow keeps the file open and waits for new lines (like tail -f).
	Follow bool
}

// DefaultOptions returns sensible tail defaults.
func DefaultOptions() Options {
	return Options{
		PollInterval: 250 * time.Millisecond,
		Follow:       false,
	}
}

// Tailer reads lines from a file, optionally following new writes.
type Tailer struct {
	path string
	opts Options
}

// New creates a new Tailer for the given file path.
func New(path string, opts Options) *Tailer {
	return &Tailer{path: path, opts: opts}
}

// Lines streams lines from the file into the returned channel.
// The channel is closed when the context is cancelled or (when not following)
// EOF is reached. Any read error is sent on errCh.
func (t *Tailer) Lines(ctx context.Context) (<-chan string, <-chan error) {
	lines := make(chan string, 64)
	errCh := make(chan error, 1)

	go func() {
		defer close(lines)
		defer close(errCh)

		f, err := os.Open(t.path)
		if err != nil {
			errCh <- err
			return
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line, err := reader.ReadString('\n')
			if len(line) > 0 {
				// Strip trailing newline.
				if line[len(line)-1] == '\n' {
					line = line[:len(line)-1]
				}
				select {
				case lines <- line:
				case <-ctx.Done():
					return
				}
			}

			if err == io.EOF {
				if !t.opts.Follow {
					return
				}
				select {
				case <-time.After(t.opts.PollInterval):
				case <-ctx.Done():
					return
				}
				continue
			}
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	return lines, errCh
}
