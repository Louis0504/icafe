package utils

import (
	"fmt"
	"reflect"
)

func RecursiveIndirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

func RecursiveIndirectType(p reflect.Type) reflect.Type {
	for p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	return p
}

func PanicTypeMissmatch(msg string, left, right reflect.Type) {
	if left.Kind() != right.Kind() {
		panic(fmt.Sprintf("%s: %v != %v", msg, left, right))
	}
}
