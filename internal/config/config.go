package config

import (
	"errors"
	"time"
)

// SplitMode controls how the log file is segmented.
type SplitMode string

const (
	SplitByTime    SplitMode = "time"
	SplitByPattern SplitMode = "pattern"
	SplitByLines   SplitMode = "lines"
)

// Config holds the full runtime configuration for a logslice run.
type Config struct {
	// Input source
	InputFile string

	// Output
	OutputDir    string
	OutputPrefix string
	DryRun       bool

	// Splitting
	Mode      SplitMode
	Window    time.Duration // used when Mode == SplitByTime
	LineLimit int           // used when Mode == SplitByLines

	// Filtering
	IncludePatterns []string
	ExcludePatterns []string

	// Timestamp parsing
	TimestampFormat string
	TimestampField  int // zero-based column index

	// Misc
	Verbose bool
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		OutputDir:    ".",
		OutputPrefix: "slice",
		Mode:         SplitByTime,
		Window:       time.Hour,
		TimestampField: 0,
	}
}

// Validate checks that the configuration is self-consistent.
func (c *Config) Validate() error {
	if c.InputFile == "" {
		return errors.New("config: input file must be specified")
	}
	switch c.Mode {
	case SplitByTime:
		if c.Window <= 0 {
			return errors.New("config: window must be positive for time-based splitting")
		}
	case SplitByLines:
		if c.LineLimit <= 0 {
			return errors.New("config: line-limit must be positive for line-based splitting")
		}
	case SplitByPattern:
		if len(c.IncludePatterns) == 0 {
			return errors.New("config: at least one include pattern required for pattern-based splitting")
		}
	default:
		return errors.New("config: unknown split mode")
	}
	return nil
}
