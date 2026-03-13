package limiter

import "sync"

type IPRateLimiter struct {
	buckets    sync.Map
	capacity   int
	refillRate int
}

func NewIPRateLimiter(capacity, refillRate int) *IPRateLimiter {
	return &IPRateLimiter{
		capacity:   capacity,
		refillRate: refillRate,
	}
}

func (rl *IPRateLimiter) getBucket(ip string) *TokenBucket {
	if bucket, ok := rl.buckets.Load(ip); ok {
		return bucket.(*TokenBucket)
	}

	tb := newTokenBucket(rl.capacity, rl.refillRate)
	actual, _ := rl.buckets.LoadOrStore(ip, tb)
	return actual.(*TokenBucket)
}

func (rl *IPRateLimiter) Allow(ip string) bool {
	return rl.getBucket(ip).allow()
}
