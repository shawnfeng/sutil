// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"testing"

	"gitlab.pri.ibanyu.com/middleware/seaweed/xlog"
)

func TestShowLog(t *testing.T) {
	t0(t)
	t1(t)
}

func t0(t *testing.T) {

	// 直接用xlog.Error, skip:4, slog.Error, skip:5, slog/slog.Error, skip:6
	xlog.SetAppLogSkip(5)
	Tracef("Tracef %s", "TT")
	Debugf("Debugf %s", "TT")
	Infof("Infof %s", "TT")
	Warnf("Warnf %s", "TT")
	Errorf("Errorf %s", "TT")
	//Fatalf("Fatalf %s", "TT")
	//Panicf("Panicf %s", "TT")

	Traceln("Traceln tt")
	Debugln("Debugln tt")
	Infoln("Infoln tt")
	Warnln("Warnln tt")
	Errorln("Errorln tt")
	//Fatalln("Fatalln tt")
	//Panicln("Panic %s", "TT")

	Infoln("FF")
	Infoln("FF")
}

func t1(t *testing.T) {

	Init("./log", "tt", "TRACE")
	Init("./log", "tt", "TRACE")
	Init("./log", "tt", "TRACE")

	Infoln("log file")

	Init("./log", "tt2", "TRACE")

	Infoln("log file2")

	Init("", "", "TRACE")

	Infoln("std out")

}
