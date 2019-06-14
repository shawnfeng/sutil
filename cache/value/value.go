// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/cache/redis"
	"github.com/shawnfeng/sutil/slog"
	"time"
)

// key类型只支持int（包含有无符号，8，16，32，64位）和string
type LoadFunc func(key interface{}) (value interface{}, err error)

type Cache struct {
	namespace string
	prefix    string
	load      LoadFunc
	expire    time.Duration
}

func NewCache(namespace, prefix string, expire time.Duration, load LoadFunc) *Cache {
	return &Cache{
		namespace: namespace,
		prefix:    prefix,
		load:      load,
		expire:    expire,
	}
}

func (m *Cache) Get(ctx context.Context, key, value interface{}) error {
	fun := "Cache.Get -->"

	span, ctx := opentracing.StartSpanFromContext(ctx, "cache.value.Get")
	if span != nil {
		defer span.Finish()
	}

	err := m.getValueFromCache(ctx, key, value)
	if err == nil {
		return nil
	}
	if err != nil && err.Error() != redis.RedisNil {
		slog.Errorf("%s cache key: %s err: %s", fun, key, err)
		return fmt.Errorf("%s cache key: %s err: %s", fun, key, err)
	}

	slog.Infof("%s miss key: %v, err: %s", fun, key, err)

	err = m.loadValueToCache(ctx, key)
	if err != nil {
		slog.Errorf("%s loadValueToCache key: %s err: %s", fun, key, err)
		return err
	}

	//简单处理interface对象构造的问题
	return m.getValueFromCache(ctx, key, value)
}

func (m *Cache) Del(ctx context.Context, key interface{}) error {
	fun := "Cache.Del -->"

	span, ctx := opentracing.StartSpanFromContext(ctx, "cache.value.Del")
	if span != nil {
		defer span.Finish()
	}

	skey, err := m.fixKey(key)
	if err != nil {
		slog.Errorf("%s fixkey, key: %v err: %s", fun, key, err)
		return err
	}

	client := redis.DefaultInstanceManager.GetInstance(ctx, m.namespace)
	if client == nil {
		slog.Errorf("%s get instance err, namespace: %s", fun, m.namespace)
		return fmt.Errorf("get instance err, namespace: %s", m.namespace)
	}

	err = client.Del(skey).Err()
	if err != nil {
		return fmt.Errorf("del cache key: %v err: %s", key, err.Error())
	}

	return nil
}

func (m *Cache) keyToString(key interface{}) (string, error) {
	switch t := key.(type) {
	case string:
		return t, nil
	case int8:
		return fmt.Sprintf("%d", key), nil
	case int16:
		return fmt.Sprintf("%d", key), nil
	case int32:
		return fmt.Sprintf("%d", key), nil
	case int64:
		return fmt.Sprintf("%d", key), nil
	case uint8:
		return fmt.Sprintf("%d", key), nil
	case uint16:
		return fmt.Sprintf("%d", key), nil
	case uint32:
		return fmt.Sprintf("%d", key), nil
	case uint64:
		return fmt.Sprintf("%d", key), nil
	case int:
		return fmt.Sprintf("%d", key), nil
	default:
		return "", errors.New("key err: unsupported type")
	}
}

func (m *Cache) fixKey(key interface{}) (string, error) {
	fun := "Cache.fixKey -->"

	skey, err := m.keyToString(key)
	if err != nil {
		slog.Errorf("%s key: %v err:%s", fun, key, err)
		return "", err
	}

	if len(m.prefix) > 0 {
		return fmt.Sprintf("%s.%s", m.prefix, skey), nil
	}

	return skey, nil
}

func (m *Cache) getValueFromCache(ctx context.Context, key, value interface{}) error {
	fun := "Cache.getValueFromCache -->"

	skey, err := m.fixKey(key)
	if err != nil {
		return err
	}

	client := redis.DefaultInstanceManager.GetInstance(ctx, m.namespace)
	if client == nil {
		slog.Errorf("%s get instance err, namespace: %s", fun, m.namespace)
		return fmt.Errorf("get instance err, namespace: %s", m.namespace)
	}

	data, err := client.Get(skey).Bytes()
	if err != nil {
		return err
	}

	slog.Infof("%s key: %v data: %s", fun, key, string(data))

	err = json.Unmarshal(data, value)
	if err != nil {
		return err
	}

	return nil
}

func (m *Cache) loadValueToCache(ctx context.Context, key interface{}) error {
	fun := "Cache.loadValueToCache -->"

	var data []byte
	value, err := m.load(key)
	if err != nil {
		slog.Warnf("%s load err, cache key:%s err:%s", fun, key, err)
		data = []byte(err.Error())

	} else {
		data, err = json.Marshal(value)
		if err != nil {
			slog.Errorf("%s marshal err, cache key: %s err:%s", fun, key, err)
			data = []byte(err.Error())
		}
	}

	skey, err := m.fixKey(key)
	if err != nil {
		slog.Errorf("%s fixkey, key: %s err:%s", fun, key, err)
		return err
	}

	client := redis.DefaultInstanceManager.GetInstance(ctx, m.namespace)
	if client == nil {
		slog.Errorf("%s get instance err, namespace: %s", fun, m.namespace)
		return fmt.Errorf("get instance err, namespace: %s", m.namespace)
	}

	rerr := client.Set(skey, data, m.expire*time.Second).Err()
	if rerr != nil {
		slog.Errorf("%s set err, cache key: %v rerr: %s", fun, key, rerr)
	}

	if err != nil {
		return err
	}

	return rerr
}
