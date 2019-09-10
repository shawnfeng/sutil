// Copyright 2013 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	//	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/shawnfeng/sutil/slog/slog"
	"time"
)

type Sql struct {
	dbType   string
	dbName   string
	dbAddr   string
	timeOut  time.Duration
	userName string
	passWord string
	db       *DB
	gormdb   *GormDB
}

func NewSql(dbtype, dbname, addr, userName, passWord string, timeout time.Duration) (*Sql, error) {
	fun := "NewSql-->"

	if timeout == 0 {
		timeout = 60 * time.Second
	}

	info := &Sql{
		dbType:   dbtype,
		dbName:   dbname,
		dbAddr:   addr,
		timeOut:  timeout,
		userName: userName,
		passWord: passWord,
	}

	var err error
	info.db, info.gormdb, err = dial(info)
	if err != nil {
		slog.Errorf(context.TODO(), "%s info:%v, err:%s", fun, *info, err.Error())
		return nil, err
	}
	return info, err
}

func dial(info *Sql) (*DB, *GormDB, error) {
	fun := "dial -->"

	sqlxdb, err := dialBySqlx(info)
	if err != nil {
		slog.Errorf(context.TODO(), "%s info:%v, dialBySqlx err:%s", fun, *info, err.Error())
		return nil, nil, err
	}

	gormdb, err := dialByGorm(info)
	if err != nil {
		slog.Errorf(context.TODO(), "%s info:%v, dialByGorm err:%s", fun, *info, err.Error())
		return nil, nil, err
	}

	return NewDB(sqlxdb), NewGormDB(gormdb), nil
}

func (m *Sql) getDB() *DB {
	return m.db
}

func (m *Sql) getGormDB() *GormDB {
	return m.gormdb
}

func (m *Sql) GetType() string {
	return m.dbType
}

func (m *Sql) Close() error {
	err1 := m.db.Close()
	err2 := m.gormdb.Close()

	if err1 != nil || err2 != nil {
		return fmt.Errorf("sqlx.Close err: %v, gorm.Close err: %v", err1, err2)
	}

	return nil
}
