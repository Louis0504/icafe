package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/Louis0504/icafe/log"
	"io"
	"net"
	"syscall"
	"time"
)

type Conn struct {
	rp      *redis.Pool
	rc      redis.Conn
	slowlog time.Duration
	tun     string // target unit name
	addr    string // target address
}

func (c *Conn) Do(ctx context.Context, cmd string, args ...interface{}) (reply interface{}, err error) {
	reply, err = c.rc.Do(cmd, args...)
	if err != nil && (isConnectionReset(err) || isConnectionEOF(err)) && isIdempotentCommand(cmd) {
		log.Error("retry on connection reset or eof: ", err)
		if closeErr := c.rc.Close(); closeErr != nil {
			log.Error(closeErr)
		}

		c.rc = c.rp.Get()
		reply, err = c.rc.Do(cmd, args...)
	}

	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *Conn) Send(ctx context.Context, cmd string, args ...interface{}) (err error) {
	return c.rc.Send(cmd, args...)
}

func (c *Conn) Flush(ctx context.Context) (err error) {
	return c.rc.Flush()
}

func (c *Conn) Receive(ctx context.Context) (reply interface{}, err error) {

	return c.rc.Receive()
}

func (c *Conn) Err(ctx context.Context) (err error) {

	return c.rc.Err()
}

func (c *Conn) Close(ctx context.Context) error {
	return c.rc.Close()
}

func isConnectionReset(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err.Error() == syscall.ECONNRESET.Error()
	}
	return false
}

func isConnectionEOF(err error) bool {
	return err == io.EOF
}
