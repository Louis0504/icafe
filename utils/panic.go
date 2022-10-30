package utils

import pkgErrors "github.com/pkg/errors"

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIfWithStack(err error) {
	if err != nil {
		panic(pkgErrors.WithStack(err))
	}
}
