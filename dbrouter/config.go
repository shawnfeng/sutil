// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"fmt"
	//	"github.com/shawnfeng/sutil/slog"
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
	GetConfig(ctx context.Context, dbType, dbName string) *Config
	GetTypeAndName(ctx context.Context, cluster, table string) (dbType, dbName string)
}

/*
var DefaultConfiger Configer

func InitDefaultConfiger(configType int, data []byte) {
	//fun := "InitDefaultConfiger -->"

	configer, err := NewConfiger(configType, data)
	if err != nil {
		slog.Errorf("%s NewConfiger, configType: %s", configType)
	}

	DefaultConfiger = configer
}
*/

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

func (m *SimpleConfig) GetConfig(ctx context.Context, dbType, dbName string) *Config {
	info := m.parser.GetConfig(dbType, dbName)
	return &Config{
		DBType:   info.dbType,
		DBAddr:   info.dbAddr,
		DBName:   info.dbName,
		UserName: info.userName,
		PassWord: info.passWord,
		TimeOut:  3 * time.Second,
	}
}

func (m *SimpleConfig) GetTypeAndName(ctx context.Context, cluster, table string) (dbType, dbName string) {
	return m.parser.GetTypeAndName(cluster, table)
}

type EtcdConfig struct {
	etcdAddr []string
}

func NewEtcdConfiger() Configer {
	return &EtcdConfig{
		etcdAddr: []string{}, //todo
	}
}

func (m *EtcdConfig) GetConfig(ctx context.Context, dbType, dbName string) *Config {
	//todo etcd router
	return &Config{
		DBType:   "",
		DBAddr:   []string{},
		UserName: "",
		PassWord: "",
		TimeOut:  3 * time.Second,
	}
}

func (m *EtcdConfig) GetTypeAndName(ctx context.Context, cluster, table string) (dbType, dbName string) {
	return "", ""
}
