package slog

import "testing"

func TestShowLog(t *testing.T) {
	Init("", "", "TRACE")

	t0(t)
}


func t0(t *testing.T) {
	Infoln("FF")

}
