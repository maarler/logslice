package filter

import (
	"strings"
	"testing"
)

func TestChain_EmptyMatchesAll(t *testing.T) {
	c := NewChain()
	lines := []string{"INFO hello", "ERROR boom", "DEBUG trace"}
	for _, l := range lines {
		if !c.Match(l) {
			t.Errorf("empty chain should match %q", l)
		}
	}
}

func TestChain_SinglePattern(t *testing.T) {
	c := NewChain()
	if err := c.AddExpr(`ERROR`, false); err != nil {
		t.Fatal(err)
	}
	if !c.Match("ERROR: disk full") {
		t.Error("expected match")
	}
	if c.Match("INFO: all fine") {
		t.Error("expected no match")
	}
}

func TestChain_MultiplePatterns(t *testing.T) {
	c := NewChain()
	_ = c.AddExpr(`ERROR`, false)
	_ = c.AddExpr(`disk`, false)

	if !c.Match("ERROR: disk full") {
		t.Error("expected match for line with ERROR and disk")
	}
	if c.Match("ERROR: network timeout") {
		t.Error("expected no match: missing 'disk'")
	}
}

func TestApplyLines(t *testing.T) {
	c := NewChain()
	_ = c.AddExpr(`ERROR`, false)

	input := []string{"INFO start", "ERROR boom", "ERROR crash", "WARN low mem"}
	res := ApplyLines(input, c)

	if len(res.Lines) != 2 {
		t.Errorf("expected 2 matched lines, got %d", len(res.Lines))
	}
	if res.Dropped != 2 {
		t.Errorf("expected 2 dropped lines, got %d", res.Dropped)
	}
}

func TestApply_Reader(t *testing.T) {
	c := NewChain()
	_ = c.AddExpr(`WARN`, false)

	input := "INFO ok\nWARN low disk\nERROR fail\nWARN high load\n"
	res, err := Apply(strings.NewReader(input), c)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(res.Lines))
	}
}
