package cache

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"encoding/json"
	"errors"
	"github.com/ymetas/icafe/redis"
	"github.com/ymetas/icafe/util"
	"reflect"
	"time"
)

type _compressMode int

const (
	compressModeNone _compressMode = 0 // 不压缩
	compressModeGZip _compressMode = 1 // GZip ：这里其实是一早实现错了， Python 的 gzip 其实是用了 zlib
	compressModeZLib _compressMode = 2 // ZLib ：保持 Python/Golang compress 兼容
)

type RWRedisCache struct {
	rwRedis *redis.RWRedis

	compressMode _compressMode
	escapeHTML   bool
	serializer   Serializer
}

var _ Cache = (*RWRedisCache)(nil)

// NewRWRedisStore NewRedisStore set default behavior as:
// 	compress:   false
//  escapeHTML: true
//  serializer: JSON
func NewRWRedisStore(rwRedis *redis.RWRedis, options ...func(*RWRedisCache)) Cache {
	store := &RWRedisCache{
		rwRedis:      rwRedis,
		compressMode: compressModeNone,
		escapeHTML:   true,
		serializer:   JSON,
	}

	for _, opt := range options {
		opt(store)
	}

	return store
}

// RWStoreCompressMode CompressMode determine whether data should be compressed before store.
func RWStoreCompressMode(compress bool) func(*RWRedisCache) {
	return func(store *RWRedisCache) {
		if compress {
			store.compressMode = compressModeGZip
		} else {
			store.compressMode = compressModeNone
		}
	}
}

// RWStorePythonCompatibleCompress 保持 Python/Golang compress 兼容
func RWStorePythonCompatibleCompress(compress bool) func(*RWRedisCache) {
	return func(store *RWRedisCache) {
		if compress {
			store.compressMode = compressModeZLib
		} else {
			store.compressMode = compressModeNone
		}
	}
}

func RWStoreEscapeHTML(escapeHTML bool) func(*RWRedisCache) {
	return func(store *RWRedisCache) {
		store.escapeHTML = escapeHTML
	}
}

func RWStoreSerializerType(t Serializer) func(*RWRedisCache) {
	return func(store *RWRedisCache) {
		store.serializer = t
	}
}

func (store *RWRedisCache) Get(ctx context.Context, key string, dst interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		panic("cache: dst must be a pointer")
	}

	value, err := redis.Bytes(store.rwRedis.Do(ctx, "GET", key))
	if err != nil {
		return err
	}

	return store.processOutputData(value, dst)
}

func (store *RWRedisCache) MustGet(ctx context.Context, key string, dst interface{}) {
	util.PanicIfError(store.Get(ctx, key, dst))
}

// GetMulti get the provided keys from redis and store it in dst.
// Because golang has no generic type, so result must be provided in params.
// dst must be a map or pointer-to-map
func (store *RWRedisCache) GetMulti(ctx context.Context, keys []string, dstMap interface{}) error {
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

	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}

	values, err := redis.ByteSlices(store.rwRedis.Do(ctx, "MGET", args...))
	if err != nil {
		return err
	}

	for i, value := range values {
		if value == nil {
			continue
		}

		v := reflect.New(dstV.Type().Elem())
		if v.Kind() != reflect.Ptr {
			v = v.Addr()
		}

		if err := store.processOutputData(value, v.Interface()); err != nil {
			return err
		}

		dstV.SetMapIndex(reflect.ValueOf(keys[i]), v.Elem())
	}

	return nil
}

func (store *RWRedisCache) MustGetMulti(ctx context.Context, keys []string, dstMap interface{}) {
	util.PanicIfError(store.GetMulti(ctx, keys, dstMap))
}

func (store *RWRedisCache) Exists(ctx context.Context, key string) (bool, error) {
	return redis.Bool(store.rwRedis.Do(ctx, "EXISTS", key))
}

func (store *RWRedisCache) MustExists(ctx context.Context, key string) bool {
	ret, err := store.Exists(ctx, key)
	util.PanicIfError(err)
	return ret
}

func (store *RWRedisCache) ExistsMulti(ctx context.Context, keys ...string) ([]bool, error) {
	if len(keys) == 0 {
		return []bool{}, nil
	}

	conn := store.rwRedis.ReadClientConn(ctx)
	defer conn.Close(ctx)

	for _, key := range keys {
		if err := conn.Send(ctx, "EXISTS", key); err != nil {
			return nil, err
		}
	}

	if err := conn.Flush(ctx); err != nil {
		return nil, err
	}

	ret := make([]bool, len(keys))
	for i := range keys {
		var err error
		if ret[i], err = redis.Bool(conn.Receive(ctx)); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (store *RWRedisCache) MustExistsMulti(ctx context.Context, keys ...string) []bool {
	ret, err := store.ExistsMulti(ctx, keys...)
	util.PanicIfError(err)
	return ret
}

func (store *RWRedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	buf, err := store.processInputData(value)
	if err != nil {
		return err
	}

	if _, err := store.rwRedis.Do(ctx, "SETEX", key, int(ttl.Seconds()), buf); err != nil {
		return err
	}

	return nil
}

func (store *RWRedisCache) MustSet(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	util.PanicIfError(store.Set(ctx, key, value, ttl))
}

func (store *RWRedisCache) SetMulti(ctx context.Context, keys []string, values interface{}, ttl time.Duration) error {
	srcV := reflect.Indirect(reflect.ValueOf(values))
	if srcV.Kind() != reflect.Slice {
		return errors.New("gache: src must be a slice or pointer-to-slice")
	}
	if srcV.Len() != len(keys) {
		return errors.New("gache: keys and src slices have different length")
	}

	conn := store.rwRedis.WriteClientConn(ctx)
	defer conn.Close(ctx)

	for index, key := range keys {
		v := srcV.Index(index)
		if v.Kind() != reflect.Ptr {
			v = v.Addr()
		}

		if buf, err := store.processInputData(v.Interface()); err != nil {
			return err
		} else {
			conn.Send(ctx, "SETEX", key, int(ttl.Seconds()), buf)
		}
	}

	if err := conn.Flush(ctx); err != nil {
		return err
	}

	return nil
}

func (store *RWRedisCache) MustSetMulti(ctx context.Context, keys []string, values interface{}, ttl time.Duration) {
	util.PanicIfError(store.SetMulti(ctx, keys, values, ttl))
}

func (store *RWRedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = interface{}(key)
	}

	_, err := store.rwRedis.Do(ctx, "DEL", args...)
	return err
}

func (store *RWRedisCache) MustDelete(ctx context.Context, keys ...string) {
	util.PanicIfError(store.Delete(ctx, keys...))
}

func (store *RWRedisCache) processOutputData(src []byte, dst interface{}) error {
	switch store.compressMode {
	case compressModeGZip:
		decompressor, err := gzip.NewReader(bytes.NewReader(src))
		if err != nil {
			return err
		}
		defer decompressor.Close()

		return json.NewDecoder(decompressor).Decode(dst)

	case compressModeZLib:
		decompressor, err := zlib.NewReader(bytes.NewReader(src))
		if err != nil {
			return err
		}
		defer decompressor.Close()

		return json.NewDecoder(decompressor).Decode(dst)

	case compressModeNone:
		return json.Unmarshal(src, dst)
	default:
		return json.Unmarshal(src, dst)
	}
}

func (store *RWRedisCache) processInputData(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}

	switch store.compressMode {
	case compressModeGZip:
		compressor := gzip.NewWriter(&buf)
		// fixme 这里用 defer close 会有问题，因为 compressor.Close() 里其实还会写一些数据到 buffer 中
		defer compressor.Close()

		encoder := json.NewEncoder(compressor)
		encoder.SetEscapeHTML(store.escapeHTML)
		if err := encoder.Encode(value); err != nil {
			return nil, err
		}

		if err := compressor.Flush(); err != nil {
			return nil, err
		}
	case compressModeZLib:
		compressor := zlib.NewWriter(&buf)

		encoder := json.NewEncoder(compressor)
		encoder.SetEscapeHTML(store.escapeHTML)
		if err := encoder.Encode(value); err != nil {
			compressor.Close()
			return nil, err
		}

		if err := compressor.Flush(); err != nil {
			compressor.Close()
			return nil, err
		}

		compressor.Close()
	default:
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(store.escapeHTML)
		if err := encoder.Encode(value); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
