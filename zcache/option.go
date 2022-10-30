package zcache

type Option func(*ZCache)

// IRateLimiter Interface of rate limiter
// TakeAvailable is a non-blocking function, it takes up to count immediately available tokens from the bucket.
// It returns the number of tokens removed, or zero if there are no available tokens.
type IRateLimiter interface {
	TakeAvailable(count int64) int64
}

func FallbackWhenError() Option {
	return Option(func(m *ZCache) {
		m.fallbackWhenError = true
	})
}

func RateLimiter(limiter IRateLimiter) Option {
	return Option(func(m *ZCache) {
		m.limiter = limiter
	})
}
