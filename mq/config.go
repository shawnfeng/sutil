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
	"strings"
	"sync"
	"time"
)

type MQType int

const (
	MqTypeKafka MQType = iota
)

func (t MQType) String() string {
	switch t {
	case MqTypeKafka:
		return "kafka"
	default:
		return ""
	}
}

const (
	ConfigTypeSimple = iota
	ConfigTypeEtcd
	ConfigTypeApollo
)

const (
	defaultTimeout = 3 * time.Second
)

type Config struct {
	MQType         MQType
	MQAddr         []string
	Topic          string
	TimeOut        time.Duration
	CommitInterval time.Duration
	Offset         int64
}

type KeyParts struct {
	Topic string
	Group string
}

var DefaultConfiger = NewSimpleConfiger()

type Configer interface {
	GetConfig(ctx context.Context, topic string) (*Config, error)
	ParseKey(ctx context.Context, k string) (*KeyParts, error)
	Watch(ctx context.Context) <-chan *center.ChangeEvent
}

func NewConfiger(configType int) (Configer, error) {
	switch configType {
	case ConfigTypeSimple:
		return NewSimpleConfiger(), nil
	case ConfigTypeEtcd:
		return NewEtcdConfiger(), nil
	case ConfigTypeApollo:
		return NewApolloConfig(), nil
	default:
		return nil, fmt.Errorf("configType %d error", configType)
	}
}

type SimpleConfig struct {
	mqAddr []string
}

func NewSimpleConfiger() Configer {
	return &SimpleConfig{
		mqAddr: []string{"prod.kafka1.ibanyu.com:9092", "prod.kafka2.ibanyu.com:9092", "prod.kafka3.ibanyu.com:9092"},
	}
}

func (m *SimpleConfig) GetConfig(ctx context.Context, topic string) (*Config, error) {
	fun := "SimpleConfig.GetConfig-->"
	slog.Infof(ctx, "%s get simple config topic:%s", fun, topic)

	return &Config{
		MQType:         MqTypeKafka,
		MQAddr:         m.mqAddr,
		Topic:          topic,
		TimeOut:        defaultTimeout,
		CommitInterval: 1 * time.Second,
		Offset:         FirstOffset,
	}, nil
}

func (m *SimpleConfig) Watch(ctx context.Context) <-chan *center.ChangeEvent {
	// noop
	return nil
}

func (m *SimpleConfig) ParseKey(ctx context.Context, k string) (*KeyParts, error) {
	fun := "SimpleConfig.ParseKey-->"
	return nil, fmt.Errorf("%s not implemented", fun)
}

type EtcdConfig struct {
	etcdAddr []string
}

func NewEtcdConfiger() Configer {
	return &EtcdConfig{
		etcdAddr: []string{}, //todo
	}
}

func (m *EtcdConfig) GetConfig(ctx context.Context, topic string) (*Config, error) {
	fun := "EtcdConfig.GetConfig-->"
	slog.Infof(ctx, "%s get etcd config topic:%s", fun, topic)

	return nil, fmt.Errorf("%s etcd config not supported", fun)
}

func (m *EtcdConfig) Watch(ctx context.Context) <-chan *center.ChangeEvent {
	// TODO:
	return nil
}

func (m *EtcdConfig) ParseKey(ctx context.Context, k string) (*KeyParts, error) {
	fun := "EtcdConfig.ParseKey-->"
	return nil, fmt.Errorf("%s not implemented", fun)
}

const defaultApolloNamespace = "infra.mq"
const apolloConfigSep = "."

type ApolloConfig struct {
	watchOnce sync.Once
	ch        chan *center.ChangeEvent
}

func NewApolloConfig() Configer {
	return &ApolloConfig{
		ch: make(chan *center.ChangeEvent),
	}
}

func (m *ApolloConfig) GetConfig(ctx context.Context, topic string) (*Config, error) {
	fun := "ApolloConfig.GetConfig-->"
	slog.Infof(ctx, "%s get mq config topic:%s", fun, topic)

	brokerKey := m.buildKey(ctx, topic, "brokers")
	var brokers []string
	for _, broker := range strings.Split(center.GetStringWithNamespace(ctx, defaultApolloNamespace, brokerKey), ",") {
		brokers = append(brokers, strings.TrimSpace(broker))
	}

	// validate config
	if len(brokers) == 0 {
		return nil, fmt.Errorf("%s no brokers config found")
	}

	return &Config{
		MQType:         MqTypeKafka,
		MQAddr:         brokers,
		Topic:          topic,
		TimeOut:        defaultTimeout,
		CommitInterval: 1 * time.Second,
		Offset:         FirstOffset,
	}, nil
}

type apolloObserver struct {
	ch chan<- *center.ChangeEvent
}

func (ob *apolloObserver) HandleChangeEvent(event *center.ChangeEvent) {
	if event.Namespace != defaultApolloNamespace {
		return
	}

	// TODO: filter different mq types
	var changes = map[string]*center.Change{}
	for k, ce := range event.Changes {
		if strings.Contains(k, fmt.Sprint(MqTypeKafka)) {
			changes[k] = ce
		}
	}

	event.Changes = changes

	ob.ch <- event
}

func (m *ApolloConfig) Watch(ctx context.Context) <-chan *center.ChangeEvent {
	m.watchOnce.Do(func() {
		center.StartWatchUpdate(ctx)
		center.RegisterObserver(ctx, &apolloObserver{m.ch})
	})
	return m.ch
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

func (m *ApolloConfig) buildKey(ctx context.Context, topic, item string) string {
	return strings.Join([]string{
		topic,
		scontext.GetGroupWithDefault(ctx, defaultGroup),
		fmt.Sprint(MqTypeKafka),
		item,
	}, apolloConfigSep)
}
