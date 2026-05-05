package limiter

import (
	"context"
	"testing"
	"time"
)

type mockClock struct {
	now int64
}

func (c *mockClock) Now() time.Time {
	return time.UnixMilli(c.now)
}

func TestTokenBucket_Tokens(t *testing.T) {
	b := &Bucket{tokens: 0, burst: 10, rate: 2, lastReplenish: 0}
	clock := &mockClock{now: 5000}

	t.Run("happy-path", func(t *testing.T) {
		ctx := context.Background()
		tokens, err := b.Tokens(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(10), tokens)
	})

	t.Run("context-cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := b.Tokens(ctx)
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func TestTokenBucket_Allow(t *testing.T) {
	b := &Bucket{tokens: 0, burst: 10, rate: 2, lastReplenish: 0}
	clock := &mockClock{now: 5000}

	t.Run("happy-path", func(t *testing.T) {
		ctx := context.Background()
		allowed := b.Allow(ctx)
		assert.True(t, allowed)
	})

	t.Run("context-cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		allowed := b.Allow(ctx)
		assert.False(t, allowed)
	})

	t.Run("rate-limit-exhaustion", func(t *testing.T) {
		b.tokens = 0
		ctx := context.Background()
		allowed := b.Allow(ctx)
		assert.False(t, allowed)
	})
}
