package config

import (
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Mode != SplitByTime {
		t.Errorf("expected default mode %q, got %q", SplitByTime, cfg.Mode)
	}
	if cfg.Window != time.Hour {
		t.Errorf("expected default window 1h, got %v", cfg.Window)
	}
	if cfg.OutputDir != "." {
		t.Errorf("unexpected default output dir: %q", cfg.OutputDir)
	}
}

func TestValidate_MissingInput(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing input file")
	}
}

func TestValidate_TimeMode(t *testing.T) {
	cfg := Default()
	cfg.InputFile = "app.log"
	cfg.Window = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero window in time mode")
	}
}

func TestValidate_LinesMode(t *testing.T) {
	cfg := Default()
	cfg.InputFile = "app.log"
	cfg.Mode = SplitByLines
	cfg.LineLimit = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero line-limit in lines mode")
	}
	cfg.LineLimit = 500
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_PatternMode(t *testing.T) {
	cfg := Default()
	cfg.InputFile = "app.log"
	cfg.Mode = SplitByPattern
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when no include patterns provided")
	}
	cfg.IncludePatterns = []string{"ERROR"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
