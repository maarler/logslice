package config

import (
	"flag"
	"testing"
	"time"
)

func newFS() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestFromFlags_Basic(t *testing.T) {
	fs := newFS()
	cfg, err := FromFlags(fs, []string{
		"-input", "server.log",
		"-mode", "time",
		"-window", "30m",
		"-prefix", "web",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.InputFile != "server.log" {
		t.Errorf("input: got %q", cfg.InputFile)
	}
	if cfg.Window != 30*time.Minute {
		t.Errorf("window: got %v", cfg.Window)
	}
	if cfg.OutputPrefix != "web" {
		t.Errorf("prefix: got %q", cfg.OutputPrefix)
	}
}

func TestFromFlags_MultiInclude(t *testing.T) {
	fs := newFS()
	cfg, err := FromFlags(fs, []string{
		"-input", "app.log",
		"-mode", "pattern",
		"-include", "ERROR",
		"-include", "WARN",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.IncludePatterns) != 2 {
		t.Errorf("expected 2 include patterns, got %d", len(cfg.IncludePatterns))
	}
}

func TestFromFlags_MissingInput(t *testing.T) {
	fs := newFS()
	_, err := FromFlags(fs, []string{"-mode", "time"})
	if err == nil {
		t.Error("expected error for missing -input")
	}
}

func TestFromFlags_InvalidWindow(t *testing.T) {
	fs := newFS()
	_, err := FromFlags(fs, []string{"-input", "app.log", "-window", "notaduration"})
	if err == nil {
		t.Error("expected error for invalid window")
	}
}

func TestFromFlags_DryRun(t *testing.T) {
	fs := newFS()
	cfg, err := FromFlags(fs, []string{"-input", "app.log", "-dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.DryRun {
		t.Error("expected DryRun to be true")
	}
}
