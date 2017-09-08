package ratelimiter

import (
	"errors"
	"math"
	"sync"
	"time"
)

type RateLimiter struct {
	rate          float64
	interval      float64
	maxPermits    float64
	storedPermits float64
	nextFree      int64
	mut           sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	m := &RateLimiter{}
	return m
}

func (m *RateLimiter) SetRate(rate float64) error {
	if rate <= 0.0 {
		return errors.New("RateLimiter: Rate must be greater than 0")
	}

	m.mut.Lock()
	defer m.mut.Unlock()

	m.rate = rate
	m.maxPermits = rate
	m.storedPermits = m.maxPermits
	m.interval = 1000000.0 / rate

	return nil
}

func (m *RateLimiter) GetRate() float64 {
	return m.rate
}

func (m *RateLimiter) Aquire(permits float64) error {
	if permits < 0 || permits > m.maxPermits {
		return errors.New("Ratelimite: Permits must be greater than 0, and smaller than maxPermit")
	}

	now := m.nowMicroSecond()

	m.mut.Lock()
	wait := m.claimNext(float64(permits), now)
	m.mut.Unlock()

	if wait > 0 {
		time.Sleep(time.Duration(wait * 1000))
	}

	return nil
}

func (m *RateLimiter) TryAcquire(permits float64, timeout int64) bool {
	if permits < 0 || permits > m.maxPermits {
		return false
	}

	now := m.nowMicroSecond()

	m.mut.Lock()
	if m.storedPermits+m.futurePermits(now+timeout*1000) < permits {
		m.mut.Unlock()
		return false
	}
	wait := m.claimNext(float64(permits), now)
	m.mut.Unlock()

	if wait > 0 {
		time.Sleep(time.Duration(wait * 1000))
	}

	return true
}

func (m *RateLimiter) futurePermits(future int64) float64 {
	return float64(future-m.nextFree) / m.interval
}

func (m *RateLimiter) sync(now int64) {
	if now > m.nextFree {
		m.storedPermits = math.Min(m.maxPermits, m.storedPermits+m.futurePermits(now))
		m.nextFree = now
	}
}

func (m *RateLimiter) nowMicroSecond() int64 {
	return time.Now().UnixNano() / 1000
}

func (m *RateLimiter) claimNext(permits float64, now int64) int64 {
	m.sync(now)

	stored := math.Min(permits, m.storedPermits)
	fresh := permits - stored

	m.nextFree += int64(fresh * m.interval)
	m.storedPermits -= stored

	wait := m.nextFree - now

	return wait
}
