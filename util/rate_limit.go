package util

import (
	"github.com/juju/ratelimit"
	"time"
)

var RateLimiters []RateLimiter

type RateLimiter struct {
	UserID int
	Bucket *ratelimit.Bucket
}

func NewRateLimiter(userID int) *RateLimiter {
	var rateLimiter RateLimiter
	rateLimiter.UserID = userID
	IntervalForFillingUp := 1 * time.Second
	var capacity int64 = 1
	var quantumPerInterval int64 = 1
	rateLimiter.Bucket = ratelimit.NewBucketWithQuantum(IntervalForFillingUp, capacity, quantumPerInterval)
	return &rateLimiter
}

func FirstRateLimiter(userID int) *RateLimiter {
	for k, RateLimiter := range RateLimiters {
		if RateLimiter.UserID == userID {
			return &RateLimiters[k]
		}
	}
	return NewRateLimiter(userID)
}
