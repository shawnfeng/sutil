package slog

import "testing"

func TestShowLog(t *testing.T) {
	t0(t)
	t1(t)
}


func t0(t *testing.T) {
	Infoln("FF")

}


func t1(t *testing.T) {

	Init("./log", "tt", "TRACE")


	Infoln("log file")


	Init("./log", "tt2", "WARN")


	Infoln("log file2")

}
