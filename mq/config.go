// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"fmt"
	"time"
)

const (
	MQ_TYPE_KAFKA = iota
)

const (
	CONFIG_TYPE_SIMPLE = iota
	CONFIG_TYPE_ETCD
)

type Config struct {
	MQType  int
	MQAddr  []string
	Topic   string
	TimeOut time.Duration
}

var DefaultConfiger = NewSimpleConfiger()

type Configer interface {
	GetConfig(topic string) *Config
}

func NewConfiger(configType int) (Configer, error) {

	switch configType {
	case CONFIG_TYPE_SIMPLE:
		return NewSimpleConfiger(), nil

	case CONFIG_TYPE_ETCD:
		return NewEtcdConfiger(), nil

	default:
		return nil, fmt.Errorf("configType %s error", configType)
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
		MQType:  MQ_TYPE_KAFKA,
		MQAddr:  m.mqAddr,
		Topic:   topic,
		TimeOut: 3 * time.Second,
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
		MQType:  MQ_TYPE_KAFKA,
		MQAddr:  []string{},
		Topic:   topic,
		TimeOut: 3 * time.Second,
	}
}
