package timewindow

import (
	"time"
)

// Entry is a timestamped line of text.
type Entry struct {
	Time time.Time
	Line string
}

// Router distributes log entries into per-window buckets.
type Router struct {
	duration time.Duration
	buckets  map[string][]string
	order    []string
}

// NewRouter creates a Router that groups entries into windows of size d.
func NewRouter(d time.Duration) *Router {
	if d <= 0 {
		panic("timewindow: router duration must be positive")
	}
	return &Router{
		duration: d,
		buckets:  make(map[string][]string),
	}
}

// Add routes entry into the appropriate time window bucket.
func (r *Router) Add(e Entry) {
	w := Bucket(e.Time, r.duration)
	k := w.Key()
	if _, exists := r.buckets[k]; !exists {
		r.order = append(r.order, k)
	}
	r.buckets[k] = append(r.buckets[k], e.Line)
}

// Keys returns window keys in the order they were first observed.
func (r *Router) Keys() []string {
	out := make([]string, len(r.order))
	copy(out, r.order)
	return out
}

// Lines returns all lines belonging to the window identified by key.
// Returns nil if the key is unknown.
func (r *Router) Lines(key string) []string {
	lines, ok := r.buckets[key]
	if !ok {
		return nil
	}
	out := make([]string, len(lines))
	copy(out, lines)
	return out
}

// Len returns the total number of entries across all buckets.
func (r *Router) Len() int {
	n := 0
	for _, v := range r.buckets {
		n += len(v)
	}
	return n
}

// Reset clears all buckets and ordering state.
func (r *Router) Reset() {
	r.buckets = make(map[string][]string)
	r.order = r.order[:0]
}
