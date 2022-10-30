package redis

import (
	"github.com/gomodule/redigo/redis"
)

var (
	Int        = redis.Int
	Int64      = redis.Int64
	Uint64     = redis.Uint64
	Float64    = redis.Float64
	String     = redis.String
	Bytes      = redis.Bytes
	Bool       = redis.Bool
	MultiBulk  = redis.MultiBulk
	Values     = redis.Values
	Float64s   = redis.Float64s
	Strings    = redis.Strings
	ByteSlices = redis.ByteSlices
	Ints       = redis.Ints
	Int64s     = redis.Int64s
	StringMap  = redis.StringMap
	IntMap     = redis.IntMap
	Int64Map   = redis.Int64Map
	Positions  = redis.Positions
)
