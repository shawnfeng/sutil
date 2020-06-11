// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"errors"
	"fmt"
	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MQType int

const (
	MQTypeKafka MQType = iota
	MQTypeDelay
)

func (t MQType) String() string {
	switch t {
	case MQTypeKafka:
		return "kafka"
	case MQTypeDelay:
		return "delay"
	default:
		return ""
	}
}

type ConfigerType int

const (
	ConfigerTypeSimple ConfigerType = iota
	ConfigerTypeEtcd
	ConfigerTypeApollo
)

func (c ConfigerType) String() string {
	switch c {
	case ConfigerTypeSimple:
		return "simple"
	case ConfigerTypeEtcd:
		return "etcd"
	case ConfigerTypeApollo:
		return "apollo"
	default:
		return "unknown"
	}
}

const (
	defaultTimeout = 3 * time.Second
	//默认1000毫秒
	defaultBatchTimeoutMs = 1000
	defaultTTR            = 3600      // 1 hour
	defaultTTL            = 3600 * 24 // 1 day
	defaultTries          = 1
)

type Config struct {
	MQType         MQType
	MQAddr         []string
	Topic          string
	TimeOut        time.Duration
	CommitInterval time.Duration
	// time interval to flush msg to broker default is 1 second
	BatchTimeout time.Duration
	Offset       int64
	OffsetAt     string
	TTR          uint32 // time to run
	TTL          uint32 // time to live
	Tries        uint16 // delay tries
	BatchSize    int
}

type KeyParts struct {
	Topic string
	Group string
}

var DefaultConfiger Configer

type Configer interface {
	Init(ctx context.Context) error
	GetConfig(ctx context.Context, topic string, mqType MQType) (*Config, error)
	ParseKey(ctx context.Context, k string) (*KeyParts, error)
	Watch(ctx context.Context) <-chan *center.ChangeEvent
}

func NewConfiger(configType ConfigerType) (Configer, error) {
	switch configType {
	case ConfigerTypeSimple:
		return NewSimpleConfiger(), nil
	case ConfigerTypeEtcd:
		return NewEtcdConfiger(), nil
	case ConfigerTypeApollo:
		return NewApolloConfiger(), nil
	default:
		return nil, fmt.Errorf("configType %d error", configType)
	}
}

type SimpleConfig struct {
	mqAddr []string
}

func NewSimpleConfiger() *SimpleConfig {
	return &SimpleConfig{
		mqAddr: []string{"prod.kafka1.ibanyu.com:9092", "prod.kafka2.ibanyu.com:9092", "prod.kafka3.ibanyu.com:9092"},
	}
}

func (m *SimpleConfig) Init(ctx context.Context) error {
	fun := "SimpleConfig.Init-->"
	slog.Infof(ctx, "%s start", fun)
	// noop
	return nil
}

func (m *SimpleConfig) GetConfig(ctx context.Context, topic string, mqType MQType) (*Config, error) {
	fun := "SimpleConfig.GetConfig-->"
	slog.Infof(ctx, "%s get simple config topic:%s", fun, topic)

	return &Config{
		MQType:         mqType,
		MQAddr:         m.mqAddr,
		Topic:          topic,
		TimeOut:        defaultTimeout,
		CommitInterval: 1 * time.Second,
		Offset:         FirstOffset,
	}, nil
}

func (m *SimpleConfig) ParseKey(ctx context.Context, k string) (*KeyParts, error) {
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

func NewEtcdConfiger() *EtcdConfig {
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

func (m *EtcdConfig) GetConfig(ctx context.Context, topic string, mqType MQType) (*Config, error) {
	fun := "EtcdConfig.GetConfig-->"
	slog.Infof(ctx, "%s get etcd config topic:%s", fun, topic)
	// TODO
	return nil, fmt.Errorf("%s etcd config not supported", fun)
}

func (m *EtcdConfig) ParseKey(ctx context.Context, k string) (*KeyParts, error) {
	fun := "EtcdConfig.ParseKey-->"
	return nil, fmt.Errorf("%s not implemented", fun)
}

func (m *EtcdConfig) Watch(ctx context.Context) <-chan *center.ChangeEvent {
	fun := "EtcdConfig.Watch-->"
	slog.Infof(ctx, "%s start", fun)
	// TODO:
	return nil
}

const (
	apolloConfigSep         = "."
	apolloBrokersSep        = ","
	apolloBrokersKey        = "brokers"
	apolloOffsetAtKey       = "offsetat"
	apolloTTRKey            = "ttr"
	apolloTTLKey            = "ttl"
	apolloTriesKey          = "tries"
	apolloBatchSizeKey      = "batchsize"
	apolloBatchTimeoutMsKey = "batchtimeoutms"
)

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

func (m *ApolloConfig) Init(ctx context.Context) (err error) {
	fun := "ApolloConfig.Init-->"
	slog.Infof(ctx, "%s start", fun)
	apolloCenter, err := center.NewConfigCenter(center.ApolloConfigCenter)
	if err != nil {
		slog.Errorf(ctx, "%s create config center err:%v", fun, err)
	}

	err = apolloCenter.Init(ctx, center.DefaultApolloMiddlewareService, []string{center.DefaultApolloMQNamespace})
	if err != nil {
		slog.Errorf(ctx, "%s init config center err:%v", fun, err)
	}

	m.center = apolloCenter
	return
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

func (m *ApolloConfig) getConfigItemWithFallback(ctx context.Context, topic string, name string, mqType MQType) (string, bool) {
	val, ok := m.center.GetStringWithNamespace(ctx, center.DefaultApolloMQNamespace, m.buildKey(ctx, topic, name, mqType))
	if !ok {
		defaultCtx := context.WithValue(ctx, scontext.ContextKeyControl, simpleContextControlRouter{defaultRouteGroup})
		val, ok = m.center.GetStringWithNamespace(defaultCtx, center.DefaultApolloMQNamespace, m.buildKey(defaultCtx, topic, name, mqType))
	}
	return val, ok
}

func (m *ApolloConfig) GetConfig(ctx context.Context, topic string, mqType MQType) (*Config, error) {
	fun := "ApolloConfig.GetConfig-->"
	slog.Infof(ctx, "%s get mq config topic:%s", fun, topic)

	brokersVal, ok := m.getConfigItemWithFallback(ctx, topic, apolloBrokersKey, mqType)
	if !ok {
		return nil, fmt.Errorf("%s no brokers config found", fun)
	}

	var brokers []string
	for _, broker := range strings.Split(brokersVal, apolloBrokersSep) {
		if broker != "" {
			brokers = append(brokers, strings.TrimSpace(broker))
		}
	}

	slog.Infof(ctx, "%s got config brokers:%s", fun, brokers)

	offsetAtVal, ok := m.getConfigItemWithFallback(ctx, topic, apolloOffsetAtKey, mqType)
	if !ok {
		slog.Infof(ctx, "%s no offsetAtVal config founds", fun)

	}
	slog.Infof(ctx, "%s got config offsetAt:%s", fun, offsetAtVal)

	ttrVal, ok := m.getConfigItemWithFallback(ctx, topic, apolloTTRKey, mqType)
	if !ok {
		slog.Infof(ctx, "%s no ttrVal config founds", fun)
	}
	ttr, err := strconv.ParseUint(ttrVal, 10, 32)
	if err != nil {
		ttr = defaultTTR
	}
	slog.Infof(ctx, "%s got config TTR:%d", fun, ttr)

	ttlVal, ok := m.getConfigItemWithFallback(ctx, topic, apolloTTLKey, mqType)
	if !ok {
		slog.Infof(ctx, "%s no ttlVal config founds", fun)
	}
	ttl, err := strconv.ParseUint(ttlVal, 10, 32)
	if err != nil {
		ttl = defaultTTL
	}
	slog.Infof(ctx, "%s got config TTL:%d", fun, ttl)

	triesVal, ok := m.getConfigItemWithFallback(ctx, topic, apolloTriesKey, mqType)
	if !ok {
		slog.Infof(ctx, "%s no triesVal config founds", fun)
	}
	tries, err := strconv.ParseUint(triesVal, 10, 16)
	if err != nil {
		tries = defaultTries
	}
	slog.Infof(ctx, "%s got config triesVal:%s", fun, triesVal)

	batchSize := defaultBatchSize
	batchSizeVal, ok := m.getConfigItemWithFallback(ctx, topic, apolloBatchSizeKey, mqType)
	if !ok {
		// do nothing
		slog.Infof(ctx, "%s has no batchsize config", fun)
	} else {
		t, err := strconv.Atoi(batchSizeVal)
		if err != nil {
			slog.Errorf(ctx, "%s got invalid batchsize config, batchsize: %s", fun, batchSizeVal)
		} else {
			batchSize = t
		}
	}
	slog.Infof(ctx, "%s got config batchSize: %d", fun, batchSize)
	batchTimeoutMs, ok := m.getConfigItemWithFallback(ctx, topic, apolloBatchTimeoutMsKey, mqType)
	if !ok {
		slog.Infof(ctx, "%s no batchTimeout config founds", fun)
	}
	batchTimeoutMsVal, err := strconv.ParseUint(batchTimeoutMs, 10, 32)
	if err != nil {
		batchTimeoutMsVal = defaultBatchTimeoutMs
	}
	slog.Infof(ctx, "%s got config batchTimeout:%d", fun, ttl)

	return &Config{
		MQType:         mqType,
		MQAddr:         brokers,
		Topic:          topic,
		TimeOut:        defaultTimeout,
		CommitInterval: 1 * time.Second,
		BatchTimeout:   time.Duration(batchTimeoutMsVal * 1000000),
		Offset:         FirstOffset,
		OffsetAt:       offsetAtVal,
		TTR:            uint32(ttr),
		TTL:            uint32(ttl),
		Tries:          uint16(tries),
		BatchSize:      batchSize,
	}, nil
}

func (m *ApolloConfig) ParseKey(ctx context.Context, key string) (*KeyParts, error) {
	fun := "ApolloConfig.ParseKey-->"
	parts := strings.Split(key, apolloConfigSep)
	numParts := len(parts)
	if numParts < 4 {
		errMsg := fmt.Sprintf("%s invalid key:%s", fun, key)
		slog.Errorln(ctx, errMsg)
		return nil, errors.New(errMsg)
	}

	return &KeyParts{
		Topic: strings.Join(parts[:numParts-3], apolloConfigSep),
		Group: parts[numParts-3],
	}, nil
}

type apolloObserver struct {
	ch chan<- *center.ChangeEvent
}

func (ob *apolloObserver) HandleChangeEvent(event *center.ChangeEvent) {
	if event.Namespace != center.DefaultApolloMQNamespace {
		return
	}

	// TODO: filter different mq types
	var changes = map[string]*center.Change{}
	for k, ce := range event.Changes {
		if strings.Contains(k, fmt.Sprint(MQTypeKafka)) || strings.Contains(k, fmt.Sprint(MQTypeDelay)) {
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

func (m *ApolloConfig) buildKey(ctx context.Context, topic, item string, mqType MQType) string {
	return strings.Join([]string{
		topic,
		scontext.GetControlRouteGroupWithDefault(ctx, defaultRouteGroup),
		fmt.Sprint(mqType),
		item,
	}, apolloConfigSep)
}
