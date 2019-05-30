// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	client *redis.Client
}

func NewCache(ctx context.Context, namespace, prefix string, expire time.Duration, load LoadFunc) (*Cache, error) {
	fun := "NewCache -->"

	client, err := redis.NewClient(ctx, namespace)
	if err != nil {
		slog.Errorf("%s NewClient, err:%s", fun, err)
		return nil, err
	}

	return &Cache{
		namespace: namespace,
		prefix:    prefix,
		load:      load,
		client:    client,
		expire:    expire,
	}, nil
}

func (m *Cache) Get(ctx context.Context, key, value interface{}) error {
	fun := "Cache.Get -->"

	err := m.getValueFromCache(key, value)
	if err == nil {
		return nil
	}
	if err != nil && err.Error() != redis.RedisNil {
		slog.Errorf("%s cache key: %s err: %s", fun, key, err)
		return fmt.Errorf("%s cache key: %s err: %s", fun, key, err)
	}

	slog.Infof("%s miss key: %v, err: %s", fun, key, err)

	return m.loadValueToCache(key)
}

func (m *Cache) Del(ctx context.Context, key interface{}) error {
	fun := "Cache.Del -->"

	skey, err := m.fixKey(key)
	if err != nil {
		slog.Errorf("%s fixkey, key: %v err: %s", fun, key, err)
		return err
	}

	err = m.client.Del(skey).Err()
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

func (m *Cache) getValueFromCache(key, value interface{}) error {

	skey, err := m.fixKey(key)
	if err != nil {
		return err
	}

	data, err := m.client.Get(skey).Bytes()
	if err != nil {
		return err
	}

	slog.Infof("key: %v data: %s", key, string(data))

	err = json.Unmarshal(data, value)
	if err != nil {
		return err
	}

	return nil
}

func (m *Cache) loadValueToCache(key interface{}) error {
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

	rerr := m.client.Set(skey, data, m.expire*time.Second).Err()
	if rerr != nil {
		slog.Errorf("%s set err, cache key: %v rerr: %s", fun, key, rerr)
	}

	if err != nil {
		return err
	}

	return rerr
}
