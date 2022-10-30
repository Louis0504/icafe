package server

import (
	"context"
	"fmt"
	"github.com/ymetas/icafe/utils"
)

func RunUntilError(ctx context.Context, fns []func(ctx context.Context) error) error {
	for _, fn := range fns {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func RunAll(ctx context.Context, fns []func(ctx context.Context) error) error {
	var eg utils.ErrorGroup
	for _, fn := range fns {
		f := fn
		eg.Go(func() error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("recovered", r)
				}
			}()
			return f(ctx)
		})
	}
	return eg.Wait()
}
