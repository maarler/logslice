package output

import (
	"fmt"
	"strings"
	"sync"
)

// Namer generates output filenames based on a pattern and a time/segment key.
type Namer struct {
	pattern string
	ext     string
	mu      sync.Mutex
	counts  map[string]int
}

const defaultPattern = "{time}"

// NewNamer creates a Namer with the given filename pattern and extension.
// If pattern is empty, defaultPattern is used.
func NewNamer(pattern, ext string) *Namer {
	if pattern == "" {
		pattern = defaultPattern
	}
	if ext == "" {
		ext = "log"
	}
	return &Namer{
		pattern: pattern,
		ext:     ext,
		counts:  make(map[string]int),
	}
}

// Generate produces a filename for the given segment key.
// Each call with the same key increments an internal index.
func (n *Namer) Generate(key string) string {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.counts[key]++
	idx := n.counts[key]

	name := n.pattern
	name = strings.ReplaceAll(name, "{time}", sanitize(key))
	name = strings.ReplaceAll(name, "{index}", fmt.Sprintf("%03d", idx))

	// If pattern contained no placeholders, append the key directly.
	if name == n.pattern {
		name = sanitize(key)
	}

	return name + "." + n.ext
}

// sanitize replaces characters unsafe for filenames with underscores.
func sanitize(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '/', '\\', ':', ' ', '\t', '\n', '\r', '*', '?', '"', '<', '>', '|':
			b.WriteRune('_')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
