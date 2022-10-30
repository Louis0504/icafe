package redis

import (
	"context"
)

// conn.go

func (c *Conn) MustDo(ctx context.Context, cmd string, args ...interface{}) interface{} {
	reply, err := c.Do(ctx, cmd, args...)
	if err != nil {
		panic(err)
	}

	return reply
}

func (c *Conn) MustSend(ctx context.Context, cmd string, args ...interface{}) {
	if err := c.Send(ctx, cmd, args...); err != nil {
		panic(err)
	}
}

func (c *Conn) MustFlush(ctx context.Context) {
	if err := c.Flush(ctx); err != nil {
		panic(err)
	}
}

func (c *Conn) MustReceive(ctx context.Context) interface{} {
	reply, err := c.Receive(ctx)
	if err != nil {
		panic(err)
	}

	return reply
}

func (c *Conn) MustErr(ctx context.Context) {
	if err := c.Err(ctx); err != nil {
		panic(err)
	}
}

func (c *Conn) MustClose(ctx context.Context) {
	if err := c.Close(ctx); err != nil {
		panic(err)
	}
}

// reply.go

func MustInt(reply interface{}, err error) int {
	v, err := Int(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustInt64(reply interface{}, err error) int64 {
	v, err := Int64(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustUint64(reply interface{}, err error) uint64 {
	v, err := Uint64(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustFloat64(reply interface{}, err error) float64 {
	v, err := Float64(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustString(reply interface{}, err error) string {
	v, err := String(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustBytes(reply interface{}, err error) []byte {
	v, err := Bytes(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustBool(reply interface{}, err error) bool {
	v, err := Bool(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustValues(reply interface{}, err error) []interface{} {
	v, err := Values(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustStrings(reply interface{}, err error) []string {
	v, err := Strings(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustByteSlices(reply interface{}, err error) [][]byte {
	v, err := ByteSlices(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustInts(reply interface{}, err error) []int {
	v, err := Ints(reply, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustStringMap(result interface{}, err error) map[string]string {
	v, err := StringMap(result, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustIntMap(result interface{}, err error) map[string]int {
	v, err := IntMap(result, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustInt64Map(result interface{}, err error) map[string]int64 {
	v, err := Int64Map(result, err)
	if err != nil {
		panic(err)
	}
	return v
}

func MustPositions(result interface{}, err error) []*[2]float64 {
	v, err := Positions(result, err)
	if err != nil {
		panic(err)
	}
	return v
}

// scan.go

func MustScan(src []interface{}, dest ...interface{}) []interface{} {
	v, err := Scan(src, dest...)
	if err != nil {
		panic(err)
	}
	return v
}

func MustScanSlice(src []interface{}, dest interface{}, fieldNames ...string) {
	if err := ScanSlice(src, dest, fieldNames...); err != nil {
		panic(err)
	}
}

func MustScanStruct(src []interface{}, dest interface{}) {
	if err := ScanStruct(src, dest); err != nil {
		panic(err)
	}
}
