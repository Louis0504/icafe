package utils

import (
	"errors"
	"strconv"
)

var ErrStrToBool = errors.New("string to bool fail")

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Float64ToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseStrBool only support "true" "1" "false" "0".
func ParseStrBool(s string) (bool, error) {
	switch s {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	default:
		return false, ErrStrToBool
	}
}

func Int64ToBool(i int64) bool {
	return i != 0
}

func BoolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
