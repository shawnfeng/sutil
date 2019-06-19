// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/slog/slog"
	stat "github.com/shawnfeng/sutil/stat"
	"github.com/shawnfeng/sutil/stime"
	"gopkg.in/mgo.v2"
)

const (
	SpanTagKeyCluster = "cluster"
	SpanTagKeyTable = "table"
)

type Router struct {
	configer  Configer
	instances *InstanceManager
	report    *stat.StatReport
}

func NewRouter(data []byte) (*Router, error) {
	configer := NewSimpleConfiger(data)
	factory := func(configer Configer) func(ctx context.Context, key string) (in Instancer, err error) {
		return func(ctx context.Context, key string) (in Instancer, err error) {
			return Factory(ctx, key, configer)
		}
	}(configer)
	return &Router{
		configer:  configer,
		instances: NewInstanceManager(factory),
		report:    stat.NewStat(),
	}, nil
}

func (m *Router) StatInfo() []*stat.QueryStat {
	return m.report.StatInfo()
}

func (m *Router) SqlExec(ctx context.Context, cluster string, query func(*DB, []interface{}) error, tables ...string) error {
	fun := "Router.SqlExec -->"

	span, ctx := opentracing.StartSpanFromContext(ctx, "dbrouter.SqlExec")
	defer span.Finish()
	span.SetTag(SpanTagKeyCluster, cluster)

	st := stime.NewTimeStat()

	if len(tables) <= 0 {
		return fmt.Errorf("tables is empty")
	}

	table := tables[0]
	span.SetTag(SpanTagKeyTable, table)
	instance := m.configer.GetInstance(ctx, cluster, table)
	in := m.instances.Get(ctx, generateKey(instance))
	if in == nil {
		return fmt.Errorf("db instance not find: instance:%s", instance)
	}

	dbsql, ok := in.(*Sql)
	if !ok {
		return fmt.Errorf("db instance type error: instance:%s, dbtype:%s", instance, in.GetType())
	}

	db := dbsql.getDB()

	defer func() {
		dur := st.Duration()
		m.report.IncQuery(cluster, table, st.Duration())
		slog.Tracef(ctx, "%s cls:%s table:%s dur:%d", fun, cluster, table, dur)
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
	fun := "Router.mongoExec -->"

	span, ctx := opentracing.StartSpanFromContext(ctx, "dbrouter.mongoExec")
	defer span.Finish()
	span.SetTag("cluster", cluster)
	span.SetTag("table", table)

	st := stime.NewTimeStat()

	instance := m.configer.GetInstance(ctx, cluster, table)
	in := m.instances.Get(ctx, generateKey(instance))
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
		m.report.IncQuery(cluster, table, st.Duration())
		slog.Tracef(ctx, "%s const:%d cls:%s table:%s dur:%d", fun, consistency, cluster, table, dur)
	}()

	return query(c)
}
