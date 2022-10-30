package utils

import "context"

func DerefCtx(a, b context.Context) context.Context {
	if a != nil {
		return a
	}
	return b
}

func DerefString(a *string, b string) string {
	if a != nil {
		return *a
	}
	return b
}

func DerefBool(a *bool, b bool) bool {
	if a != nil {
		return *a
	}
	return b
}

func DerefInt(a *int, b int) int {
	if a != nil {
		return *a
	}
	return b
}
