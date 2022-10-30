package redis

import (
	"context"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/juju/ratelimit"
)

type RWRedis struct {
	name string
	writeAddr string
	readAddrs []string
	writePool *Pool
	readPools []*Pool

	count          uint64
	length         uint64
	retryLimiter   *ratelimit.Bucket
	retryOnTimeout bool
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

/*
func init() {
	writeAddr := "redis://x:22465"
	readAddrs := []string{
		"redis://x:22466",
		"redis://x:22467",
		"redis://x:22468",
	}
	r = NewRWRedis(writeAddr, readAddrs, ReadTimeout(1*time.Second))
}
*/

func NewRWRedis(name, writeAddr string, readAddrs []string, opts ...Option) *RWRedis {
	var readPools []*Pool
	if len(readAddrs) == 0 {
		readPools = append(readPools, NewPool(writeAddr, opts...))
	} else {
		for _, addr := range readAddrs {
			readPools = append(readPools, NewPool(addr, opts...))
		}
	}

	c := &Config{
		maxRetryPerSecond: 5,
		retryOnTimeout:    false,
	}
	for _, o := range opts {
		o(c)
	}

	return &RWRedis{
		name: 		name,
		writeAddr: writeAddr,
		readAddrs: readAddrs,
		writePool: NewPool(writeAddr, opts...),
		readPools: readPools,

		count:          rand.Uint64(),
		length:         uint64(len(readPools)),
		retryOnTimeout: c.retryOnTimeout,
		// 每秒 c.maxRetryPerSecond 个令牌，最多 10 个
		retryLimiter: ratelimit.NewBucketWithRate(float64(c.maxRetryPerSecond), 10),
	}
}

func (r *RWRedis) getReadClient(ctx context.Context) *Pool {
	c := atomic.AddUint64(&r.count, 1)
	return r.readPools[c%r.length]
}

func (r *RWRedis) WriteClientConn(ctx context.Context) *Conn {
	return r.writePool.Get()
}

func (r *RWRedis) ReadClientConn(ctx context.Context) *Conn {
	return r.getReadClient(ctx).Get()
}

func (r *RWRedis) ConnWrite(ctx context.Context, execFunc func(context.Context, *Conn)) {
	conn := r.WriteClientConn(ctx)
	defer conn.Close(ctx)
	execFunc(ctx, conn)
}

func (r *RWRedis) ConnRead(ctx context.Context, execFunc func(context.Context, *Conn)) {
	conn := r.ReadClientConn(ctx)
	defer conn.Close(ctx)
	execFunc(ctx, conn)
}

func (r *RWRedis) do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
	var conn *Conn
	if isReadCommand(cmd) {
		conn = r.getReadClient(ctx).Get()
	} else {
		conn = r.writePool.Get()
	}
	defer conn.Close(ctx)

	reply, err := conn.Do(ctx, cmd, args...)
	return reply, err
}

func (r *RWRedis) Do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
	reply, err := r.do(ctx, cmd, args...)
	if _, ok := err.(*net.OpError); ok && r.retryOnTimeout && r.retryLimiter.TakeAvailable(1) > 0 {
		reply, err = r.do(ctx, cmd, args...)
	}
	return reply, err
}

func (r *RWRedis) MustDo(ctx context.Context, cmd string, args ...interface{}) interface{} {
	reply, err := r.Do(ctx, cmd, args...)
	if err != nil {
		panic(err)
	}

	return reply
}
