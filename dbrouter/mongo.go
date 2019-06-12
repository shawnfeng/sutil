// Copyright 2014 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"gopkg.in/mgo.v2"
	"sync"
	"time"
)

type dbMongo struct {
	dbType   string
	dbName   string
	dialInfo *mgo.DialInfo

	sessMu  sync.RWMutex
	session [3]*mgo.Session
}

func (m *dbMongo) GetType() string {
	return m.dbType
}

func NewMongo(dbtype, dbname, user, passwd string, addrs []string, timeout time.Duration) (*dbMongo, error) {

	if timeout == 0 {
		timeout = 5 * time.Second
	}

	info := &mgo.DialInfo{
		Addrs:     addrs,
		Timeout:   timeout,
		Database:  dbname,
		Username:  user,
		Password:  passwd,
		PoolLimit: 128,
	}

	return &dbMongo{
		dbType:   dbtype,
		dbName:   dbname,
		dialInfo: info,
	}, nil

}

type mode int

const (
	eventual  mode = 0
	monotonic mode = 1
	strong    mode = 2
)

func dialConsistency(info *mgo.DialInfo, consistency mode) (session *mgo.Session, err error) {

	session, err = mgo.DialWithInfo(info)
	if err != nil {
		return
	}
	session.SetSyncTimeout(1 * time.Minute)
	session.SetSocketTimeout(1 * time.Minute)

	switch consistency {
	case eventual:
		session.SetMode(mgo.Eventual, true)
	case monotonic:
		session.SetMode(mgo.Monotonic, true)
	case strong:
		session.SetMode(mgo.Strong, true)
	}

	return
}

func dialConsistencyWithUrl(url string, timeout time.Duration, consistency mode) (session *mgo.Session, err error) {

	session, err = mgo.DialWithTimeout(url, timeout)
	if err != nil {
		return
	}
	// 看Dial内部的实现
	session.SetSyncTimeout(1 * time.Minute)
	// 不设置这个在执行写入，表不存在时候会报 read tcp 127.0.0.1:27017: i/o timeout
	session.SetSocketTimeout(1 * time.Minute)

	switch consistency {
	case eventual:
		session.SetMode(mgo.Eventual, true)
	case monotonic:
		session.SetMode(mgo.Monotonic, true)
	case strong:
		session.SetMode(mgo.Strong, true)
	}

	return
}

func (m *dbMongo) checkGetSession(consistency mode) *mgo.Session {
	m.sessMu.RLock()
	defer m.sessMu.RUnlock()

	return m.session[consistency]

}

func (m *dbMongo) initSession(consistency mode) (*mgo.Session, error) {
	m.sessMu.Lock()
	defer m.sessMu.Unlock()
	//fmt.Println("CCCCCC", m.session)

	if m.session[consistency] != nil {
		return m.session[consistency], nil
	} else {
		s, err := dialConsistency(m.dialInfo, consistency)
		if err != nil {
			return nil, err
		} else {
			m.session[consistency] = s
			return m.session[consistency], nil
		}
	}
}

func (m *dbMongo) getSession(consistency mode) (*mgo.Session, error) {
	if s := m.checkGetSession(consistency); s != nil {
		return s, nil
	} else {
		return m.initSession(consistency)
	}
}

func (m *dbMongo) Close() error {
	return nil
}
