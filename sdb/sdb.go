// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package sdb

import (
	"context"
	"github.com/shawnfeng/sutil/slog"
	"gopkg.in/mgo.v2"
)


type MgoDb struct {
	addr string
	mgoSession *mgo.Session
}

func (m *MgoDb) LoadDb(addr string) error {
	fun := "MgoDb.LoadDb"
	if m.mgoSession != nil {
		m.mgoSession.Close()
		slog.Warnf(context.TODO(), "%s old mongodb load close", fun)
	}

	session, err := mgo.Dial(addr)
	if err != nil {
		slog.Warnf(context.TODO(), "%s mongodb load err:%s", fun, err)
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	m.mgoSession = session

	slog.Infof(context.TODO(), "%s load mongo:%s", fun, addr)
	return nil

}



