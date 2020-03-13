// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/cache/constants"
	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"strings"
	"sync"
	"time"
)

const (
	apolloConfigSep = "."

	apolloConfigKeyAddr       = "addr"
	apolloConfigKeyPoolSize   = "poolsize"
	apolloConfigKeyTimeout    = "timeout"
	apolloConfigKeyUseWrapper = "usewrapper"

	defaultPoolSize          = 128
	defaultTimeoutNumSeconds = 3
	defaultUseWrapper        = true
)

type Config struct {
	addr       string
	namespace  string
	poolSize   int
	timeout    time.Duration
	useWrapper bool
}

type KeyParts struct {
	Namespace string
	Group     string
}

var DefaultConfiger Configer

type Configer interface {
	Init(ctx context.Context) error
	GetConfig(ctx context.Context, namespace string) (*Config, error)
	ParseKey(ctx context.Context, key string) (*KeyParts, error)
	Watch(ctx context.Context) <-chan *center.ChangeEvent
}

func NewConfiger(configType constants.ConfigerType) (Configer, error) {
	switch configType {
	case constants.ConfigerTypeSimple:
		return NewSimpleConfiger(), nil
	case constants.ConfigerTypeEtcd:
		return NewEtcdConfiger(), nil
	case constants.ConfigerTypeApollo:
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

func (m *SimpleConfig) Init(ctx context.Context) error {
	fun := "SimpleConfig.Init-->"
	slog.Infof(ctx, "%s start", fun)
	// noop
	return nil
}

func (m *SimpleConfig) GetConfig(ctx context.Context, namespace string) (*Config, error) {
	addr := ""
	if namespace == "base/report" {
		addr = "common.codis.pri.ibanyu.com:19000"
		//addr = "core.codis.pri.ibanyu.com:19000"
	}

	if namespace == "base/growthsystem" {
		addr = "common.codis.pri.ibanyu.com:19000"
	}

	return &Config{
		addr:      addr,
		namespace: namespace,
		timeout:   defaultTimeoutNumSeconds * time.Second,
		poolSize:  defaultPoolSize,
	}, nil
}

func (m *SimpleConfig) ParseKey(ctx context.Context, key string) (*KeyParts, error) {
	fun := "SimpleConfig.ParseKey-->"
	return nil, fmt.Errorf("%s not implemented", fun)
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

func (m *EtcdConfig) Init(ctx context.Context) error {
	fun := "EtcdConfig.Init-->"
	slog.Infof(ctx, "%s start", fun)
	// TODO
	return nil
}

func (m *EtcdConfig) GetConfig(ctx context.Context, namespace string) (*Config, error) {
	fun := "EtcdConfig.GetConfig-->"
	slog.Infof(ctx, "%s get etcd config namespace:%s", fun, namespace)
	//todo etcd router
	return nil, fmt.Errorf("%s etcd config not supported", fun)
}

func (m *EtcdConfig) ParseKey(ctx context.Context, key string) (*KeyParts, error) {
	fun := "EtcdConfig.ParseKey-->"
	return nil, fmt.Errorf("%s not implemented", fun)
}

func (m *EtcdConfig) Watch(ctx context.Context) <-chan *center.ChangeEvent {
	fun := "EtcdConfig.Watch-->"
	slog.Infof(ctx, "%s start", fun)
	// TODO
	return nil
}


type ApolloConfig struct {
	watchOnce sync.Once
	ch        chan *center.ChangeEvent
	center    center.ConfigCenter
}

func NewApolloConfiger() *ApolloConfig {
	return &ApolloConfig{
		ch: make(chan *center.ChangeEvent),
	}
}

func (m *ApolloConfig) Init(ctx context.Context) error {
	fun := "ApolloConfig.Init-->"
	apolloCenter, err := center.NewConfigCenter(center.ApolloConfigCenter)
	if err != nil {
		slog.Errorf(ctx, "%s create config center err:%v", fun, err)
	}

	err = apolloCenter.Init(ctx, center.DefaultApolloMiddlewareService, []string{center.DefaultApolloCacheNamespace})
	if err != nil {
		slog.Errorf(ctx, "%s init config center err:%v", fun, err)
	}

	m.center = apolloCenter
	return err
}

type simpleContextControlRouter struct {
	group string
}

func (s simpleContextControlRouter) GetControlRouteGroup() (string, bool) {
	return s.group, true
}

func (s simpleContextControlRouter) SetControlRouteGroup(group string) error {
	s.group = group
	return nil
}

func (m *ApolloConfig) getConfigStringItemWithFallback(ctx context.Context, namespace, name string) (string, bool) {
	val, ok := m.center.GetStringWithNamespace(ctx, center.DefaultApolloCacheNamespace, m.buildKey(ctx, namespace, name))
	if !ok {
		defaultCtx := context.WithValue(ctx, scontext.ContextKeyControl, simpleContextControlRouter{defaultGroup})
		val, ok = m.center.GetStringWithNamespace(defaultCtx, center.DefaultApolloCacheNamespace, m.buildKey(defaultCtx, namespace, name))
	}
	return val, ok
}

func (m *ApolloConfig) getConfigIntItemWithFallback(ctx context.Context, namespace, name string) (int, bool) {
	val, ok := m.center.GetIntWithNamespace(ctx, center.DefaultApolloCacheNamespace, m.buildKey(ctx, namespace, name))
	if !ok {
		defaultCtx := context.WithValue(ctx, scontext.ContextKeyControl, simpleContextControlRouter{defaultGroup})
		val, ok = m.center.GetIntWithNamespace(defaultCtx, center.DefaultApolloCacheNamespace, m.buildKey(defaultCtx, namespace, name))
	}
	return val, ok
}

func (m *ApolloConfig) getConfigBoolItemWithFallback(ctx context.Context, namespace, name string) (bool, bool) {
	val, ok := m.center.GetBoolWithNamespace(ctx, center.DefaultApolloCacheNamespace, m.buildKey(ctx, namespace, name))
	if !ok {
		defaultCtx := context.WithValue(ctx, scontext.ContextKeyControl, simpleContextControlRouter{defaultGroup})
		val, ok = m.center.GetBoolWithNamespace(defaultCtx, center.DefaultApolloCacheNamespace, m.buildKey(defaultCtx, namespace, name))
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

	useWrapper, ok := m.getConfigBoolItemWithFallback(ctx, namespace, apolloConfigKeyUseWrapper)
	if !ok {
		useWrapper = defaultUseWrapper
		slog.Infof(ctx, "%s no usewrapper config found, use default:%v", fun, useWrapper)
	}
	slog.Infof(ctx, "%s got config usewrapper:%v", fun, useWrapper)

	return &Config{
		addr:       addr,
		namespace:  namespace,
		poolSize:   poolSize,
		timeout:    time.Duration(timeout) * time.Second,
		useWrapper: useWrapper,
	}, nil
}

func (m *ApolloConfig) ParseKey(ctx context.Context, key string) (*KeyParts, error) {
	fun := "ApolloConfig.ParseKey-->"
	parts := strings.Split(key, apolloConfigSep)
	numParts := len(parts)

	if numParts < 4 {
		err := fmt.Errorf("%s invalid key:%s", fun, key)
		slog.Errorf(ctx, "%s err:%v", fun, err)
		return nil, err
	}

	return &KeyParts{
		Namespace: strings.Join(parts[:numParts-3], apolloConfigSep),
		Group:     parts[numParts-3],
	}, nil
}

type apolloObserver struct {
	ch chan<- *center.ChangeEvent
}

func (ob *apolloObserver) HandleChangeEvent(event *center.ChangeEvent) {
	if event.Namespace != center.DefaultApolloCacheNamespace {
		return
	}

	var changes = map[string]*center.Change{}
	for k, ce := range event.Changes {
		if strings.Contains(k, fmt.Sprint(constants.CacheTypeRedis)) {
			changes[k] = ce
		}
	}

	event.Changes = changes
	ob.ch <- event
}

func (m *ApolloConfig) Watch(ctx context.Context) <-chan *center.ChangeEvent {
	fun := "ApolloConfig.Watch-->"
	m.watchOnce.Do(func() {
		slog.Infof(ctx, "%s start", fun)
		m.center.StartWatchUpdate(ctx)
		m.center.RegisterObserver(ctx, &apolloObserver{m.ch})
	})
	return m.ch
}

func (m *ApolloConfig) buildKey(ctx context.Context, namespace, item string) string {
	return strings.Join([]string{
		namespace,
		scontext.GetControlRouteGroupWithDefault(ctx, defaultGroup),
		fmt.Sprint(constants.CacheTypeRedis),
		item,
	}, apolloConfigSep)
}
