package server

import (
	"context"
	"github.com/YLeseclaireurs/icafe/utils"
	"sync"

	"github.com/YLeseclaireurs/icafe/log"
)

type Bundle interface {
	Type() string
	Name() string
	Run(ctx context.Context) error
	Stop() context.Context
}

type Container struct {
	bundles  []Bundle
	bundleWg sync.WaitGroup
}

func New() *Container {
	return &Container{}
}

func bundleDesc(b Bundle) string {
	return b.Type() + "[" + b.Name() + "]"
}

func (c *Container) AddBundle(bundles ...Bundle) {
	c.bundles = append(c.bundles, bundles...)
}

func (c *Container) StartAll(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	for _, b := range c.bundles {
		bundle := b
		log.Info(ctx, "Start bundle:", bundleDesc(bundle))
		c.bundleWg.Add(1)
		go func() {
			defer c.bundleWg.Done()
			if err := bundle.Run(ctx); err != nil {
				log.Errorf("Run bundle:%s failed error: %s", bundleDesc(bundle), err.Error())
			}
		}()
		log.Info(ctx, "Bundle started:", bundleDesc(bundle))
	}

	go func() {
		c.bundleWg.Wait()
		cancel()
	}()

	return ctx
}

func (c *Container) StopAll(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	var eg utils.ErrorGroup
	for _, b := range c.bundles {
		bundle := b
		log.Info(ctx, "Stop bundle:", bundleDesc(bundle))
		eg.Go(func() error {
			stopCtx := bundle.Stop()
			<-stopCtx.Done()
			log.Info(ctx, "Bundle stopped:", bundleDesc(bundle))
			return nil
		})
	}

	go func() {
		if err := eg.Wait(); err != nil {
			log.Error(ctx, "Stop bundel error", err)
		}
		cancel()
		log.Info(ctx, "All bundle stopped")
	}()

	return ctx
}
