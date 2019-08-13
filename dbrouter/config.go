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
}

func NewConfiger(configType int, data []byte, dbChangeChan chan dbInstanceChange) (Configer, error) {

	switch configType {
	case CONFIG_TYPE_SIMPLE:
		close(dbChangeChan)
		return NewSimpleConfiger(data), nil

	case CONFIG_TYPE_ETCD:
		return NewEtcdConfiger(context.TODO(), dbChangeChan), nil

	default:
		return nil, fmt.Errorf("configType %d error", configType)
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
	parser *Parser
}

func NewEtcdConfiger(ctx context.Context, dbChangeChan chan dbInstanceChange) Configer {
	fun := "NewEtcdConfiger -->"
	etcdConfig := &EtcdConfig{
		etcdAddr: []string{"http://infra0.etcd.ibanyu.com:20002", "http://infra1.etcd.ibanyu.com:20002", "http://infra2.etcd.ibanyu.com:20002", "http://infra3.etcd.ibanyu.com:20002", "http://infra4.etcd.ibanyu.com:20002", "http://old0.etcd.ibanyu.com:20002", "http://old1.etcd.ibanyu.com:20002", "http://old2.etcd.ibanyu.com:20002"},
	}
	err := etcdConfig.init(ctx, dbChangeChan)
	if err != nil {
		slog.Errorf(ctx, "%s init etcd configer err: %s", fun, err.Error())
		return nil
	}
	return etcdConfig
}

func (m *EtcdConfig) init(ctx context.Context, dbChangeChan chan dbInstanceChange) error {
	fun := "EtcdConfig.init -->"
	etcdInstance, err := setcd.NewEtcdInstance(m.etcdAddr)
	if err != nil {
		return err
	}

	initCh := make(chan error)
	var initOnce sync.Once
	etcdInstance.Watch(ctx, "/roc/db/route", func(response *client.Response) {
		slog.Infof(ctx, "get db conf: %s", response.Node.Value)
		parser, er := NewParserEtcd([]byte(response.Node.Value))
		initOnce.Do(func() {
			initCh <- er
		})

		if er != nil {
			slog.Errorf(ctx, "%s init db parser err: ", fun, er.Error())
		} else {
			slog.Infof(ctx, "succeed to init new parser")
			if m.parser != nil {
				dbInsChange := compareParsers(*m.parser, *parser)
				slog.Infof(ctx, "parser changes: %+v", dbInsChange)
				m.parser = parser
				dbChangeChan <- dbInsChange
			} else {
				m.parser = parser
			}
		}
	})
	// 做一次同步，等parser初始化完成
	err = <- initCh
	close(initCh)
	return err
}

func (m *EtcdConfig) GetConfig(ctx context.Context, instance string) *Config {
	fun := "EtcdConfig.GetConfig --> "
	var info *dbInsInfo
	// TODO 确定是否压测的标识
	group := scontext.GetGroup(ctx)
	switch group {
	case "":
		info = m.parser.getConfig(instance)
	case "default":
		info = m.parser.getConfig(instance)
	case "xxx":
		info = m.parser.GetShadowConfig(instance)
	default:
		// TODO 这种情况容不容易出现？
		slog.Errorf(ctx, "%s invalid context group: %s", fun, group)
		info = &dbInsInfo{}
	}
	/*if isTest, ok := ctx.Value("xxx").(bool); ok {
		if isTest {
			info = m.parser.GetShadowConfig(instance)
		} else {
			info = m.parser.getConfig(instance)
		}
	} else {
		info = m.parser.getConfig(instance)
	}*/
	//todo etcd router
	return &Config{
		DBType:   info.DBType,
		DBAddr:   info.DBAddr,
		UserName: info.UserName,
		PassWord: info.PassWord,
		TimeOut:  3 * time.Second,
	}
}

func (m *EtcdConfig) GetInstance(ctx context.Context, cluster, table string) (instance string) {
	return m.parser.GetInstance(cluster, table)
}
