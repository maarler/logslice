package dedupe

import "io"

// LineSource is anything that yields lines one at a time.
type LineSource interface {
	// ReadLine returns the next line and whether more lines follow.
	ReadLine() (string, bool)
}

// Filter wraps a Deduper and exposes helper functions for common use-cases.
type Filter struct {
	d *Deduper
}

// NewFilter creates a Filter backed by a new Deduper.
func NewFilter(opts Options) *Filter {
	return &Filter{d: New(opts)}
}

// Apply reads all lines from r, writes non-duplicate lines to w, and returns
// the number of lines suppressed.
func (f *Filter) Apply(r io.Reader, w io.Writer) (suppressed int64, err error) {
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 1)
	var line []byte

	for {
		n, rerr := r.Read(tmp)
		if n > 0 {
			if tmp[0] == '\n' {
				text := string(line)
				if !f.d.IsDuplicate(text) {
					buf = append(buf[:0], line...)
					buf = append(buf, '\n')
					if _, werr := w.Write(buf); werr != nil {
						return f.d.Suppressed(), werr
					}
				}
				line = line[:0]
			} else {
				line = append(line, tmp[0])
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return f.d.Suppressed(), rerr
		}
	}
	// Handle final line without trailing newline.
	if len(line) > 0 {
		if !f.d.IsDuplicate(string(line)) {
			if _, werr := w.Write(append(line, '\n')); werr != nil {
				return f.d.Suppressed(), werr
			}
		}
	}
	return f.d.Suppressed(), nil
}

// ApplyLines filters a slice of strings and returns unique lines.
func (f *Filter) ApplyLines(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if !f.d.IsDuplicate(l) {
			out = append(out, l)
		}
	}
	return out
}
