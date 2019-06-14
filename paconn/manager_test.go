// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package paconn

import (
	"context"
	"testing"
	"fmt"
	"time"

	"github.com/shawnfeng/sutil/slog"
)

func newagentcb(ag *Agent) {
	fun := "newagentcb"

	slog.Infof(context.TODO(), "%s ag:%s", fun, ag)

	err := ag.Oneway(0, []byte("Hello Fuck You"), time.Millisecond*100)
	if err != nil {
		slog.Errorln(context.TODO(), err)
	}

}

func TestMan(t *testing.T) {
	fun := "TestMan"

	agm, err := NewAgentManager(
		":",
		time.Second * 60 *15,
		0,
		newagentcb,
		serverNotifyOneway,
		serverNotify,
		serverClose,

	)


	slog.Infof(context.TODO(), "%s %s %v", fun, agm.Listenport(), err)

	ag, err := NewAgentFromAddr(
		fmt.Sprintf("127.0.0.1:%s", agm.Listenport()),
		time.Second * 60 *15,
		time.Second * 5,
		clientNotifyOneway,
		clientNotify,
		clientClose,
	)

	if err != nil {
		t.Errorf("%s Dial err:%s", fun, err)
	}


	slog.Infoln(context.TODO(), ag)


	err = ag.Oneway(0, []byte("NT"), time.Millisecond*100)
	if err != nil {
		slog.Infoln(context.TODO(), err)
		t.Errorf("%s oneway %s", fun, err)
	}

	slog.Infof(context.TODO(), "%s ^^^^^^^^^^^^^^^^ oneway", fun)
	btype, res, err := ag.Twoway(2, []byte("NT"), time.Millisecond*100)
	if err != nil {
		slog.Warnln(context.TODO(), err)
		t.Errorf("%s twoway %s", fun, err)
	}

	if btype != 3 {
		t.Errorf("%s twoway rv btype", fun)
	}


	slog.Infof(context.TODO(), "%s twoway btype:%d res:%s", fun, btype, res)

	ag.Close()

	time.Sleep(time.Second * time.Duration(5))
}
