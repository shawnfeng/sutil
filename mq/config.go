// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"strings"
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

var DefaultConfiger = NewSimpleConfiger()

type Configer interface {
	GetConfig(ctx context.Context, topic string) *Config
}

func NewConfiger(configType int) (Configer, error) {

	switch configType {
	case ConfigTypeSimple:
		return NewSimpleConfiger(), nil

	case ConfigTypeEtcd:
		return NewEtcdConfiger(), nil

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

func (m *SimpleConfig) GetConfig(ctx context.Context, topic string) *Config {
	fun := "SimpleConfig.GetConfig-->"
	slog.Infof(ctx, "%s get simple config topic:%s", fun, topic)

	return &Config{
		MQType:         MqTypeKafka,
		MQAddr:         m.mqAddr,
		Topic:          topic,
		TimeOut:        defaultTimeout,
		CommitInterval: 1 * time.Second,
		Offset:         FirstOffset,
	}
}

type EtcdConfig struct {
	etcdAddr []string
}

func NewEtcdConfiger() Configer {
	return &EtcdConfig{
		etcdAddr: []string{}, //todo
	}
}

func (m *EtcdConfig) GetConfig(ctx context.Context, topic string) *Config {
	fun := "EtcdConfig.GetConfig-->"
	slog.Infof(ctx, "%s get etcd config topic:%s", fun, topic)

	//todo etcd router
	return &Config{
		MQType:  MqTypeKafka,
		MQAddr:  []string{},
		Topic:   topic,
		TimeOut: defaultTimeout,
	}
}

const defaultApolloNamespace = "infra.mq"

func getApolloMQConfigKey(topic, group, mqType, key string) string {
	return strings.Join([]string{
		topic,
		group,
		mqType,
		key,
	}, ".")
}

type ApolloConfig struct {}

func (m *ApolloConfig) GetConfig(ctx context.Context, topic string) *Config {
	fun := "ApolloConfig.GetConfig-->"
	slog.Infof(ctx, "%s get mq config topic:%s", fun, topic)

	group := scontext.GetGroup(ctx)
	if group == "" {
		group = "default"
	}

	brokerKey := getApolloMQConfigKey(topic, group, fmt.Sprint(MqTypeKafka), "brokers")
	var brokers []string
	for _, broker := range strings.Split(center.GetStringWithNamespace(context.TODO(), defaultApolloNamespace, brokerKey), ",") {
		brokers = append(brokers, strings.TrimSpace(broker))
	}

	return &Config{
		MQType:         MqTypeKafka,
		MQAddr:         brokers,
		Topic:          topic,
		TimeOut:        defaultTimeout,
		CommitInterval: 1 * time.Second,
		Offset:         FirstOffset,
	}
}

func NewApolloConfig() Configer {
	return &ApolloConfig{}
}
