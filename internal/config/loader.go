package config

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// multiFlag allows a flag to be specified multiple times.
type multiFlag []string

func (m *multiFlag) String() string  { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }

// FromFlags parses os.Args via the provided FlagSet and returns a Config.
// Pass flag.CommandLine for normal CLI usage, or a fresh flag.FlagSet for tests.
func FromFlags(fs *flag.FlagSet, args []string) (*Config, error) {
	cfg := Default()

	var (
		windowStr string
		mode      string
		includes  multiFlag
		excludes  multiFlag
	)

	fs.StringVar(&cfg.InputFile, "input", "", "path to the input log file (required)")
	fs.StringVar(&cfg.OutputDir, "output-dir", cfg.OutputDir, "directory for output slices")
	fs.StringVar(&cfg.OutputPrefix, "prefix", cfg.OutputPrefix, "filename prefix for output slices")
	fs.BoolVar(&cfg.DryRun, "dry-run", false, "print actions without writing files")
	fs.StringVar(&mode, "mode", string(cfg.Mode), "split mode: time|lines|pattern")
	fs.StringVar(&windowStr, "window", cfg.Window.String(), "time window for time-based splitting (e.g. 1h, 30m)")
	fs.IntVar(&cfg.LineLimit, "line-limit", cfg.LineLimit, "lines per slice for line-based splitting")
	fs.StringVar(&cfg.TimestampFormat, "ts-format", "", "custom timestamp format (Go reference time)")
	fs.IntVar(&cfg.TimestampField, "ts-field", cfg.TimestampField, "zero-based field index of the timestamp")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "enable verbose output")
	fs.Var(&includes, "include", "include pattern (may be repeated)")
	fs.Var(&excludes, "exclude", "exclude pattern (may be repeated)")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("config: flag parse error: %w", err)
	}

	cfg.Mode = SplitMode(mode)
	cfg.IncludePatterns = []string(includes)
	cfg.ExcludePatterns = []string(excludes)

	if windowStr != "" {
		d, err := time.ParseDuration(windowStr)
		if err != nil {
			return nil, fmt.Errorf("config: invalid window %q: %w", windowStr, err)
		}
		cfg.Window = d
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}
