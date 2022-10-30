package redis

import "time"

type Config struct {
	maxIdle           int
	maxActive         int
	wait              bool
	maxRetryPerSecond int64
	retryOnTimeout    bool
	maxConnLifetime   time.Duration
	idleTimeout       time.Duration
	connectTimeout    time.Duration
	writeTimeout      time.Duration
	readTimeout       time.Duration
	slowlog           time.Duration
}

type Option func(*Config)

// MaxIdle Maximum number of idle connections in the rp.
// Default: 10
func MaxIdle(maxIdle int) Option {
	return Option(func(c *Config) {
		c.maxIdle = maxIdle
	})
}

// MaxActive Maximum number of connections allocated by the rp at a given time.
// When zero, there is no limit on the number of connections in the rp.
// Default: 100
func MaxActive(maxActive int) Option {
	return Option(func(c *Config) {
		c.maxActive = maxActive
	})
}

// Wait If wait is true and the rp is at the maxActive limit, then Get() waits
// for a connection to be returned to the rp before returning.
// Default: true
func Wait(wait bool) Option {
	return Option(func(c *Config) {
		c.wait = wait
	})
}

// MaxConnLifetime Close connections after maxConnLifeTime.
// If the value is zero, it doesn't work.
func MaxConnLifetime(maxConnLifeTime time.Duration) Option {
	return Option(func(c *Config) {
		c.maxConnLifetime = maxConnLifeTime
	})
}

// IdleTimeout Close connections after remaining idle for this duration. If the value
// is zero, then idle connections are not closed. Applications should set
// the timeout to a value less than the server's timeout.
// Default: 240s
func IdleTimeout(idleTimeout time.Duration) Option {
	return Option(func(c *Config) {
		c.idleTimeout = idleTimeout
	})
}

// Slowlog set the threshold for redis query.
// Default: 100ms
func Slowlog(threshold time.Duration) Option {
	return Option(func(c *Config) {
		c.slowlog = threshold
	})
}

// ConnectTimeout set the timeout of connecting to redis.
// Default: 100ms
func ConnectTimeout(d time.Duration) Option {
	return Option(func(c *Config) {
		c.connectTimeout = d
	})
}

// WriteTimeout set the timeout of write request to redis.
// Default: 200ms
func WriteTimeout(d time.Duration) Option {
	return Option(func(c *Config) {
		c.writeTimeout = d
	})
}

// ReadTimeout set the timeout of read response from redis.
// Default: 200ms
func ReadTimeout(d time.Duration) Option {
	return Option(func(c *Config) {
		c.readTimeout = d
	})
}

// MaxRetryPerSecond set max retry times per second for RWRedis
// Default: 5
func MaxRetryPerSecond(maxRetryPerSecond int64) Option {
	return Option(func(c *Config) {
		c.maxRetryPerSecond = maxRetryPerSecond
	})
}

// RetryOnTimeout set retry on timeout for RWRedis
// Default: false
// Description: 开启此设置的服务，必须保证对 redis 执行的所有操作为幂等的。其中如有例如 incr 还有 lpush 等非幂等操作，需要评估是否可以开启超时后重试的逻辑
func RetryOnTimeout(c *Config) {
	c.retryOnTimeout = true
}
