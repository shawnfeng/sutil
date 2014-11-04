package paconn

import (
	"fmt"
	"net"
	"time"
	"testing"

	"sutil/snetutil"
	"sutil/slog"
)

var usetestaddrPort string

func serverNotify(a *Agent, data []byte) []byte {

	fun := "serverNotify"
	slog.Infof("%s %s %v %s", fun, a, data, data)

	if "NT" == string(data) {
		slog.Infof("%s use not timeout", fun)
		return []byte("OK")
	} else {
		slog.Infof("%s use timeout", fun)
		time.Sleep(time.Second * time.Duration(1))
		return []byte("TIMEOUT")
	}



}


func serverNotifyOneway(a *Agent, data []byte) {

	fun := "serverNotifyOneway"
	slog.Infof("%s %s %v %s", fun, a, data, data)


}


func serverClose(a *Agent, data []byte, err error) {
	fun := "serverClose"

	slog.Infof("%s %s %v %s %s", fun, a, data, data, err)
}

func WaitLink(t *testing.T) {
	fun := "WaitLink"
	addr := fmt.Sprintf(":%d", 0)
	//addr := ":0"

	tcpAddr, error := net.ResolveTCPAddr("tcp", addr)

	slog.Infof("%s %s %v %s %s", fun, tcpAddr, error, tcpAddr.Network(), tcpAddr.String())

	if error != nil {
		slog.Panicf("%s Error: Could not resolve address %s", fun, error)
	}


	netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String())

	slog.Infof("%s listen:%s", fun, netListen.Addr())
	if error != nil {
		slog.Panicf("%s Error: Could not Listen %s", fun, error)

	}
	defer netListen.Close()

	addr = netListen.Addr().String()

	port := snetutil.IpAddrPort(addr)
	slog.Infoln(port)


	usetestaddrPort = fmt.Sprintf("%s:%s", "127.0.0.1", port)
	for {
		//slog.Infof("%s Waiting for clients", fun)
		conn, error := netListen.Accept()
		if error != nil {
			slog.Warnf("%s Agent error: ", fun, error)
		} else {

			id, ag := NewAgent(
				conn,
				0,
				serverNotifyOneway,
				serverNotify,
				serverClose,
		
			)

			slog.Infoln("S:", id, ag)

		}
	}

}

func clientNotifyOneway(a *Agent, data []byte) {

	fun := "clientNotifyOneway"
	slog.Infof("%s %s %v %s", fun, a, data, data)

}


func clientNotify(a *Agent, data []byte) []byte {

	fun := "clientNotify"
	slog.Infof("%s %s %v %s", fun, a, data, data)

	return []byte("OK")

}


func clientClose(a *Agent, data []byte, err error) {
	fun := "clientClose"

	slog.Infof("%s %s %v %s %s", fun, a, data, data, err)
}


func clientAgent(t *testing.T) {

	fun := "clientAgent"

	id, ag, err := NewAgentFromAddr(
		usetestaddrPort,
		1000 * 5,
		clientNotifyOneway,
		clientNotify,
		clientClose,
	)

	if err != nil {
		t.Errorf("%s Dial err:%s", fun, err)
	}


	slog.Infoln(id, ag)


	err = ag.Oneway([]byte("NT"), 100)
	if err != nil {
		slog.Infoln(err)
		t.Errorf("%s oneway %s", fun, err)
	}

	slog.Infof("%s ^^^^^^^^^^^^^^^^ oneway", fun)
	res, err := ag.Twoway([]byte("NT"), 100)
	if err != nil {
		slog.Warnln(err)
		t.Errorf("%s twoway %s", fun, err)
	}

	slog.Infof("%s twoway res:%s", fun, res)

	ag.Close()

}

func TestAgent(t *testing.T) {

	go WaitLink(t)
	time.Sleep(time.Millisecond * time.Duration(100))

	clientAgent(t)

	time.Sleep(time.Second * time.Duration(5))
}

