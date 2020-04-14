// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/shawnfeng/sutil/slog/slog"
	"github.com/shawnfeng/sutil/stat"
	"github.com/shawnfeng/sutil/stime"
	"gopkg.in/mgo.v2"
)

const (
	spanLogKeyCluster = "cluster"
	spanLogKeyTable   = "table"
)

type Router struct {
	configer  Configer
	instances *InstanceManager
	report    *stat.StatReport
}

type dbConfigChange struct {
	dbInstanceChange map[string][]string
	dbGroups         []string
}

func NewRouter(data []byte) (*Router, error) {
	// TODO config type由哪里决定
	var dbChangeChan = make(chan dbConfigChange)
	configer, err := NewConfiger(CONFIG_TYPE_ETCD, data, dbChangeChan)
	if err != nil {
		return nil, err
	}

	factory := func(configer Configer) func(ctx context.Context, key, group string) (in Instancer, err error) {
		return func(ctx context.Context, key, group string) (in Instancer, err error) {
			return Factory(ctx, key, group, configer)
		}
	}(configer)

	return &Router{
		configer:  configer,
		instances: NewInstanceManager(factory, dbChangeChan, configer.GetGroups(context.TODO())),
		report:    stat.NewStat(),
	}, nil
}

func (m *Router) StatInfo() []*stat.QueryStat {
	return m.report.StatInfo()
}

func (m *Router) sqlPrepare(ctx context.Context, cluster, table string) (db *DB, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "dbrouter.sqlPrepare")
	defer span.Finish()

	instance := m.configer.GetInstance(ctx, cluster, table)
	in := m.instances.Get(ctx, generateKey(instance))
	if in == nil {
		err = fmt.Errorf("db instance not find: cluster:%s table:%s instance:%s", cluster, table, instance)
		return
	}

	dbsql, ok := in.(*Sql)
	if !ok {
		err = fmt.Errorf("db instance type error: cluster:%s table:%s instance:%s, dbtype:%s", cluster, table, instance, in.GetType())
		return
	}

	db = dbsql.getDB()
	return
}

func (m *Router) SqlExec(ctx context.Context, cluster string, query func(*DB, []interface{}) error, tables ...string) error {
	fun := "Router.SqlExec -->"

	span, ctx := opentracing.StartSpanFromContext(ctx, "dbrouter.SqlExec")
	defer span.Finish()

	st := stime.NewTimeStat()

	if len(tables) <= 0 {
		return fmt.Errorf("tables is empty")
	}
	table := tables[0]

	span.LogFields(
		log.String(spanLogKeyCluster, cluster),
		log.String(spanLogKeyTable, table))

	// check breaker
	if !Entry(cluster, table){
		slog.Errorf(ctx, "%s trigger tidb breaker, because too many timeout sqls, cluster: %s, table: %s", fun, cluster, table)
		return errors.New("sql cause breaker, because too many timeout")
	}

	db, err := m.sqlPrepare(ctx, cluster, table)
	if err != nil {
		return err
	}

	defer func() {
		dur := st.Duration()
		m.report.IncQuery(cluster, table, st.Duration())
		slog.Tracef(ctx, "%s cls:%s table:%s dur:%d", fun, cluster, table, dur)
	}()

	var tmptables []interface{}
	for _, item := range tables {
		tmptables = append(tmptables, item)
	}
	err = query(db, tmptables)
	statReqErr(cluster, table, err)
	// record breaker
	statBreaker(cluster, table, err)
	return err
}

func (m *Router) ormPrepare(ctx context.Context, cluster, table string) (db *GormDB, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "dbrouter.ormPrepare")
	defer span.Finish()

	instance := m.configer.GetInstance(ctx, cluster, table)
	in := m.instances.Get(ctx, generateKey(instance))
	if in == nil {
		err = fmt.Errorf("db instance not find: instance:%s", instance)
		return
	}

	dbsql, ok := in.(*Sql)
	if !ok {
		err = fmt.Errorf("db instance type error: instance:%s, dbtype:%s", instance, in.GetType())
		return
	}

	db = dbsql.getGormDB()
	return
}

func (m *Router) OrmExec(ctx context.Context, cluster string, query func(*GormDB, []interface{}) error, tables ...string) error {
	fun := "Router.OrmExec -->"

	span, ctx := opentracing.StartSpanFromContext(ctx, "dbrouter.OrmExec")
	defer span.Finish()

	st := stime.NewTimeStat()

	if len(tables) <= 0 {
		return fmt.Errorf("tables is empty")
	}

	table := tables[0]
	span.LogFields(
		log.String(spanLogKeyCluster, cluster),
		log.String(spanLogKeyTable, table))

	// check breaker
	if !Entry(cluster, table){
		slog.Errorf(ctx, "%s trigger tidb breaker, because too many timeout sqls, cluster: %s, table: %s", fun, cluster, table)
		return errors.New("sql cause breaker, because too many timeout")
	}

	db, err := m.ormPrepare(ctx, cluster, table)
	if err != nil {
		return err
	}

	defer func() {
		dur := st.Duration()
		m.report.IncQuery(cluster, table, st.Duration())
		slog.Tracef(ctx, "%s cls:%s table:%s dur:%d", fun, cluster, table, dur)
	}()

	var tmptables []interface{}
	for _, item := range tables {
		tmptables = append(tmptables, item)
	}
	err = query(db, tmptables)
	statReqErr(cluster, table, err)
	// stat breaker
	statBreaker(cluster, table, err)
	return err
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

func (m *Router) mongoPrepare(ctx context.Context, consistency mode, cluster, table string) (sess *mgo.Session, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "dbrouter.mongoPrepare")
	defer span.Finish()

	instance := m.configer.GetInstance(ctx, cluster, table)
	in := m.instances.Get(ctx, generateKey(instance))
	if in == nil {
		err = fmt.Errorf("db instance not find: cluster:%s table:%s", cluster, table)
		return
	}

	db, ok := in.(*dbMongo)
	if !ok {
		err = fmt.Errorf("db instance type error: cluster:%s table:%s type:%s", cluster, table, in.GetType())
		return
	}

	ss, err := db.getSession(consistency)
	if err != nil {
		return
	}

	if ss == nil {
		err = fmt.Errorf("db instance session empty: cluster:%s table:%s type:%s", cluster, table, in.GetType())
		return
	}

	sess = ss.Copy()
	return
}

func (m *Router) mongoExec(ctx context.Context, consistency mode, cluster, table string, query func(*mgo.Collection) error) error {
	fun := "Router.mongoExec -->"
	if !Entry(cluster, table){
		slog.Errorf(ctx, "%s trigger mongodb breaker, because too many timeout query, cluster: %s, table: %s", fun, cluster, table)
		return errors.New("mongo query cause breaker, because too many timeout")
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "dbrouter.mongoExec")
	defer span.Finish()
	span.LogFields(
		log.String(spanLogKeyCluster, cluster),
		log.String(spanLogKeyTable, table))

	st := stime.NewTimeStat()

	sess, err := m.mongoPrepare(ctx, consistency, cluster, table)
	if err != nil {
		return err
	}
	defer sess.Close()
	coll := sess.DB("").C(table)

	defer func() {
		dur := st.Duration()
		m.report.IncQuery(cluster, table, st.Duration())
		slog.Tracef(ctx, "%s const:%d cls:%s table:%s dur:%d", fun, consistency, cluster, table, dur)
	}()
	err = query(coll)
	statReqErr(cluster, table, err)
	statBreaker(cluster, table, err)
	return err
}
