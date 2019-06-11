// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/slog"
	"github.com/shawnfeng/sutil/stime"
	"gopkg.in/mgo.v2"
)

type Router struct {
	configer  Configer
	instances *InstanceManager
	stat      *statReport
}

func NewRouter(data []byte) (*Router, error) {

	configer := NewSimpleConfiger(data)
	factory := func(configer Configer) func(ctx context.Context, key string) (in Instancer, err error) {
		return func(ctx context.Context, key string) (in Instancer, err error) {
			return Factory(ctx, key, configer)
		}
	}(configer)
	return &Router{
		configer:  NewSimpleConfiger(data),
		instances: NewInstanceManager(factory),
		stat:      newStat(),
	}, nil
}

func (m *Router) StatInfo() []*QueryStat {
	return m.stat.statInfo()
}

func (m *Router) SqlExec(ctx context.Context, cluster string, query func(*DB, []interface{}) error, tables ...string) error {
	fun := "Router.SqlExec -->"

	st := stime.NewTimeStat()

	if len(tables) <= 0 {
		return fmt.Errorf("tables is empty")
	}

	table := tables[0]
	dbType, dbName := m.configer.GetTypeAndName(ctx, cluster, table)
	in := m.instances.Get(ctx, generateKey(dbType, dbName))
	if in == nil {
		return fmt.Errorf("db instance not find: dbname:%s", dbName)
	}

	dbsql, ok := in.(*Sql)
	if !ok {
		return fmt.Errorf("db instance type error: dbname:%s, dbtype:%s", dbName, in.GetType())
	}

	db := dbsql.getDB()

	defer func() {
		dur := st.Duration()
		m.stat.incQuery(cluster, table, st.Duration())
		slog.Infof("%s type:%s dbname:%s query:%d", fun, in.GetType(), dbName, dur)
	}()

	var tmptables []interface{}
	for _, item := range tables {
		tmptables = append(tmptables, item)
	}

	return query(db, tmptables)
}

func (m *Router) MongoExecEventual(ctx context.Context, cluster, table string, query func(*mgo.Collection) error) error {
	return m.mongoExec(ctx, eventual, cluster, table, query)
}

func (m *Router) MongoExecMonotonic(ctx context.Context, cluster, table string, query func(*mgo.Collection) error) error {
	return m.mongoExec(ctx, monotonic, cluster, table, query)
}

func (m *Router) MongoExecStrong(ctx context.Context, cluster, table string, query func(*mgo.Collection) error) error {
	return m.mongoExec(ctx, strong, cluster, table, query)
}

func (m *Router) mongoExec(ctx context.Context, consistency mode, cluster, table string, query func(*mgo.Collection) error) error {
	st := stime.NewTimeStat()

	dbType, dbName := m.configer.GetTypeAndName(ctx, cluster, table)
	in := m.instances.Get(ctx, generateKey(dbType, dbName))
	if in == nil {
		return fmt.Errorf("db instance not find: cluster:%s table:%s", cluster, table)
	}

	db, ok := in.(*dbMongo)
	if !ok {
		return fmt.Errorf("db instance type error: cluster:%s table:%s type:%s", cluster, table, in.GetType())
	}

	sess, err := db.getSession(consistency)
	if err != nil {
		return err
	}

	if sess == nil {
		return fmt.Errorf("db instance session empty: cluster:%s table:%s type:%s", cluster, table, in.GetType())
	}

	sessionCopy := sess.Copy()
	defer sessionCopy.Close()
	c := sessionCopy.DB("").C(table)

	defer func() {
		dur := st.Duration()
		m.stat.incQuery(cluster, table, st.Duration())
		slog.Tracef("[MONGO] const:%d cls:%s table:%s dur:%d", consistency, cluster, table, dur)
	}()

	return query(c)
}
