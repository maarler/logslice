package splitter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// Options configures the log splitting behavior.
type Options struct {
	WindowSize  time.Duration
	OutputDir   string
	Prefix      string
	TimeFormat  string
}

// Splitter splits a log file into time-windowed segments.
type Splitter struct {
	opts Options
}

// New creates a new Splitter with the given options.
func New(opts Options) *Splitter {
	if opts.OutputDir == "" {
		opts.OutputDir = "."
	}
	if opts.Prefix == "" {
		opts.Prefix = "segment"
	}
	return &Splitter{opts: opts}
}

// Split reads from r and writes segmented output files based on time windows.
func (s *Splitter) Split(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)

	var (
		currentWindow time.Time
		currentFile   *os.File
		writer        *bufio.Writer
		outFiles      []string
	)

	for scanner.Scan() {
		line := scanner.Text()
		ts, err := parser.ParseTimestampWithFormat(line, s.opts.TimeFormat)
		if err != nil {
			if writer != nil {
				fmt.Fprintln(writer, line)
			}
			continue
		}

		window := ts.Truncate(s.opts.WindowSize)
		if window != currentWindow {
			if writer != nil {
				writer.Flush()
				currentFile.Close()
			}
			currentWindow = window
			name := filepath.Join(s.opts.OutputDir,
				fmt.Sprintf("%s_%s.log", s.opts.Prefix, window.UTC().Format("20060102T150405Z")))
			currentFile, err = os.Create(name)
			if err != nil {
				return outFiles, fmt.Errorf("create segment file: %w", err)
			}
			writer = bufio.NewWriter(currentFile)
			outFiles = append(outFiles, name)
		}
		fmt.Fprintln(writer, line)
	}

	if writer != nil {
		writer.Flush()
		currentFile.Close()
	}
	return outFiles, scanner.Err()
}
