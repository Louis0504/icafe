package utils

import (
	"strconv"
	"strings"
)

func JoinInt64(sep string, values ...int64) string {
	if len(values) == 0 {
		return ""
	}

	s := make([]string, len(values))
	for idx, val := range values {
		s[idx] = strconv.FormatInt(val, 10)
	}
	return strings.Join(s, sep)
}

func JoinString(sep string, values ...string) string {
	if len(values) == 0 {
		return ""
	}

	return strings.Join(values, sep)
}

func NilOrEmptyString(strPtr *string) bool {
	return strPtr == nil || *strPtr == ""
}
