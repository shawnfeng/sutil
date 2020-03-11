// Copyright 2013 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shawnfeng/sutil/slog/slog"
	"time"
)

type GormDB struct {
	*gorm.DB
}

func NewGormDB(gormdb *gorm.DB) *GormDB {
	db := &GormDB{
		gormdb,
	}
	return db
}

func dialByGorm(info *Sql) (db *gorm.DB, err error) {
	fun := "dialByGorm -->"

	var dataSourceName string
	if info.dbType == DB_TYPE_MYSQL { // timeout: 3s readTimeout: 5s writeTimeout: 5s, TODO: dynamic config
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci&timeout=3s&readTimeout=5s&writeTimeout=5s", info.userName, info.passWord, info.dbAddr, info.dbName)

	} else if info.dbType == DB_TYPE_POSTGRES {
		dataSourceName = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			info.userName, info.passWord, info.dbAddr, info.dbName)
	}

	slog.Infof(context.TODO(), "%s dbtype:%s datasourcename:%s", fun, info.dbType, dataSourceName)
	gormdb, err := gorm.Open(info.dbType, dataSourceName)
	if err == nil {
		gormdb.DB().SetMaxIdleConns(8)
		gormdb.DB().SetMaxOpenConns(128)
		gormdb.DB().SetConnMaxLifetime(time.Hour*6)
	}

	return gormdb, err
}
