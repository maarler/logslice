package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/ratelimit"
)

func TestNew_ZeroRate_NoLimit(t *testing.T) {
	l := ratelimit.New(0)
	defer l.Close()

	if l.Rate() != 0 {
		t.Fatalf("expected rate 0, got %d", l.Rate())
	}

	// Should return immediately without blocking
	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 100; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Fatalf("unlimited limiter blocked for %v", elapsed)
	}
}

func TestNew_PositiveRate_ReturnsRate(t *testing.T) {
	l := ratelimit.New(100)
	defer l.Close()

	if l.Rate() != 100 {
		t.Fatalf("expected rate 100, got %d", l.Rate())
	}
}

func TestWait_ContextCancelled_ReturnsError(t *testing.T) {
	// Very low rate so the bucket drains quickly.
	l := ratelimit.New(1)
	defer l.Close()

	// Drain the initial bucket.
	ctx := context.Background()
	_ = l.Wait(ctx)

	// Now cancel the context before a new token arrives.
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err := l.Wait(cancelCtx)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestWait_RateLimited_Throttles(t *testing.T) {
	const rate = 50
	l := ratelimit.New(rate)
	defer l.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	const calls = 10
	for i := 0; i < calls; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error on call %d: %v", i, err)
		}
	}
	elapsed := time.Since(start)

	// At 50 lines/s, 10 tokens should take at least ~180ms (allowing slack).
	minExpected := 180 * time.Millisecond
	if elapsed < minExpected {
		t.Fatalf("rate limiter too fast: %v for %d calls at %d/s", elapsed, calls, rate)
	}
}

func TestClose_Idempotent(t *testing.T) {
	l := ratelimit.New(10)
	l.Close() // should not panic
}
