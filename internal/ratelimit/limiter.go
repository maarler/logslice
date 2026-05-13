// Package ratelimit provides a token-bucket rate limiter for controlling
// the throughput of log lines processed per second.
package ratelimit

import (
	"context"
	"time"
)

// Limiter controls the rate at which log lines are processed.
type Limiter struct {
	tokens   chan struct{}
	rate     int
	quit     chan struct{}
}

// New creates a Limiter that allows up to linesPerSecond lines per second.
// If linesPerSecond is zero or negative, no rate limiting is applied.
func New(linesPerSecond int) *Limiter {
	l := &Limiter{
		rate: linesPerSecond,
		quit: make(chan struct{}),
	}
	if linesPerSecond > 0 {
		l.tokens = make(chan struct{}, linesPerSecond)
		go l.refill()
	}
	return l
}

// Wait blocks until a token is available or ctx is cancelled.
// Returns ctx.Err() if the context is cancelled before a token is acquired.
// If no rate limit is configured, Wait returns immediately.
func (l *Limiter) Wait(ctx context.Context) error {
	if l.tokens == nil {
		return nil
	}
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close stops the background refill goroutine.
func (l *Limiter) Close() {
	if l.tokens != nil {
		close(l.quit)
	}
}

// Rate returns the configured lines-per-second limit. Zero means unlimited.
func (l *Limiter) Rate() int {
	return l.rate
}

// refill adds tokens to the bucket at the configured rate using a ticker.
func (l *Limiter) refill() {
	interval := time.Second / time.Duration(l.rate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default:
				// bucket full, discard token
			}
		case <-l.quit:
			return
		}
	}
}
