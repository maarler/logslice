// Package checkpoint persists the last successfully processed byte offset
// for a log file so that logslice can resume after an interruption.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// State holds the persisted position for a single source file.
type State struct {
	Path   string    `json:"path"`
	Offset int64     `json:"offset"`
	SavedAt time.Time `json:"saved_at"`
}

// Store reads and writes checkpoint state to a JSON file on disk.
type Store struct {
	file string
}

// NewStore returns a Store that persists state in the given file.
func NewStore(file string) *Store {
	return &Store{file: file}
}

// Load reads the persisted State from disk.
// It returns a zero-value State and no error when the file does not exist yet.
func (s *Store) Load() (State, error) {
	data, err := os.ReadFile(s.file)
	if errors.Is(err, os.ErrNotExist) {
		return State{}, nil
	}
	if err != nil {
		return State{}, err
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return State{}, err
	}
	return st, nil
}

// Save persists st to disk, creating parent directories as needed.
func (s *Store) Save(st State) error {
	st.SavedAt = time.Now().UTC()
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.file), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.file, data, 0o644)
}

// Delete removes the checkpoint file if it exists.
func (s *Store) Delete() error {
	err := os.Remove(s.file)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
