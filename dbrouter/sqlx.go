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
	"github.com/shawnfeng/sutil/slog/slog"
	"time"
)

type DB struct {
	*sqlx.DB
}

func NewDB(sqlxdb *sqlx.DB) *DB {
	db := &DB{
		sqlxdb,
	}
	return db
}

func dialBySqlx(info *Sql) (db *sqlx.DB, err error) {
	fun := "dialBySqlx -->"

	var dataSourceName string
	if info.dbType == DB_TYPE_MYSQL { // timeout: 3s readTimeout: 5s writeTimeout: 5s, TODO: dynamic config
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&timeout=3s&readTimeout=5s&writeTimeout=5s", info.userName, info.passWord, info.dbAddr, info.dbName)

	} else if info.dbType == DB_TYPE_POSTGRES {
		dataSourceName = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			info.userName, info.passWord, info.dbAddr, info.dbName)
	}

	slog.Infof(context.TODO(), "%s dbtype:%s datasourcename:%s", fun, info.dbType, dataSourceName)
	sqlxdb, err := sqlx.Connect(info.dbType, dataSourceName)
	if err == nil {
		sqlxdb.SetMaxIdleConns(16)
		sqlxdb.SetMaxOpenConns(128)
		sqlxdb.SetConnMaxLifetime(time.Hour * 6)
	}
	return sqlxdb, err
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
