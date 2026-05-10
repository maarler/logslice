package mask

import (
	"fmt"
	"strings"
)

// RuleConfig holds the raw configuration for a single masking rule,
// suitable for unmarshalling from flags or a config file.
type RuleConfig struct {
	Label       string `json:"label"       yaml:"label"`
	Pattern     string `json:"pattern"     yaml:"pattern"`
	Replacement string `json:"replacement" yaml:"replacement"`
}

// Options groups all masking-related configuration.
type Options struct {
	// UsePreset loads the built-in sensitive-data rules (IP, email, tokens).
	UsePreset bool
	// Rules are additional user-defined masking rules.
	Rules []RuleConfig
}

// Build constructs a Masker from the Options.
// Returns nil, nil when masking is entirely disabled (no preset, no rules).
func (o *Options) Build() (*Masker, error) {
	var all []*Rule

	if o.UsePreset {
		pm, err := Preset()
		if err != nil {
			return nil, fmt.Errorf("mask preset: %w", err)
		}
		all = append(all, pm.rules...)
	}

	for _, rc := range o.Rules {
		if strings.TrimSpace(rc.Pattern) == "" {
			return nil, fmt.Errorf("mask rule %q: pattern must not be empty", rc.Label)
		}
		repl := rc.Replacement
		if repl == "" {
			repl = "[REDACTED]"
		}
		r, err := NewRule(rc.Label, rc.Pattern, repl)
		if err != nil {
			return nil, fmt.Errorf("mask rule %q: %w", rc.Label, err)
		}
		all = append(all, r)
	}

	if len(all) == 0 {
		return nil, nil
	}
	return New(all...), nil
}

// DefaultOptions returns an Options with the preset enabled and no custom rules.
func DefaultOptions() Options {
	return Options{UsePreset: false}
}
