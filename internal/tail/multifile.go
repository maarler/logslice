package tail

import (
	"context"
	"sync"
)

// Entry pairs a line with the source file it came from.
type Entry struct {
	Path string
	Line string
}

// MultiFile fans-in lines from multiple files into a single channel.
// Each file is tailed with the same Options. Lines arrive in the order
// they are produced; no ordering across files is guaranteed.
type MultiFile struct {
	paths []string
	opts  Options
}

// NewMultiFile creates a MultiFile tailer for the given paths.
func NewMultiFile(paths []string, opts Options) *MultiFile {
	return &MultiFile{paths: paths, opts: opts}
}

// Lines starts a goroutine per file and merges their output.
// The returned channel is closed once all files are exhausted or ctx is done.
func (m *MultiFile) Lines(ctx context.Context) (<-chan Entry, <-chan error) {
	out := make(chan Entry, 128)
	errCh := make(chan error, len(m.paths))

	var wg sync.WaitGroup
	for _, p := range m.paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			t := New(path, m.opts)
			lines, errs := t.Lines(ctx)
			for {
				select {
				case l, ok := <-lines:
					if !ok {
						if err := <-errs; err != nil {
							errCh <- err
						}
						return
					}
					select {
					case out <- Entry{Path: path, Line: l}:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(out)
		close(errCh)
	}()

	return out, errCh
}
