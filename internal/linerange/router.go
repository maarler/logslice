package linerange

import (
	"fmt"
	"io"
	"strings"
)

// Router applies a set of named line ranges to an input stream,
// routing matching lines to the corresponding output writer.
type Router struct {
	routes []route
}

type route struct {
	name  string
	r     Range
	w     io.Writer
}

// NewRouter constructs a Router from a map of name -> range spec.
// Each spec is parsed via Parse; an error is returned on the first
// invalid spec.
func NewRouter(specs map[string]string, writers map[string]io.Writer) (*Router, error) {
	rt := &Router{}
	for name, spec := range specs {
		r, err := Parse(spec)
		if err != nil {
			return nil, fmt.Errorf("linerange router: spec %q for %q: %w", spec, name, err)
		}
		w, ok := writers[name]
		if !ok {
			return nil, fmt.Errorf("linerange router: no writer for %q", name)
		}
		rt.routes = append(rt.routes, route{name: name, r: r, w: w})
	}
	return rt, nil
}

// Route reads lines from src and writes each line to every writer whose
// range contains that line number (1-based). Lines not matched by any
// route are discarded.
func (rt *Router) Route(src io.Reader) error {
	data, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("linerange router: read: %w", err)
	}
	lines := strings.Split(string(data), "\n")
	// Remove trailing empty element produced by a final newline.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	for i, line := range lines {
		lineNo := i + 1
		for _, ro := range rt.routes {
			if ro.r.Contains(lineNo) {
				if _, werr := fmt.Fprintln(ro.w, line); werr != nil {
					return fmt.Errorf("linerange router: write to %q: %w", ro.name, werr)
				}
			}
		}
	}
	return nil
}
