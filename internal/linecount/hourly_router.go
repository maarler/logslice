package linecount

import (
	"fmt"
	"io"
	"time"
)

// HourlyRouter routes log lines to per-hour writers.
// Writers are created on demand via the Factory function.
type HourlyRouter struct {
	factory func(hour time.Time) (io.WriteCloser, error)
	writers map[string]io.WriteCloser
}

// NewHourlyRouter creates a router that calls factory whenever a new hour
// bucket is encountered. The caller must invoke Close to flush all writers.
func NewHourlyRouter(factory func(hour time.Time) (io.WriteCloser, error)) *HourlyRouter {
	return &HourlyRouter{
		factory: factory,
		writers: make(map[string]io.WriteCloser),
	}
}

// Route writes line to the writer corresponding to the given timestamp.
func (r *HourlyRouter) Route(t time.Time, line string) error {
	k := hourKey(t)
	w, ok := r.writers[k]
	if !ok {
		var err error
		w, err = r.factory(t.UTC().Truncate(time.Hour))
		if err != nil {
			return fmt.Errorf("hourly router: open writer for %s: %w", k, err)
		}
		r.writers[k] = w
	}
	_, err := fmt.Fprintln(w, line)
	return err
}

// Keys returns the hour keys that have received at least one line.
func (r *HourlyRouter) Keys() []string {
	keys := make([]string, 0, len(r.writers))
	for k := range r.writers {
		keys = append(keys, k)
	}
	return keys
}

// Close closes all open writers.
func (r *HourlyRouter) Close() error {
	var first error
	for k, w := range r.writers {
		if err := w.Close(); err != nil && first == nil {
			first = fmt.Errorf("close writer %s: %w", k, err)
		}
		delete(r.writers, k)
	}
	return first
}
