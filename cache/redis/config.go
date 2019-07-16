// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/cache"
	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"strings"
	"sync"
	"time"
)

const (
	defaultPoolSize          = 128
	defaultTimeoutNumSeconds = 1
)

type Config struct {
	addr      string
	namespace string
	poolSize  int
	timeout   time.Duration
}

var DefaultConfiger = NewSimpleConfiger()

type Configer interface {
	GetConfig(ctx context.Context, namespace string) (*Config, error)
	Watch(ctx context.Context) <-chan *center.ChangeEvent
}

func NewConfiger(configType cache.ConfigerType) (Configer, error) {
	switch configType {
	case cache.ConfigerTypeSimple:
		return NewSimpleConfiger(), nil
	case cache.ConfigerTypeEtcd:
		return NewEtcdConfiger(), nil
	case cache.ConfigerTypeApollo:
		return NewApolloConfiger(), nil
	default:
		return nil, fmt.Errorf("configType %d error", configType)
	}
}

type SimpleConfig struct {
}

func NewSimpleConfiger() Configer {
	return &SimpleConfig{}
}

func (m *SimpleConfig) GetConfig(ctx context.Context, namespace string) (*Config, error) {
	addr := ""
	if namespace == "base/report" {
		addr = "common.codis.pri.ibanyu.com:19000"
		//addr = "core.codis.pri.ibanyu.com:19000"
	}
	return &Config{
		addr:      addr,
		namespace: namespace,
		timeout:   defaultTimeoutNumSeconds * time.Second,
		poolSize:  defaultPoolSize,
	}, nil
}

func (m *SimpleConfig) Watch(ctx context.Context) <-chan *center.ChangeEvent {
	fun := "SimpleConfig.Watch-->"
	slog.Infof(ctx, "%s start", fun)
	// noop
	return nil
}

type EtcdConfig struct {
	etcdAddr []string
}

func NewEtcdConfiger() Configer {
	return &EtcdConfig{
		etcdAddr: []string{}, //todo
	}
}

func (m *EtcdConfig) GetConfig(ctx context.Context, namespace string) (*Config, error) {
	fun := "EtcdConfig.GetConfig-->"
	slog.Infof(ctx, "%s get etcd config namespace:%s", fun, namespace)
	//todo etcd router
	return nil, fmt.Errorf("%s etcd config not supported", fun)
}

func (m *EtcdConfig) Watch(ctx context.Context) <-chan *center.ChangeEvent {
	fun := "EtcdConfig.Watch-->"
	slog.Infof(ctx, "%s start", fun)
	// TODO
	return nil
}

const (
	defaultApolloNamespace = "infra.cache"
	apolloConfigSep        = "."

	apolloConfigKeyAddr     = "addr"
	apolloConfigKeyPoolSize = "poolsize"
	apolloConfigKeyTimeout  = "timeout"
)

type ApolloConfig struct {
	watchOnce sync.Once
	ch        chan *center.ChangeEvent
}

func NewApolloConfiger() *ApolloConfig {
	return &ApolloConfig{
		ch: make(chan *center.ChangeEvent),
	}
}

type simpleContextController struct {
	group string
}

func (s simpleContextController) GetGroup() string {
	return s.group
}

func (m *ApolloConfig) getConfigStringItemWithFallback(ctx context.Context, namespace, name string) (string, bool) {
	slog.Infof(ctx, "build key %s", m.buildKey(ctx, namespace, name))
	val, ok := center.GetStringWithNamespace(ctx, defaultApolloNamespace, m.buildKey(ctx, namespace, name))
	if !ok {
		defaultCtx := context.WithValue(ctx, scontext.ContextKeyControl, simpleContextController{defaultGroup})
		val, ok = center.GetStringWithNamespace(defaultCtx, defaultApolloNamespace, m.buildKey(defaultCtx, namespace, name))
	}
	return val, ok
}

func (m *ApolloConfig) getConfigIntItemWithFallback(ctx context.Context, namespace, name string) (int, bool) {
	val, ok := center.GetIntWithNamespace(ctx, defaultApolloNamespace, m.buildKey(ctx, namespace, name))
	if !ok {
		defaultCtx := context.WithValue(ctx, scontext.ContextKeyControl, simpleContextController{defaultGroup})
		val, ok = center.GetIntWithNamespace(defaultCtx, defaultApolloNamespace, m.buildKey(defaultCtx, namespace, name))
	}
	return val, ok
}

func (m *ApolloConfig) GetConfig(ctx context.Context, namespace string) (*Config, error) {
	fun := "ApolloConfig.GetConfig-->"
	slog.Infof(ctx, "%s get apollo config namespace:%s", fun, namespace)

	addr, ok := m.getConfigStringItemWithFallback(ctx, namespace, apolloConfigKeyAddr)
	if !ok {
		return nil, fmt.Errorf("%s no addr config found", fun)
	}
	slog.Infof(ctx, "%s got config addr:%s", fun, addr)

	poolSize, ok := m.getConfigIntItemWithFallback(ctx, namespace, apolloConfigKeyPoolSize)
	if !ok {
		poolSize = defaultPoolSize
		slog.Infof(ctx, "%s no poolSize config found, use default:%d", fun, defaultPoolSize)
	} else {
		slog.Infof(ctx, "%s got config poolSize:%d", fun, poolSize)
	}

	timeout, ok := m.getConfigIntItemWithFallback(ctx, namespace, apolloConfigKeyTimeout)
	if !ok {
		timeout = defaultTimeoutNumSeconds
		slog.Infof(ctx, "%s no timeout config found, use default:%v secs", fun, timeout)
	}
	slog.Infof(ctx, "%s got config timeout:%v seconds", fun, timeout)

	return &Config{
		addr:      addr,
		namespace: namespace,
		poolSize:  poolSize,
		timeout:   time.Duration(timeout) * time.Second,
	}, nil
}

func (m *ApolloConfig) Watch(ctx context.Context) <-chan *center.ChangeEvent {
	// TODO
	return nil
}

func (m *ApolloConfig) buildKey(ctx context.Context, namespace, item string) string {
	return strings.Join([]string{
		namespace,
		scontext.GetGroupWithDefault(ctx, defaultGroup),
		fmt.Sprint(cache.CacheTypeRedis),
		item,
	}, apolloConfigSep)
}
