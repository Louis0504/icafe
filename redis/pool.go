package redis

import (
	"net/url"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Pool struct {
	rp      *redis.Pool
	slowlog time.Duration
	tun     string // target unit name
	addr    string // target address
}

func NewPool(rawurl string, opts ...Option) *Pool {
	config := &Config{
		maxIdle:         128,
		maxActive:       128,
		wait:            false,
		idleTimeout:     10 * time.Minute,
		connectTimeout:  100 * time.Millisecond,
		writeTimeout:    200 * time.Millisecond,
		readTimeout:     200 * time.Millisecond,
		slowlog:         100 * time.Millisecond,
		maxConnLifetime: 5 * time.Minute,
	}
	for _, o := range opts {
		o(config)
	}

	return &Pool{
		rp: &redis.Pool{
			MaxIdle:         config.maxIdle,
			MaxActive:       config.maxActive,
			IdleTimeout:     config.idleTimeout,
			Wait:            config.wait,
			MaxConnLifetime: config.maxConnLifetime,

			Dial: func() (redis.Conn, error) {
				return redis.DialURL(
					rawurl,
					redis.DialWriteTimeout(config.writeTimeout),
					redis.DialReadTimeout(config.readTimeout),
					redis.DialConnectTimeout(config.connectTimeout),
				)
			},

			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				// Reduce real PING command send.
				if time.Since(t) < 10*time.Second {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		},
		slowlog: config.slowlog,
		tun:     genTargetUnitName(rawurl),
		addr:    rawurl,
	}
}

func (p *Pool) Get() *Conn {
	return &Conn{
		rp:      p.rp,
		rc:      p.rp.Get(),
		slowlog: p.slowlog,
		tun:     p.tun,
		addr:    p.addr,
	}
}

func genTargetUnitName(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "redis"
	}
	return "redis_" + strings.Replace(strings.Replace(u.Host, ".", "_", -1), ":", "_", -1)
}
