package redis

import (
	"github.com/gomodule/redigo/redis"
)

func Scan(src []interface{}, dest ...interface{}) ([]interface{}, error) {
	return redis.Scan(src, dest...)
}

func ScanSlice(src []interface{}, dest interface{}, fieldNames ...string) error {
	return redis.ScanSlice(src, dest, fieldNames...)
}

func ScanStruct(src []interface{}, dest interface{}) error {
	return redis.ScanStruct(src, dest)
}
