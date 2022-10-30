package utils

import (
	"golang.org/x/sync/errgroup"
)

type ErrorGroup struct {
	errgroup.Group
}

func (g *ErrorGroup) Go(f func() error) {
	g.Group.Go(func() (err error) {
		return SafelyRun(func() {
			err = f()
		})
	})
}
