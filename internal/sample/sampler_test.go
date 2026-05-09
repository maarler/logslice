package sample

import (
	"bytes"
	"strings"
	"testing"
)

func makeLines(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		sb.WriteString(strings.Repeat("x", 10) + "\n")
	}
	return sb.String()
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n")
}

func TestNew_InvalidRate(t *testing.T) {
	_, err := New(Options{Strategy: StrategyRandom, Rate: 0})
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
	_, err = New(Options{Strategy: StrategyRandom, Rate: 1.5})
	if err == nil {
		t.Fatal("expected error for rate=1.5")
	}
}

func TestNew_InvalidN(t *testing.T) {
	_, err := New(Options{Strategy: StrategyNth, N: 0})
	if err == nil {
		t.Fatal("expected error for N=0")
	}
}

func TestSampler_NthStrategy(t *testing.T) {
	s, err := New(Options{Strategy: StrategyNth, N: 3})
	if err != nil {
		t.Fatal(err)
	}
	input := makeLines(9) // lines 1-9; every 3rd: 3,6,9 => 3 lines
	var out bytes.Buffer
	read, written, err := s.Sample(strings.NewReader(input), &out)
	if err != nil {
		t.Fatal(err)
	}
	if read != 9 {
		t.Errorf("read: want 9, got %d", read)
	}
	if written != 3 {
		t.Errorf("written: want 3, got %d", written)
	}
	if countLines(out.String()) != 3 {
		t.Errorf("output lines: want 3, got %d", countLines(out.String()))
	}
}

func TestSampler_RandomStrategy_KeepsAll(t *testing.T) {
	s, err := New(Options{Strategy: StrategyRandom, Rate: 1.0, Seed: 42})
	if err != nil {
		t.Fatal(err)
	}
	input := makeLines(100)
	var out bytes.Buffer
	read, written, err := s.Sample(strings.NewReader(input), &out)
	if err != nil {
		t.Fatal(err)
	}
	if read != 100 {
		t.Errorf("read: want 100, got %d", read)
	}
	if written != 100 {
		t.Errorf("written: want 100, got %d", written)
	}
}

func TestSampler_EmptyInput(t *testing.T) {
	s, err := New(Options{Strategy: StrategyNth, N: 2})
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	read, written, err := s.Sample(strings.NewReader(""), &out)
	if err != nil {
		t.Fatal(err)
	}
	if read != 0 || written != 0 {
		t.Errorf("want 0/0, got %d/%d", read, written)
	}
}

func TestSampler_RandomStrategy_ReducesLines(t *testing.T) {
	s, err := New(Options{Strategy: StrategyRandom, Rate: 0.1, Seed: 1})
	if err != nil {
		t.Fatal(err)
	}
	input := makeLines(1000)
	var out bytes.Buffer
	_, written, err := s.Sample(strings.NewReader(input), &out)
	if err != nil {
		t.Fatal(err)
	}
	// With rate=0.1 and 1000 lines expect roughly 100; allow wide margin
	if written < 50 || written > 200 {
		t.Errorf("written out of expected range [50,200]: %d", written)
	}
}
