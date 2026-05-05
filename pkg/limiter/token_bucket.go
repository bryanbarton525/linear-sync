package limiter

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrBucketFull  = errors.New("bucket is full")
	ErrBucketEmpty = errors.New("bucket is empty")
)

// Bucket represents a token bucket rate limiter.
type Bucket struct {
	tokens        int64
	burst         int64
	rate          int64
	lastReplenish int64
	mutex         sync.Mutex
}

// Tokens returns the current number of tokens in the bucket.
func (b *Bucket) Tokens(ctx context.Context) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()

	now := time.Now().UnixMilli()
	elapsed := now - b.lastReplenish
	tokensToAdd := elapsed * b.rate / 1000
	b.tokens = min(b.burst, b.tokens+tokensToAdd)
	b.lastReplenish = now
	return b.tokens, nil
}

// Allow checks if a request can be allowed and consumes a token.
func (b *Bucket) Allow(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()

	now := time.Now().UnixMilli()
	elapsed := now - b.lastReplenish
	tokensToAdd := elapsed * b.rate / 1000
	b.tokens = min(b.burst, b.tokens+tokensToAdd)
	b.lastReplenish = now
	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
