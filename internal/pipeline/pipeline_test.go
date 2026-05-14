package pipeline_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/pipeline"
)

func feedLines(lines []string) <-chan string {
	ch := make(chan string, len(lines))
	for _, l := range lines {
		ch <- l
	}
	close(ch)
	return ch
}

func collectLines(ch <-chan string) []string {
	var out []string
	for l := range ch {
		out = append(out, l)
	}
	return out
}

func TestPipeline_NoStages_PassThrough(t *testing.T) {
	src := feedLines([]string{"a", "b", "c"})
	dst := make(chan string, 10)

	p := pipeline.New(8)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.Run(ctx, src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	close(dst)
	got := collectLines(dst)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestPipeline_SingleStage_Transforms(t *testing.T) {
	src := feedLines([]string{"hello", "world"})
	dst := make(chan string, 10)

	upper := pipeline.FromProcessor(strings.ToUpper)
	p := pipeline.New(8, upper)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.Run(ctx, src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	close(dst)
	got := collectLines(dst)
	if len(got) != 2 || got[0] != "HELLO" || got[1] != "WORLD" {
		t.Fatalf("unexpected output: %v", got)
	}
}

func TestPipeline_StageDropsLines(t *testing.T) {
	src := feedLines([]string{"keep", "drop", "keep"})
	dst := make(chan string, 10)

	filter := pipeline.FromProcessor(func(line string) string {
		if line == "drop" {
			return ""
		}
		return line
	})
	p := pipeline.New(8, filter)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.Run(ctx, src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	close(dst)
	got := collectLines(dst)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines after drop, got %d: %v", len(got), got)
	}
}

func TestPipeline_ContextCancel(t *testing.T) {
	blocking := make(chan string) // never closed
	dst := make(chan string, 10)

	p := pipeline.New(8, pipeline.FromProcessor(strings.ToUpper))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := p.Run(ctx, blocking, dst)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}
