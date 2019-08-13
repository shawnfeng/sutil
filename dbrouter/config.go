// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/setcd"
	"github.com/shawnfeng/sutil/slog/slog"
	"sync"
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
	GetConfigByGroup(ctx context.Context, instance, group string) *Config
	GetGroups(ctx context.Context) []string
}

func NewConfiger(configType int, data []byte, dbChangeChan chan dbConfigChange) (Configer, error) {

	switch configType {
	case CONFIG_TYPE_SIMPLE:
		close(dbChangeChan)
		return NewSimpleConfiger(data)

	case CONFIG_TYPE_ETCD:
		return NewEtcdConfiger(context.TODO(), dbChangeChan)

	default:
		return nil, fmt.Errorf("configType %d error", configType)
	}
}

type SimpleConfig struct {
	parser *Parser
}

func NewSimpleConfiger(data []byte) (Configer, error) {
	parser, err := NewParser(data)
	return &SimpleConfig{
		parser: parser,
	}, err
}

func (m *SimpleConfig) GetConfig(ctx context.Context, instance string) *Config {
	group := scontext.GetGroup(ctx)
	info := m.parser.GetConfig(instance, group)
	return &Config{
		DBType:   info.DBType,
		DBAddr:   info.DBAddr,
		DBName:   info.DBName,
		UserName: info.UserName,
		PassWord: info.PassWord,
		TimeOut:  3 * time.Second,
	}
}

func (m *SimpleConfig) GetConfigByGroup(ctx context.Context, instance, group string) *Config {
	info := m.parser.GetConfig(instance, group)
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

func (m *SimpleConfig) GetGroups(ctx context.Context) []string {
	var groups []string
	for group, _ := range m.parser.dbIns {
		groups = append(groups, group)
	}
	return groups
}

type EtcdConfig struct {
	etcdAddr []string
	parser *Parser
}

func NewEtcdConfiger(ctx context.Context, dbChangeChan chan dbConfigChange) (Configer, error) {
	fun := "NewEtcdConfiger -->"
	etcdConfig := &EtcdConfig{
		// TODO etcd address如何获取
		etcdAddr: []string{"http://infra0.etcd.ibanyu.com:20002", "http://infra1.etcd.ibanyu.com:20002", "http://infra2.etcd.ibanyu.com:20002", "http://infra3.etcd.ibanyu.com:20002", "http://infra4.etcd.ibanyu.com:20002", "http://old0.etcd.ibanyu.com:20002", "http://old1.etcd.ibanyu.com:20002", "http://old2.etcd.ibanyu.com:20002"},
	}
	err := etcdConfig.init(ctx, dbChangeChan)
	if err != nil {
		slog.Errorf(ctx, "%s init etcd configer err: %s", fun, err.Error())
		return nil, err
	}
	return etcdConfig, nil
}

func (m *EtcdConfig) init(ctx context.Context, dbChangeChan chan dbConfigChange) error {
	fun := "EtcdConfig.init -->"
	etcdInstance, err := setcd.NewEtcdInstance(m.etcdAddr)
	if err != nil {
		return err
	}

	initCh := make(chan error)
	var initOnce sync.Once
	etcdInstance.Watch(ctx, "/roc/db/route", func(response *client.Response) {
		slog.Infof(ctx, "get db conf: %s", response.Node.Value)
		parser, er := NewParser([]byte(response.Node.Value))

		if er != nil {
			slog.Errorf(ctx, "%s init db parser err: ", fun, er.Error())
		} else {
			slog.Infof(ctx, "succeed to init new parser")
			if m.parser != nil {
				dbConfigChange := compareParsers(*m.parser, *parser)
				slog.Infof(ctx, "parser changes: %+v", dbConfigChange)
				m.parser = parser
				dbChangeChan <- dbConfigChange
			} else {
				m.parser = parser
			}
		}

		initOnce.Do(func() {
			initCh <- er
		})
	})
	// 做一次同步，等parser初始化完成
	err = <- initCh
	close(initCh)
	return err
}

func (m *EtcdConfig) GetConfig(ctx context.Context, instance string) *Config {
	group := scontext.GetGroup(ctx)
	info := m.parser.GetConfig(instance, group)
	return &Config{
		DBType:   info.DBType,
		DBAddr:   info.DBAddr,
		DBName:   info.DBName,
		UserName: info.UserName,
		PassWord: info.PassWord,
		TimeOut:  3 * time.Second,
	}
}

func (m *EtcdConfig) GetConfigByGroup(ctx context.Context, instance, group string) *Config {
	info := m.parser.GetConfig(instance, group)
	return &Config{
		DBType:   info.DBType,
		DBAddr:   info.DBAddr,
		DBName:   info.DBName,
		UserName: info.UserName,
		PassWord: info.PassWord,
		TimeOut:  3 * time.Second,
	}
}

func (m *EtcdConfig) GetInstance(ctx context.Context, cluster, table string) (instance string) {
	return m.parser.GetInstance(cluster, table)
}

func (m *EtcdConfig) GetGroups(ctx context.Context) []string {
	var groups []string
	for group, _ := range m.parser.dbIns {
		groups = append(groups, group)
	}
	return groups
}
