package output

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Writer manages output files for log segments.
type Writer struct {
	dir     string
	prefix  string
	files   map[string]*os.File
	DryRun  bool
	Stdout  io.Writer
}

// New creates a new Writer that writes segments to dir.
func New(dir, prefix string) *Writer {
	return &Writer{
		dir:    dir,
		prefix: prefix,
		files:  make(map[string]*os.File),
		Stdout: os.Stdout,
	}
}

// Write writes a line to the segment file identified by key.
// In DryRun mode lines are written to Stdout instead.
func (w *Writer) Write(key, line string) error {
	if w.DryRun {
		_, err := fmt.Fprintf(w.Stdout, "[%s] %s\n", key, line)
		return err
	}

	f, err := w.file(key)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f, line)
	return err
}

// Close flushes and closes all open segment files.
func (w *Writer) Close() error {
	var firstErr error
	for _, f := range w.files {
		if err := f.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	w.files = make(map[string]*os.File)
	return firstErr
}

// Keys returns the segment keys that have been written to.
func (w *Writer) Keys() []string {
	keys := make([]string, 0, len(w.files))
	for k := range w.files {
		keys = append(keys, k)
	}
	return keys
}

func (w *Writer) file(key string) (*os.File, error) {
	if f, ok := w.files[key]; ok {
		return f, nil
	}
	if err := os.MkdirAll(w.dir, 0o755); err != nil {
		return nil, fmt.Errorf("output: create dir %q: %w", w.dir, err)
	}
	name := filepath.Join(w.dir, fmt.Sprintf("%s%s.log", w.prefix, key))
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("output: open %q: %w", name, err)
	}
	w.files[key] = f
	return f, nil
}
