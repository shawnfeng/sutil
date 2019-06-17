// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package paconn

import (
	"context"
	"fmt"
	"net"
	"time"
	"testing"

	"github.com/shawnfeng/sutil/snetutil"
	"github.com/shawnfeng/sutil/slog/slog"
)


func serverNotify(a *Agent, btype int32, data []byte) (int32, []byte) {

	fun := "serverNotify"
	slog.Infof(context.TODO(), "%s %s %d %v %s", fun, a, btype, data, data)

	if "NT" == string(data) {
		slog.Infof(context.TODO(), ">>>>%s use not timeout", fun)
		return btype+1, []byte("OK")
	} else {
		slog.Infof(context.TODO(), "%s use timeout", fun)
		time.Sleep(time.Second * time.Duration(1))
		return btype+1, []byte("TIMEOUT")
	}



}


func serverNotifyOneway(a *Agent, btype int32, data []byte) {

	fun := "serverNotifyOneway"
	slog.Infof(context.TODO(), "%s %s %v %s", fun, a, data, data)


}


func serverClose(a *Agent, data []byte, err error) {
	fun := "serverClose"

	slog.Infof(context.TODO(), "%s %s %v %s %s", fun, a, data, data, err)
}

func WaitLink(t *testing.T) string {
	fun := "WaitLink"
	addr := fmt.Sprintf(":%d", 0)
	//addr := ":0"

	tcpAddr, error := net.ResolveTCPAddr("tcp", addr)

	slog.Infof(context.TODO(), "%s %s %v %s %s", fun, tcpAddr, error, tcpAddr.Network(), tcpAddr.String())

	if error != nil {
		slog.Panicf(context.TODO(), "%s Error: Could not resolve address %s", fun, error)
	}


	netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String())

	slog.Infof(context.TODO(), "%s listen:%s", fun, netListen.Addr())
	if error != nil {
		slog.Panicf(context.TODO(), "%s Error: Could not Listen %s", fun, error)

	}


	addr = netListen.Addr().String()

	port := snetutil.IpAddrPort(addr)
	slog.Infoln(context.TODO(), port)


	usetestaddrPort := fmt.Sprintf("%s:%s", "127.0.0.1", port)
	go func() {
		defer netListen.Close()
	for {
		//slog.Infof(context.TODO(), "%s Waiting for clients", fun)
		conn, error := netListen.Accept()
		if error != nil {
			slog.Warnf(context.TODO(), "%s Agent error: %s", fun, error)
			t.Errorf("%s", error)
			return
		} else {

			ag := NewAgent(
				conn,
				time.Second * 60 * 15,
				0,
				serverNotifyOneway,
				serverNotify,
				serverClose,
		
			)

			slog.Infoln(context.TODO(), "S:", ag)

		}
	}
	}()

	return usetestaddrPort

}

func clientNotifyOneway(a *Agent, btype int32, data []byte) {

	fun := "clientNotifyOneway"
	slog.Infof(context.TODO(), "%s %s %d %v %s", fun, a, btype, data, data)

}


func clientNotify(a *Agent, btype int32, data []byte) (int32, []byte) {

	fun := "clientNotify"
	slog.Infof(context.TODO(), "%s %s %d %v %s", fun, a, btype, data, data)

	return btype+1, []byte("OK")

}


func clientClose(a *Agent, data []byte, err error) {
	fun := "clientClose"

	slog.Infof(context.TODO(), "%s %s %v %s %s", fun, a, data, data, err)
}


func clientAgent(t *testing.T, addrport string) {

	fun := "clientAgent"

	ag, err := NewAgentFromAddr(
		addrport,
		time.Second * 60 * 15,
		time.Second * 5,
		clientNotifyOneway,
		clientNotify,
		clientClose,
	)

	if err != nil {
		t.Errorf("%s Dial err:%s ag:%s", fun, err, ag)
		return
	}


	slog.Infoln(context.TODO(), ag)


	err = ag.Oneway(1, []byte("NT"), time.Millisecond*100)
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
		t.Errorf("%s twoway rv btype:%d ", fun, btype)
	}

	slog.Infof(context.TODO(), "%s twoway btype:%d res:%s", fun, btype, res)

	ag.Close()

}

func TestAgent(t *testing.T) {

	addrport := WaitLink(t)
	time.Sleep(time.Millisecond * time.Duration(100))

	clientAgent(t, addrport)

	time.Sleep(time.Second * time.Duration(5))
}

