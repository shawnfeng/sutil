// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"fmt"
	"time"
)

const (
	DB_TYPE_MONGO    = "mongo"
	DB_TYPE_MYSQL    = "mysql"
	DB_TYPE_POSTGRES = "postgres"
)

const (
	CONFIG_TYPE_SIMPLE = iota
	CONFIG_TYPE_ETCD
)

type Config struct {
	DBName   string
	DBType   string
	DBAddr   []string
	UserName string
	PassWord string
	TimeOut  time.Duration
}

type Configer interface {
	GetConfig(ctx context.Context, instance string) *Config
	GetInstance(ctx context.Context, cluster, table string) (instance string)
}

func NewConfiger(configType int, data []byte) (Configer, error) {

	switch configType {
	case CONFIG_TYPE_SIMPLE:
		return NewSimpleConfiger(data), nil

	case CONFIG_TYPE_ETCD:
		return NewEtcdConfiger(), nil

	default:
		return nil, fmt.Errorf("configType %s error", configType)
	}
}

type SimpleConfig struct {
	parser *Parser
}

func NewSimpleConfiger(data []byte) Configer {
	parser, _ := NewParser(data)
	return &SimpleConfig{
		parser: parser,
	}
}

func (m *SimpleConfig) GetConfig(ctx context.Context, instance string) *Config {
	info := m.parser.GetConfig(instance)
	return &Config{
		DBType:   info.DBType,
		DBAddr:   info.DBAddr,
		DBName:   info.DBName,
		UserName: info.UserName,
		PassWord: info.PassWord,
		TimeOut:  3 * time.Second,
	}
}

func (m *SimpleConfig) GetInstance(ctx context.Context, cluster, table string) (instance string) {
	instance = m.parser.GetInstance(cluster, table)
	return instance
}

type EtcdConfig struct {
	etcdAddr []string
}

func NewEtcdConfiger() Configer {
	return &EtcdConfig{
		etcdAddr: []string{}, //todo
	}
}

func (m *EtcdConfig) GetConfig(ctx context.Context, instance string) *Config {
	//todo etcd router
	return &Config{
		DBType:   "",
		DBAddr:   []string{},
		UserName: "",
		PassWord: "",
		TimeOut:  3 * time.Second,
	}
}

func (m *EtcdConfig) GetInstance(ctx context.Context, cluster, table string) (instance string) {
	return ""
}
