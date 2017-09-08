package ratelimiter

import (
	"testing"
	"time"
)

func TestLimiter1(t *testing.T) {
	rl := NewRateLimiter()
	rl.SetRate(1)
	var result bool
	result = rl.TryAcquire(1, 0)
	if !result {
		t.Error("Allow: false, want true")
	}
	result = rl.TryAcquire(1, 0)
	if result {
		t.Error("Allow: true, want false")
	}

	time.Sleep(1 * time.Second)
	result = rl.TryAcquire(1, 0)
	if !result {
		t.Error("Allow: false, want true")
	}
	result = rl.TryAcquire(1, 0)
	if result {
		t.Error("Allow: true, want false")
	}
}
