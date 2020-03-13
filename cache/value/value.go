// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO: 删除旧版本的 cache sdk
// TODO: 将 redis 文件夹中的 config.go, instance.go, redis.go 都拿到 sdk 根目录下
//       config 和 instance 不应该属于 redis

package value

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/cache"
	"github.com/shawnfeng/sutil/cache/constants"
	"github.com/shawnfeng/sutil/cache/redis"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"github.com/shawnfeng/sutil/stime"
	"time"
)

// key类型只支持int（包含有无符号，8，16，32，64位）和string
type LoadFunc func(ctx context.Context, key interface{}) (value interface{}, err error)

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

func (m *Cache) getInstanceConf(ctx context.Context) *redis.InstanceConf {
	return &redis.InstanceConf{
		Group:     scontext.GetControlRouteGroupWithDefault(ctx, constants.DefaultRouteGroup),
		Namespace: m.namespace,
		Wrapper:   cache.WrapperTypeCache,
	}
}

func (m *Cache) Get(ctx context.Context, key, value interface{}) error {
	fun := "Cache.Get -->"
	// TODO 目前统计的是cache层的Get，后面需要拆分为redis层、cache层
	command := "cache.value.Get"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()

	err := m.getValueFromCache(ctx, key, value)
	if err == nil {
		_metricHits.With("namespace", m.namespace, "command", command).Inc()
		return nil
	}

	if err.Error() != redis.RedisNil {
		statReqErr(m.namespace, command, err)
		slog.Errorf(ctx, "%s cache key: %v err: %v", fun, key, err)
		return fmt.Errorf("%s cache key: %v err: %v", fun, key, err)
	}
	_metricMiss.With("namespace", m.namespace, "command", command).Inc()

	data, err := m.loadValueToCache(ctx, key)
	if err != nil {
		statReqErr(m.namespace, command, err)
		slog.Errorf(ctx, "%s loadValueToCache key: %v err: %v", fun, key, err)
		return err
	}

	err = json.Unmarshal(data, value)
	if err != nil {
		statReqErr(m.namespace, command, err)
		return errors.New(string(data))
	}

	return nil
}

func (m *Cache) Del(ctx context.Context, key interface{}) error {
	fun := "Cache.Del -->"
	command := "cache.value.Del"

	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()

	skey, err := m.prefixKey(key)
	if err != nil {
		statReqErr(m.namespace, command, err)
		slog.Errorf(ctx, "%s fixkey, key: %v err: %v", fun, key, err)
		return err
	}

	client, err := redis.DefaultInstanceManager.GetInstance(ctx, m.getInstanceConf(ctx))
	if err != nil {
		statReqErr(m.namespace, command, err)
		slog.Errorf(ctx, "%s get instance err, namespace: %s", fun, m.namespace)
		return err
	}

	err = client.Del(ctx, skey).Err()
	if err != nil {
		statReqErr(m.namespace, command, err)
		return fmt.Errorf("del cache key: %v err: %s", key, err.Error())
	}

	return nil
}

func (m *Cache) Load(ctx context.Context, key interface{}) error {
	command := "cache.value.Load"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()

	_, err := m.loadValueToCache(ctx, key)
	statReqErr(m.namespace, command, err)

	return err
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

func (m *Cache) prefixKey(key interface{}) (string, error) {
	fun := "Cache.prefixKey -->"

	skey, err := m.keyToString(key)
	if err != nil {
		slog.Errorf(context.TODO(), "%s key: %v err:%s", fun, key, err)
		return "", err
	}

	if len(m.prefix) > 0 {
		return fmt.Sprintf("%s.%s", m.prefix, skey), nil
	}

	return skey, nil
}

func (m *Cache) getValueFromCache(ctx context.Context, key, value interface{}) error {
	fun := "Cache.getValueFromCache -->"

	skey, err := m.prefixKey(key)
	if err != nil {
		return err
	}

	client, err := redis.DefaultInstanceManager.GetInstance(ctx, m.getInstanceConf(ctx))
	if err != nil {
		slog.Errorf(ctx, "%s get instance err, namespace: %s", fun, m.namespace)
		return err
	}

	data, err := client.Get(ctx, skey).Bytes()
	if err != nil {
		return err
	}

	//slog.Infof(ctx, "%s key: %v data: %s", fun, key, string(data))

	err = json.Unmarshal(data, value)
	if err != nil {
		return errors.New(string(data))
	}

	return nil
}

func (m *Cache) loadValueToCache(ctx context.Context, key interface{}) (data []byte, err error) {
	fun := "Cache.loadValueToCache -->"
	expire := m.expire

	value, err := m.load(ctx, key)
	if err != nil {
		slog.Warnf(ctx, "%s load err, cache key:%v err:%v", fun, key, err)
		data = []byte(err.Error())
		expire = constants.CacheDirtyExpireTime

	} else {
		data, err = json.Marshal(value)
		if err != nil {
			slog.Errorf(ctx, "%s marshal err, cache key:%v err:%v", fun, key, err)
			data = []byte(err.Error())
			expire = constants.CacheDirtyExpireTime
		}
	}

	skey, err := m.prefixKey(key)
	if err != nil {
		slog.Errorf(ctx, "%s fixkey, key: %v err:%v", fun, key, err)
		return nil, err
	}

	client, err := redis.DefaultInstanceManager.GetInstance(ctx, m.getInstanceConf(ctx))
	if err != nil {
		slog.Errorf(ctx, "%s get instance err, namespace: %s", fun, m.namespace)
		return nil, err
	}

	rerr := client.Set(ctx, skey, data, expire).Err()
	if rerr != nil {
		slog.Errorf(ctx, "%s set err, cache key:%v rerr:%v", fun, key, rerr)
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

func SetConfiger(ctx context.Context, configerType constants.ConfigerType) error {
	fun := "Cache.SetConfiger-->"
	configer, err := redis.NewConfiger(configerType)
	if err != nil {
		slog.Errorf(ctx, "%s create configer err:%v", fun, err)
		return err
	}
	slog.Infof(ctx, "%s %v configer created", fun, configerType)
	err = configer.Init(ctx)
	if err != nil {
		slog.Errorf(ctx, "%s init configer err:%v", fun, err)
	}
	redis.DefaultConfiger = configer
	return err
}

func WatchUpdate(ctx context.Context) {
	go redis.DefaultInstanceManager.Watch(ctx)
}

func init() {
	fun := "value.init -->"
	ctx := context.Background()
	err := SetConfiger(ctx, constants.ConfigerTypeApollo)
	if err != nil {
		slog.Errorf(ctx, "%s set cache configer:%v err:%v", fun, constants.ConfigerTypeApollo, err)
	} else {
		slog.Infof(ctx, "%s cache configer:%v been set", fun, constants.ConfigerTypeApollo)
	}
	WatchUpdate(ctx)
}
