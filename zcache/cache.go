package zcache

import (
	"context"
	"time"
)

type KeyFunc func() string
type KeyMultiFunc func(interface{}) string

type FallbackFunc func() (interface{}, error)
type FallbackMultiFunc func(interface{}) (interface{}, error)

type ZCacher interface {
	Get(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, dst interface{}, ttl *time.Duration) error
	MustGet(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, dst interface{}, ttl *time.Duration)

	GetMulti(ctx context.Context, ids interface{},
		keyFunc KeyMultiFunc,
		fallbackFunc FallbackMultiFunc,
		dstMap interface{},
		ttl *time.Duration) error
	MustGetMulti(ctx context.Context, ids interface{},
		keyFunc KeyMultiFunc,
		fallbackFunc FallbackMultiFunc,
		dstMap interface{},
		ttl *time.Duration)

	Evict(ctx context.Context, keyFunc KeyFunc) error
	MustEvict(ctx context.Context, keyFunc KeyFunc)

	EvictMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc) error
	MustEvictMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc)

	Refresh(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, i interface{}, ttl *time.Duration) error
	MustRefresh(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, i interface{}, ttl *time.Duration)

	RefreshMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc, fallbackFunc FallbackMultiFunc, i interface{}, ttl *time.Duration) error
	MustRefreshMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc, fallbackFunc FallbackMultiFunc, i interface{}, ttl *time.Duration)
}
