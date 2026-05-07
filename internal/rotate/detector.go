package rotate

import (
	"os"
	"sync"
)

// FileState holds the last known inode and size of a file.
type FileState struct {
	Inode uint64
	Size  int64
}

// Detector watches a file path and detects rotation (truncation or replacement).
type Detector struct {
	mu    sync.Mutex
	path  string
	state FileState
}

// NewDetector creates a Detector seeded with the current state of path.
// Returns an error if the file cannot be stat'd.
func NewDetector(path string) (*Detector, error) {
	d := &Detector{path: path}
	state, err := statFile(path)
	if err != nil {
		return nil, err
	}
	d.state = state
	return d, nil
}

// Check compares the current file state against the last known state.
// Returns (rotated, error). rotated is true when the file has been replaced
// (different inode) or truncated (size smaller than before).
func (d *Detector) Check() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	current, err := statFile(d.path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}

	rotated := current.Inode != d.state.Inode || current.Size < d.state.Size
	d.state = current
	return rotated, nil
}

// Reset updates the stored state to the current file state.
func (d *Detector) Reset() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	state, err := statFile(d.path)
	if err != nil {
		return err
	}
	d.state = state
	return nil
}
