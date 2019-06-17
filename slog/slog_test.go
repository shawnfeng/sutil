// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package slog

import "testing"

func TestShowLog(t *testing.T) {
	t0(t)
	t1(t)
}


func t0(t *testing.T) {

	Tracef("Tracef %s", "TT")
	Debugf("Debugf %s", "TT")
	Infof("Infof %s", "TT")
	Warnf("Warnf %s", "TT")
	Errorf("Errorf %s", "TT")
	Fatalf("Fatalf %s", "TT")
	//Panicf("Panicf %s", "TT")


	Traceln("Trace %s", "TT")
	Debugln("Debug %s", "TT")
	Infoln("Info %s", "TT")
	Warnln("Warn %s", "TT")
	Errorln("Error", "TT")
	Fatalln("Fatal", "TT")
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
