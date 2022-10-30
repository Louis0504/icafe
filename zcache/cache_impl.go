package zcache

import (
	"context"
	"reflect"
	"time"

	"github.com/bluele/gcache"

	"github.com/garyburd/redigo/redis"
	"github.com/snowmetas/cafe-go/cache"
)

type ZCache struct {
	cache             cache.Cache
	expire            time.Duration // ZCache实例默认的兜底过期时间，如果命令没有设置过期时间，默认用这个兜底时间
	fallbackWhenError bool
	limiter           IRateLimiter
}

var _ ZCacher = (*ZCache)(nil)

func NewZCache(cache cache.Cache, expire time.Duration, opts ...Option) *ZCache {
	r := &ZCache{
		cache:  cache,
		expire: expire,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

func (c *ZCache) canFallbackWhenError() bool {
	if !c.fallbackWhenError {
		return false
	}
	if c.limiter == nil {
		return true
	}
	return c.limiter.TakeAvailable(1) != 0
}

func (c *ZCache) get(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, dst interface{}, expire time.Duration) error {
	dstPtrV := reflect.ValueOf(dst)
	if dstPtrV.Kind() != reflect.Ptr {
		panic("memcache: dst must be a pointer")
	}

	key := keyFunc()
	err := c.cache.Get(ctx, key, dst)
	if err == nil {
		return nil
	} else if err == redis.ErrNil || err == gcache.KeyNotFoundError || c.canFallbackWhenError() {
		fallbackResult, err := fallbackFunc()
		if err != nil {
			return err
		}

		// check nil
		fV := reflect.ValueOf(fallbackResult)
		if !fV.IsValid() || fV.Kind() == reflect.Ptr && fV.IsNil() {
			return nil
		}

		// check fallback type
		dstV := reflect.Indirect(reflect.ValueOf(dst))
		// 此处不需要强一致类型 **struct 和 *struct 是可以的
		left := recursiveIndirectType(dstV.Type())
		right := recursiveIndirectType(fV.Type())
		if left.Kind() != right.Kind() {
			panicTypeError("type of fallback result error", left, right)
		}
		dstV.Set(fV)

		_ = c.cache.Set(ctx, key, dst, expire)
	} else {
		return err
	}

	return nil
}

func (c *ZCache) Get(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, dst interface{}, ttl *time.Duration) error {
	expire := c.expire
	if ttl != nil {
		expire = *ttl
	}
	return c.get(ctx, keyFunc, fallbackFunc, dst, expire)
}

func (c *ZCache) MustGet(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, dst interface{}, ttl *time.Duration) {
	panicError(c.Get(ctx, keyFunc, fallbackFunc, dst, ttl))
}

func (c *ZCache) getMulti(ctx context.Context, ids interface{},
	keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc,
	dstMap interface{}, expire time.Duration) error {

	// dst must be a map or pointer-to-map
	dstPtrV := reflect.ValueOf(dstMap)
	dstV := reflect.Indirect(dstPtrV)
	if dstV.Kind() != reflect.Map {
		panic("memcache: dst must be a map or pointer-to-map")
	}

	// nil map
	if dstPtrV.Kind() != reflect.Ptr && dstV.IsNil() {
		panic("memcache: dst must not be a nil map")
	}

	// auto generate map
	if dstPtrV.Kind() == reflect.Ptr && dstV.IsNil() {
		m := reflect.MakeMap(reflect.MapOf(dstV.Type().Key(), dstV.Type().Elem()))
		dstV.Set(m)
	}

	dstType := recursiveIndirectType(reflect.TypeOf(dstMap))
	mapKeyType := dstType.Key()
	mapValueType := dstType.Elem()

	// check slice type
	idsT := reflect.ValueOf(ids)
	if idsT.Kind() != reflect.Slice {
		panic("memcache: ids must be a slice")
	}

	length := idsT.Len()
	if length == 0 {
		return nil
	}

	// check type of id vs key type of dstMap
	left := idsT.Type().Elem()
	right := mapKeyType
	if left.Kind() != right.Kind() {
		panicTypeError("type of ids elem and key not equal", left, right)
	}

	actualIds := reflect.MakeSlice(reflect.SliceOf(left), length, length)
	for i := 0; i < length; i++ {
		actualIds.Index(i).Set(idsT.Index(i))
	}

	// generate keys
	keys := make([]string, length)
	revertKeyMap := make(map[string]interface{}) // key -> id
	for i := 0; i < length; i++ {
		id := actualIds.Index(i).Interface()
		key := keyFunc(id)
		keys[i] = key
		revertKeyMap[key] = id
	}

	// 这里要处理 *map 的情况
	dstValueT := reflect.TypeOf(dstMap)
	if dstPtrV.Kind() == reflect.Ptr {
		dstValueT = dstValueT.Elem()
	}
	dstValueT = dstValueT.Elem()

	// get from cache
	cacheDstV := reflect.MakeMap(reflect.MapOf(reflect.ValueOf("").Type(), dstValueT))
	err := c.cache.GetMulti(ctx, keys, cacheDstV.Interface())
	if err != nil && !c.canFallbackWhenError() {
		return err
	}

	// check miss
	cacheMissIdsV := reflect.MakeSlice(reflect.SliceOf(mapKeyType), 0, 8)
	for _, key := range keys {
		id := revertKeyMap[key]

		// find value from map
		v := cacheDstV.MapIndex(reflect.ValueOf(key))
		if !v.IsValid() {
			cacheMissIdsV = reflect.Append(cacheMissIdsV, reflect.ValueOf(id))
		} else {
			dstV.SetMapIndex(reflect.ValueOf(id), v)
		}
	}

	// fallback && set cache
	if cacheMissIdsV.Len() == 0 {
		return nil
	}

	fallbackResult, err := fallbackFunc(cacheMissIdsV.Interface())
	if err != nil {
		return err
	}

	// check map type
	fallbackResultV := reflect.ValueOf(fallbackResult)
	if fallbackResultV.Kind() != reflect.Map {
		panic("memcache: type of fallbackResultV must be map")
	}

	if fallbackResultV.Type().Key().Kind() != mapKeyType.Kind() ||
		fallbackResultV.Type().Elem().Kind() != mapValueType.Kind() {
		panic("memcache: key type and value type of fallbackResult is not equal to map")
	}

	if fallbackResultV.Len() == 0 {
		return nil
	}

	cacheKeys := make([]string, 0)
	cacheSrcV := reflect.MakeSlice(reflect.SliceOf(dstValueT), 0, 8)

	left = recursiveIndirectType(recursiveIndirect(reflect.ValueOf(dstMap)).Type().Elem())
	mapKeys := fallbackResultV.MapKeys()

	for _, keyV := range mapKeys {
		valueV := fallbackResultV.MapIndex(keyV)

		if !valueV.IsValid() || valueV.Kind() == reflect.Ptr && valueV.IsNil() {
			continue
		}

		key := keyFunc(keyV.Interface())
		cacheKeys = append(cacheKeys, key)
		cacheSrcV = reflect.Append(cacheSrcV, valueV)
		dstV.SetMapIndex(keyV, valueV)
	}

	if cacheSrcV.Len() > 0 {
		_ = c.cache.SetMulti(ctx, cacheKeys, cacheSrcV.Interface(), expire)
	}

	return nil

}

func (c *ZCache) GetMulti(ctx context.Context, ids interface{},
	keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc,
	dstMap interface{},
	ttl *time.Duration) error {
	expire := c.expire
	if ttl != nil {
		expire = *ttl
	}
	return c.getMulti(ctx, ids, keyFunc, fallbackFunc, dstMap, expire)
}

func (c *ZCache) MustGetMulti(ctx context.Context, ids interface{},
	keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc,
	dstMap interface{},
	ttl *time.Duration) {
	panicError(c.GetMulti(ctx, ids, keyFunc, fallbackFunc, dstMap, ttl))
}

func (c *ZCache) Evict(ctx context.Context, keyFunc KeyFunc) error {
	return c.cache.Delete(ctx, keyFunc())
}

func (c *ZCache) MustEvict(ctx context.Context, keyFunc KeyFunc) {
	panicError(c.Evict(ctx, keyFunc))
}

func (c *ZCache) EvictMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc) error {
	// check slice type
	idsV := reflect.ValueOf(ids)
	if idsV.Kind() != reflect.Slice {
		panic("memcache: ids must be a slice")
	}

	length := idsV.Len()
	if length == 0 {
		return nil
	}

	actualIds := reflect.MakeSlice(reflect.SliceOf(idsV.Type().Elem()), length, length)
	for i := 0; i < length; i++ {
		actualIds.Index(i).Set(idsV.Index(i))
	}

	keys := make([]string, length)
	for i := 0; i < length; i++ {
		keys[i] = keyFunc(actualIds.Index(i).Interface())
	}

	return c.cache.Delete(ctx, keys...)
}

func (c *ZCache) MustEvictMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc) {
	panicError(c.EvictMulti(ctx, ids, keyFunc))
}

func (c *ZCache) refresh(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, p reflect.Type, expire time.Duration) error {
	key := keyFunc()
	err := c.cache.Delete(ctx, key)
	if err != nil {
		return err
	}

	fallbackResult, err := fallbackFunc()
	if err != nil {
		return err
	}

	// check nil
	fV := reflect.ValueOf(fallbackResult)
	if !fV.IsValid() || fV.Kind() == reflect.Ptr && fV.IsNil() {
		return nil
	}

	// 右侧 fallback 可以使用指针
	left := p
	right := reflect.Indirect(fV).Type()
	if left.Kind() != right.Kind() {
		panicTypeError("type of fallback result error", left, right)
	}

	err = c.cache.Set(ctx, key, fallbackResult, expire)
	if err != nil {
		return err
	}

	return nil
}

func (c *ZCache) Refresh(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, i interface{}, ttl *time.Duration) error {
	expire := c.expire
	if ttl != nil {
		expire = *ttl
	}
	return c.refresh(ctx, keyFunc, fallbackFunc, reflect.TypeOf(i), expire)
}

func (c *ZCache) MustRefresh(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, i interface{}, ttl *time.Duration) {
	panicError(c.Refresh(ctx, keyFunc, fallbackFunc, i, ttl))
}

func (c *ZCache) refreshMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc, dst interface{}, expire time.Duration) error {

	dstV := reflect.ValueOf(dst)
	dstT := dstV.Type()
	// check dst
	if dstV.Kind() != reflect.Map {
		panic("dst must be map")
	}
	if dstV.IsNil() {
		panic("dst map must be initialized")
	}
	dstKeyT := dstT.Key()
	dstValT := dstT.Elem()

	// check ids is slice and ids element type equal to map key type
	idsV := reflect.ValueOf(ids)
	idsT := idsV.Type()
	if idsV.Kind() != reflect.Slice {
		panic("ids must be slice")
	}
	length := idsV.Len()
	if length == 0 {
		return nil
	}
	idsValT := idsT.Elem()
	if idsValT != dstKeyT {
		panic("id's element type must equal to dst map's key type")
	}

	// gen key
	keys := make([]string, length)
	for i := 0; i < length; i++ {
		id := idsV.Index(i).Interface()
		keys[i] = keyFunc(id)
	}

	// delete cache
	err := c.cache.Delete(ctx, keys...)
	if err != nil {
		return err
	}

	// fallback
	r, err := fallbackFunc(ids)
	if err != nil {
		return err
	}

	// check fallback with dst
	rV := reflect.ValueOf(r)
	rT := rV.Type()
	if rV.Kind() != reflect.Map {
		panic("fallback's type must be map")
	}
	length = rV.Len()
	if length == 0 {
		return nil
	}
	rKeyT := rT.Key()
	rValT := rT.Elem()
	if rKeyT != dstKeyT {
		panic("fallback's key type must equal to dst's key type")
	}
	if rValT != dstValT {
		panic("fallback's value type must equal to dst's value type")
	}

	keys = make([]string, length)
	values := reflect.MakeSlice(reflect.SliceOf(rValT), 0, length)
	for i, idV := range rV.MapKeys() {
		id := idV.Interface()
		key := keyFunc(id)
		keys[i] = key
		values = reflect.Append(values, rV.MapIndex(idV))
	}
	return c.cache.SetMulti(ctx, keys, values.Interface(), expire)
}

func (c *ZCache) RefreshMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc, dst interface{}, ttl *time.Duration) error {
	expire := c.expire
	if ttl != nil {
		expire = *ttl
	}
	return c.refreshMulti(ctx, ids, keyFunc, fallbackFunc, dst, expire)
}

func (c *ZCache) MustRefreshMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc, dst interface{}, ttl *time.Duration) {
	expire := c.expire
	if ttl != nil {
		expire = *ttl
	}
	panicError(c.refreshMulti(ctx, ids, keyFunc, fallbackFunc, dst, expire))
}
