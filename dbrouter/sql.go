// Copyright 2013 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shawnfeng/sutil/slog"
	"time"
)

type DB struct {
	*sqlx.DB
}

func (db *DB) NamedExecWrapper(tables []interface{}, query string, arg interface{}) (sql.Result, error) {
	query = fmt.Sprintf(query, tables...)
	return db.DB.NamedExec(query, arg)
}

func (db *DB) NamedQueryWrapper(tables []interface{}, query string, arg interface{}) (*sqlx.Rows, error) {
	query = fmt.Sprintf(query, tables...)
	return db.DB.NamedQuery(query, arg)
}

func (db *DB) SelectWrapper(tables []interface{}, dest interface{}, query string, args ...interface{}) error {
	query = fmt.Sprintf(query, tables...)
	return db.DB.Select(dest, query, args...)
}

func (db *DB) ExecWrapper(tables []interface{}, query string, args ...interface{}) (sql.Result, error) {
	query = fmt.Sprintf(query, tables...)
	return db.DB.Exec(query, args...)
}

func (db *DB) QueryRowxWrapper(tables []interface{}, query string, args ...interface{}) *sqlx.Row {
	query = fmt.Sprintf(query, tables...)
	return db.DB.QueryRowx(query, args...)
}

func (db *DB) QueryxWrapper(tables []interface{}, query string, args ...interface{}) (*sqlx.Rows, error) {
	query = fmt.Sprintf(query, tables...)
	return db.DB.Queryx(query, args...)
}

func (db *DB) GetWrapper(tables []interface{}, dest interface{}, query string, args ...interface{}) error {
	query = fmt.Sprintf(query, tables...)
	return db.DB.Get(dest, query, args...)
}

func NewDB(sqlxdb *sqlx.DB) *DB {
	db := &DB{
		sqlxdb,
	}
	return db
}

type Sql struct {
	dbType   string
	dbName   string
	dbAddr   string
	timeOut  time.Duration
	userName string
	passWord string
	db       *DB
}

func NewSql(dbtype, dbname, addr, userName, passWord string, timeout time.Duration) (*Sql, error) {
	fun := "NewSql-->"

	if timeout == 0 {
		timeout = 3 * time.Second
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
	info.db, err = dial(info)
	if err != nil {
		slog.Errorf(context.TODO(), "%s info:%v, err:%s", fun, *info, err.Error())
		return nil, err
	}
	info.db.SetMaxIdleConns(8)
	return info, err
}

func dial(info *Sql) (db *DB, err error) {
	fun := "dial-->"

	var dataSourceName string
	if info.dbType == DB_TYPE_MYSQL {
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s)/%s", info.userName, info.passWord, info.dbAddr, info.dbName)

	} else if info.dbType == DB_TYPE_POSTGRES {
		dataSourceName = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			info.userName, info.passWord, info.dbAddr, info.dbName)
	}

	slog.Infof(context.TODO(), "%s dbtype:%s datasourcename:%s", fun, info.dbType, dataSourceName)
	sqlxdb, err := sqlx.Connect(info.dbType, dataSourceName)
	return NewDB(sqlxdb), err
}

func (m *Sql) getDB() *DB {
	return m.db
}

func (m *Sql) GetType() string {
	return m.dbType
}

func (m *Sql) Close() error {
	return m.db.Close()
}
