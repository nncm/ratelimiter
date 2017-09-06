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
	m.interval = 1000000.0 / rate

	return nil
}

func (m *RateLimiter) GetRate() float64 {
	return m.rate
}

func (m *RateLimiter) Aquire(permits int) error {
	if permits < 0 {
		return errors.New("Ratelimite: Permits must be greater than 0")
	}

	wait := m.claimNext(float64(permits))

	if wait > 0 {
		time.Sleep(time.Duration(wait * 1000))
	}

	return nil
}

func (m *RateLimiter) TryAcquire(permits int, timeout int) bool {
	if permits < 0 {
		return false
	}

	nowTime := time.Now()
	now := nowTime.Unix()*1e6 + int64(nowTime.Nanosecond())/1e3
	if m.nextFree > now+int64(timeout)*1000 {
		return false
	}

	m.Aquire(permits)

	return true
}

func (m *RateLimiter) sync(now int64) {
	if now > m.nextFree {
		m.storedPermits = math.Min(m.maxPermits, m.storedPermits+float64(now-m.nextFree)/m.interval)
		m.nextFree = now
	}
}

func (m *RateLimiter) claimNext(permits float64) int64 {
	m.mut.Lock()
	defer m.mut.Unlock()

	nowTime := time.Now()
	now := nowTime.Unix()*1e6 + int64(nowTime.Nanosecond())/1e3

	m.sync(now)

	wait := m.nextFree - now

	stored := math.Min(permits, m.storedPermits)
	fresh := permits - stored

	m.nextFree += int64(fresh * m.interval)
	m.storedPermits -= stored

	return wait
}
