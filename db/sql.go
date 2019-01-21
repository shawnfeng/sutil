// Copyright 2013 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shawnfeng/sutil/slog"
	"github.com/shawnfeng/sutil/stime"
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

type dbSql struct {
	dbType   string
	dbName   string
	dbAddr   string
	timeOut  time.Duration
	userName string
	passWord string
	db       *DB
}

func (m *dbSql) getType() string {
	return m.dbType
}

func NewdbSql(dbtype, dbname, addr, userName, passWord string, timeout time.Duration) (*dbSql, error) {

	if timeout == 0 {
		timeout = 3 * time.Second
	}

	info := &dbSql{
		dbType:   dbtype,
		dbName:   dbname,
		dbAddr:   addr,
		timeOut:  timeout,
		userName: userName,
		passWord: passWord,
	}

	var err error
	info.db, err = dial(info)
	info.db.SetMaxIdleConns(32)
	return info, err
}

func dial(info *dbSql) (db *DB, err error) {
	fun := "dial-->"

	var dataSourceName string
	if info.dbType == DB_TYPE_MYSQL {
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s)/%s", info.userName, info.passWord, info.dbAddr, info.dbName)

	} else if info.dbType == DB_TYPE_POSTGRES {
		dataSourceName = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			info.userName, info.passWord, info.dbAddr, info.dbName)
	}

	slog.Infof("%s dbtype:%s datasourcename:%s", fun, info.dbType, dataSourceName)
	sqlxdb, err := sqlx.Connect(info.dbType, dataSourceName)
	return NewDB(sqlxdb), err
}

func (m *dbSql) getDB() *DB {
	return m.db
}

func (m *DBInstanceManager) SqlExec(dbName string, query func(*DB, []interface{}) error, tables ...string) error {
	fun := "SqlExec-->"

	st := stime.NewTimeStat()

	if len(tables) <= 0 {
		return fmt.Errorf("tables is empty")
	}

	ins := m.get(DB_TYPE_MYSQL, dbName)
	if ins == nil {
		return fmt.Errorf("db instance not find: dbname:%s", dbName)
	}

	dbsql, ok := ins.(*dbSql)
	if !ok {
		return fmt.Errorf("db instance type error: dbname:%s, dbtype:%s", dbName, ins.getType())
	}

	db := dbsql.getDB()

	defer func() {
		dur := st.Duration()
		slog.Infof("%s type:%s dbname:%s query:%d", fun, ins.getType(), dbName, dur)
	}()

	var tmptables []interface{}
	for _, item := range tables {
		tmptables = append(tmptables, item)
	}

	return query(db, tmptables)
}
