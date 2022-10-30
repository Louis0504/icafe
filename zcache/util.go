package zcache

import (
	"fmt"
	"reflect"
)

func recursiveIndirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

func recursiveIndirectType(p reflect.Type) reflect.Type {
	for p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	return p
}

func panicTypeError(msg string, left, right reflect.Type) {
	panic(fmt.Sprintf("%s: %v != %v", msg, left, right))
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
