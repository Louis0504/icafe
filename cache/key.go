package cache

import (
	"bytes"
	"fmt"
)

// NewKey GenerateKey is a helper method for generate consistent key through provided args.
// Attention:
// GenerateKey("fn", map[string]int{"a": 1, "b": 2})
// is not always equal to
// GenerateKey("fn", map[string]int{"b": 2, "a": 1})
// So don't pass a map, you should flatten the map by yourself.
func NewKey(funcName string, args ...interface{}) (cacheKey string) {
	buf := bytes.NewBufferString(funcName)
	for _, arg := range args {
		buf.WriteString(fmt.Sprintf("|%v", arg))
	}
	return buf.String()
}
