// Package dedupe provides log line deduplication using a sliding hash window.
package dedupe

import (
	"hash/fnv"
	"sync"
)

// Options configures the Deduper.
type Options struct {
	// WindowSize is the number of recent hashes to remember.
	WindowSize int
	// Consecutive, when true, only suppresses back-to-back identical lines.
	Consecutive bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		WindowSize:  512,
		Consecutive: false,
	}
}

// Deduper filters duplicate log lines within a sliding window.
type Deduper struct {
	opts   Options
	mu     sync.Mutex
	window []uint64
	pos    int
	seen   map[uint64]struct{}
	last   uint64
	count  int64
}

// New creates a new Deduper with the given options.
func New(opts Options) *Deduper {
	if opts.WindowSize <= 0 {
		opts.WindowSize = DefaultOptions().WindowSize
	}
	return &Deduper{
		opts:   opts,
		window: make([]uint64, opts.WindowSize),
		seen:   make(map[uint64]struct{}, opts.WindowSize),
	}
}

// IsDuplicate returns true if the line has been seen within the current window.
func (d *Deduper) IsDuplicate(line string) bool {
	h := hash(line)
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.opts.Consecutive {
		if h == d.last {
			d.count++
			return true
		}
		d.last = h
		return false
	}

	if _, ok := d.seen[h]; ok {
		d.count++
		return true
	}

	// Evict oldest entry when window is full.
	evict := d.window[d.pos]
	if evict != 0 {
		delete(d.seen, evict)
	}
	d.window[d.pos] = h
	d.seen[h] = struct{}{}
	d.pos = (d.pos + 1) % d.opts.WindowSize
	return false
}

// Suppressed returns the total number of lines suppressed so far.
func (d *Deduper) Suppressed() int64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.count
}

// Reset clears all state.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.window = make([]uint64, d.opts.WindowSize)
	d.seen = make(map[uint64]struct{}, d.opts.WindowSize)
	d.pos = 0
	d.last = 0
	d.count = 0
}

func hash(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
