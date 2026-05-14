package pipeline_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/pipeline"
)

func TestFromProcessor_TransformsLines(t *testing.T) {
	in := make(chan string, 3)
	out := make(chan string, 3)

	in <- "foo"
	in <- "bar"
	in <- "baz"
	close(in)

	stage := pipeline.FromProcessor(strings.ToUpper)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := stage(ctx, in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	close(out)

	var got []string
	for l := range out {
		got = append(got, l)
	}
	if len(got) != 3 || got[0] != "FOO" {
		t.Fatalf("unexpected output: %v", got)
	}
}

func TestFromProcessor_DropsEmptyResult(t *testing.T) {
	in := make(chan string, 2)
	out := make(chan string, 2)

	in <- "keep"
	in <- "skip"
	close(in)

	stage := pipeline.FromProcessor(func(line string) string {
		if line == "skip" {
			return ""
		}
		return line
	})

	ctx := context.Background()
	if err := stage(ctx, in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	close(out)

	var got []string
	for l := range out {
		got = append(got, l)
	}
	if len(got) != 1 || got[0] != "keep" {
		t.Fatalf("expected [keep], got %v", got)
	}
}

func TestFromProcessor_ContextCancel(t *testing.T) {
	in := make(chan string) // blocks forever
	out := make(chan string, 1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stage := pipeline.FromProcessor(strings.TrimSpace)
	err := stage(ctx, in, out)
	if err == nil {
		t.Fatal("expected error on cancelled context")
	}
}
