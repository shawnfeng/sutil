// Copyright 2014 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"github.com/shawnfeng/sutil/slog"
	"sync"
	"time"
)

const (
	DB_TYPE_MONGO    = "mongo"
	DB_TYPE_MYSQL    = "mysql"
	DB_TYPE_POSTGRES = "postgres"
)

type ConfigInfo struct {
	DBAddr   string
	UserName string
	PassWord string
	TimeOut  time.Duration
}

type DBConfiger interface {
	GetConfig(dbType, dbName string) *ConfigInfo
}

type DefaultDBConfig struct {
	userName string
	passWord string
}

func NewDefaultDBConfig(servName string) DBConfiger {
	passWord := "MegQqwb@RPZWI55N"
	return &DefaultDBConfig{
		userName: servName,
		passWord: passWord,
	}
}

func (m *DefaultDBConfig) GetConfig(dbType, dbName string) *ConfigInfo {
	return &ConfigInfo{
		DBAddr:   "common.kingshard.pri.ibanyu.com:9090",
		UserName: m.userName,
		PassWord: m.passWord,
		TimeOut:  3 * time.Second,
	}
}

type DBInstance interface {
	getType() string
}

func NewInstance(dbtype, dbname, addr, userName, passWord string, timeout time.Duration) (DBInstance, error) {

	switch dbtype {
	case DB_TYPE_POSTGRES:
		return NewdbSql(dbtype, dbname, addr, userName, passWord, timeout)

	case DB_TYPE_MYSQL:
		return NewdbSql(dbtype, dbname, addr, userName, passWord, timeout)

	default:
		return nil, fmt.Errorf("dbtype %s error", dbtype)
	}
}

type DBInstanceManager struct {
	instances sync.Map
	config    DBConfiger
}

func NewDBInstanceManager(servName string) *DBInstanceManager {
	config := NewDefaultDBConfig(servName)
	return &DBInstanceManager{
		config: config,
	}
}

func (m *DBInstanceManager) buildKey(dbtype, dbname string) string {
	return fmt.Sprintf("%s.%s", dbtype, dbname)
}

func (m *DBInstanceManager) add(dbtype, dbname string, ins DBInstance) {
	m.instances.Store(m.buildKey(dbtype, dbname), ins)
}

func (m *DBInstanceManager) get(dbtype, dbname string) DBInstance {
	fun := "DBInstanceManager.get-->"

	var err error
	var dbIn DBInstance
	name := m.buildKey(dbtype, dbname)
	in, ok := m.instances.Load(name)
	if ok == false {

		configInfo := m.config.GetConfig(dbtype, dbname)
		dbIn, err = NewInstance(dbtype, dbname, configInfo.DBAddr, configInfo.UserName, configInfo.PassWord, configInfo.TimeOut)
		if err != nil {
			slog.Errorf("%s NewInstance err, dbname: %s, err: %s", fun, dbname, err.Error())
			return nil
		}

		m.instances.Store(name, dbIn)
		return dbIn
	}

	dbIn, ok = in.(DBInstance)
	if ok == false {
		slog.Errorf("%s ins.(DBInstance) err, name: %s", fun, name)
		return nil
	}

	return dbIn
}
