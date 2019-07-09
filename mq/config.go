// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"fmt"
	"time"
)

const (
	MqTypeKafka = iota
)

const (
	ConfigTypeSimple = iota
	ConfigTypeEtcd
)

type Config struct {
	MQType         int
	MQAddr         []string
	Topic          string
	TimeOut        time.Duration
	CommitInterval time.Duration
	Offset         int64
}

var DefaultConfiger = NewSimpleConfiger()

type Configer interface {
	GetConfig(topic string) *Config
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

func (m *SimpleConfig) GetConfig(topic string) *Config {
	return &Config{
		MQType:         MqTypeKafka,
		MQAddr:         m.mqAddr,
		Topic:          topic,
		TimeOut:        3 * time.Second,
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

func (m *EtcdConfig) GetConfig(topic string) *Config {
	//todo etcd router
	return &Config{
		MQType:  MqTypeKafka,
		MQAddr:  []string{},
		Topic:   topic,
		TimeOut: 3 * time.Second,
	}
}
