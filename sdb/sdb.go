package sdb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/shawnfeng/sutil/slog"
)


type MgoDb struct {
	addr string
	mgoSession *mgo.Session
}

func (m *MgoDb) LoadDb(addr string) error {
	fun := "MgoDb.LoadDb"
	if m.mgoSession != nil {
		m.mgoSession.Close()
		slog.Warnf("%s old mongodb load close", fun)
	}

	session, err := mgo.Dial(addr)
	if err != nil {
		slog.Warnf("%s mongodb load err:%s", fun, err)
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	m.mgoSession = session

	slog.Infof("%s load mongo:%s", fun, addr)
	return nil

}



