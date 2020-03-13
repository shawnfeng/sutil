// Copyright 2014 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	defaultSyncTimeout   = time.Second * 10
	defaultSocketTimeout = time.Second * 10
	defaultPoolLimit     = 128
	defaultDialTimeout   = time.Second * 5
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
		timeout = defaultDialTimeout
	}

	info := &mgo.DialInfo{
		Addrs:     addrs,
		Timeout:   timeout,
		Database:  dbname,
		Username:  user,
		Password:  passwd,
		PoolLimit: defaultPoolLimit,
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
	session.SetSyncTimeout(defaultSyncTimeout)
	session.SetSocketTimeout(defaultSocketTimeout)

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
