package box

import (
	"context"
)

// 必须用指针这种方式来区别未设置还是0值.
type options struct {
	appName *string
	ctx     context.Context

	withConfig bool

	// pprof port
	profilerPort *int

	includePaths []string

	// lifecycle hooks
	beforeStart []func(ctx context.Context) error
	afterStart  []func(ctx context.Context) error
	beforeStop  []func(ctx context.Context) error
	afterStop   []func(ctx context.Context) error
}

type Option func(*options)

func Name(n string) Option {
	return func(opts *options) {
		opts.appName = &n
	}
}

func WithContext(ctx context.Context) Option {
	return func(opts *options) {
		opts.ctx = ctx
	}
}

func WithConfig() Option {
	return func(opts *options) {
		opts.withConfig = true
	}
}

func WithProfiler(port int) Option {
	return func(opts *options) {
		opts.profilerPort = &port
	}
}

func SentryIncludePaths(paths ...string) Option {
	return func(o *options) {
		o.includePaths = paths
	}
}

func BeforeStart(fn func(ctx context.Context) error) Option {
	return func(opts *options) {
		opts.beforeStart = append(opts.beforeStart, fn)
	}
}

func AfterStart(fn func(ctx context.Context) error) Option {
	return func(opts *options) {
		opts.afterStart = append(opts.afterStart, fn)
	}
}

func BeforeStop(fn func(ctx context.Context) error) Option {
	return func(opts *options) {
		opts.beforeStop = append(opts.beforeStop, fn)
	}
}

func AfterStop(fn func(ctx context.Context) error) Option {
	return func(opts *options) {
		opts.afterStop = append(opts.afterStop, fn)
	}
}
