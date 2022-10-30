package cache

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/bluele/gcache"
	//"github.com/astaxie/beego"
	"github.com/snowmetas/cafe-go/util"
)

func recursiveIndirectType(p reflect.Type) reflect.Type {
	for p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	return p
}

type LocalCache struct {
	store gcache.Cache
}

func NewLocalStore(size int) *LocalCache {
	return &LocalCache{
		store: gcache.New(size).LRU().Build(),
	}
}

func (ls *LocalCache) Get(_ context.Context, key string, dst interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		panic("cache: dst must be a pointer")
	}

	val, err := ls.store.Get(key)
	if err != nil {
		return err
	}

	dstV := reflect.Indirect(reflect.ValueOf(dst))

	newVal := reflect.Indirect(reflect.ValueOf(val))

	left := recursiveIndirectType(dstV.Type())
	right := recursiveIndirectType(newVal.Type())
	if left.Kind() != right.Kind() {
		panic(fmt.Sprintf("%s: %v != %v", "type of fallback result error", left, right))
	}

	dstV.Set(newVal)

	return nil
}

func (ls *LocalCache) MustGet(ctx context.Context, key string, dst interface{}) {
	util.PanicIfError(ls.Get(ctx, key, dst))
}

func (ls *LocalCache) GetMulti(ctx context.Context, keys []string, dstMap interface{}) error {
	dstPtrV := reflect.ValueOf(dstMap)
	dstV := reflect.Indirect(dstPtrV)
	if dstV.Kind() != reflect.Map {
		panic("cache: dst must be a map or pointer-to-map")
	}

	// nil map
	if dstPtrV.Kind() != reflect.Ptr && dstV.IsNil() {
		panic("cache: dst must not be a nil map")
	}

	if dstPtrV.Kind() == reflect.Ptr && dstV.IsNil() {
		m := reflect.MakeMap(reflect.MapOf(dstV.Type().Key(), dstV.Type().Elem()))
		dstV.Set(m)
	}

	for i, key := range keys {
		v := reflect.New(dstV.Type().Elem())
		if v.Kind() != reflect.Ptr {
			v = v.Addr()
		}

		val, err := ls.store.Get(key)
		if err != nil {
			continue
		}

		dstV.SetMapIndex(reflect.ValueOf(keys[i]), reflect.ValueOf(val))
	}
	return nil
}

func (ls *LocalCache) MustGetMulti(ctx context.Context, keys []string, dstMap interface{}) {
	util.PanicIfError(ls.GetMulti(ctx, keys, dstMap))
}

func (ls *LocalCache) Exists(_ context.Context, key string) (bool, error) {
	return ls.store.Has(key), nil
}

func (ls *LocalCache) MustExists(ctx context.Context, key string) bool {
	ok, err := ls.Exists(ctx, key)
	if err != nil {
		panic(err)
	}
	return ok
}

func (ls *LocalCache) ExistsMulti(ctx context.Context, keys ...string) ([]bool, error) {
	if len(keys) == 0 {
		return []bool{}, nil
	}

	var results []bool
	for _, key := range keys {
		ok, err := ls.Exists(ctx, key)
		if err != nil {
			return nil, err
		}
		results = append(results, ok)
	}

	return results, nil
}

func (ls *LocalCache) MustExistsMulti(ctx context.Context, keys ...string) []bool {
	ret, err := ls.ExistsMulti(ctx, keys...)
	util.PanicIfError(err)
	return ret
}

func (ls *LocalCache) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	err := ls.store.SetWithExpire(key, value, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (ls *LocalCache) MustSet(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	util.PanicIfError(ls.Set(ctx, key, value, ttl))
}

func (ls *LocalCache) SetMulti(ctx context.Context, keys []string, values interface{}, ttl time.Duration) error {
	srcV := reflect.Indirect(reflect.ValueOf(values))

	if srcV.Kind() != reflect.Slice {
		return errors.New("cache: src must be a slice or pointer-to-slice")
	}

	if srcV.Len() != len(keys) {
		return errors.New("cache: keys and src slices have different length")
	}

	for index, key := range keys {
		v := srcV.Index(index)

		err := ls.Set(ctx, key, v.Interface(), ttl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ls *LocalCache) MustSetMulti(ctx context.Context, keys []string, values interface{}, ttl time.Duration) {
	util.PanicIfError(ls.SetMulti(ctx, keys, values, ttl))
}

func (ls *LocalCache) Delete(_ context.Context, keys ...string) error {
	for _, key := range keys {
		ls.store.Remove(key)
	}

	return nil
}

func (ls *LocalCache) MustDelete(ctx context.Context, keys ...string) {
	err := ls.Delete(ctx, keys...)
	if err != nil {
		panic(err)
	}
}
