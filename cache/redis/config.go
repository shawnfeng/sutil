// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"fmt"
	"time"
)

const (
	CONFIG_TYPE_SIMPLE = iota
	CONFIG_TYPE_ETCD
)

type Config struct {
	addr      string
	namespace string
	poolSize  int
	timeout   time.Duration
}

var DefaultConfiger = NewSimpleConfiger()

type Configer interface {
	GetConfig(ctx context.Context, namespace string) *Config
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
}

func NewSimpleConfiger() Configer {
	return &SimpleConfig{}
}

func (m *SimpleConfig) GetConfig(ctx context.Context, namespace string) *Config {
	addr := ""
	if namespace == "base/report" {
		addr = "common.codis.pri.ibanyu.com:19000"
		//addr = "core.codis.pri.ibanyu.com:19000"
	}
	return &Config{
		addr:      addr,
		namespace: namespace,
		timeout:   1 * time.Second,
		poolSize:  128,
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

func (m *EtcdConfig) GetConfig(ctx context.Context, namespace string) *Config {
	//todo etcd router
	return &Config{}
}
