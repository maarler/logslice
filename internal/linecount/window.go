package linecount

import (
	"io"
	"strings"
	"sync"
)

// WindowEntry records line and byte counts for a single named window.
type WindowEntry struct {
	Key   string
	Lines int64
	Bytes int64
}

// WindowCounter accumulates per-window statistics in insertion order.
type WindowCounter struct {
	mu      sync.Mutex
	index   map[string]int // key -> position in windows slice
	windows []*WindowEntry
}

// NewWindowCounter returns an initialised WindowCounter.
func NewWindowCounter() *WindowCounter {
	return &WindowCounter{
		index: make(map[string]int),
	}
}

// Add records a single line of byteLen bytes under key.
func (wc *WindowCounter) Add(key string, byteLen int64) {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	if idx, ok := wc.index[key]; ok {
		wc.windows[idx].Lines++
		wc.windows[idx].Bytes += byteLen
		return
	}
	wc.index[key] = len(wc.windows)
	wc.windows = append(wc.windows, &WindowEntry{Key: key, Lines: 1, Bytes: byteLen})
}

// Entries returns a snapshot of all window entries in insertion order.
func (wc *WindowCounter) Entries() []WindowEntry {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	out := make([]WindowEntry, len(wc.windows))
	for i, w := range wc.windows {
		out[i] = *w
	}
	return out
}

// CountWindowReader reads all lines from r, groups them by the key returned
// by keyFn, and returns a populated WindowCounter.
func CountWindowReader(r io.Reader, keyFn func(line string) string) (*WindowCounter, error) {
	wc := NewWindowCounter()
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(buf), "\n") {
		if line == "" {
			continue
		}
		key := keyFn(line)
		if key == "" {
			continue
		}
		wc.Add(key, int64(len(line)+1))
	}
	return wc, nil
}
