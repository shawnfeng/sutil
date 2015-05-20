// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package paconn

import (
	"net"
	"time"
	"errors"
	"sync"

	"github.com/shawnfeng/sutil/snetutil"
	"github.com/shawnfeng/sutil/slog"
)

type AgentManager struct {

	agentMu sync.Mutex
	agents map[string]*Agent

	cbNewagent func(*Agent)

	cbOnewaynotify FunAgOnewaynotify
	cbTwowaynotify FunAgTwowaynotify
	cbClose FunAgClose
	

	addrListen net.Addr

}

func (m *AgentManager) Agents () map[string]*Agent {
	m.agentMu.Lock()
	defer m.agentMu.Unlock()

	rv := make(map[string]*Agent)

	for k, v := range m.agents {
		rv[k] = v
	}

	return rv
}


func (m *AgentManager) Listenport () string {
	return snetutil.IpAddrPort(m.addrListen.String())
}

func (m *AgentManager) callbackOneway (a *Agent, btype int32, recv []byte) {

	if m.cbOnewaynotify != nil {
		m.cbOnewaynotify(a, btype, recv)
	}

}

func (m *AgentManager) callbackTwoway (a *Agent, btype int32, recv []byte) (int32, []byte) {
	if m.cbTwowaynotify != nil {
		return m.cbTwowaynotify(a, btype, recv)
	}

	return 0, nil
}

func (m *AgentManager) callbackClose (a *Agent, pack []byte, err error) {
	fun := "AgentManager.callbackClose"

	slog.Infof("%s close:%s pack:%v", fun, a, pack)

	ok := func() bool {
		m.agentMu.Lock()
		defer m.agentMu.Unlock()
		_, o := m.agents[a.Id()]
		if o {
			delete(m.agents, a.Id())
		}
		return o
	}()


	if ok && m.cbClose != nil {
		m.cbClose(a, pack, err)

	} else {
		slog.Errorf("%s delete not find:%s", fun, a)
	}

}

func (m *AgentManager) getAgent(aid string) (*Agent, bool) {
	m.agentMu.Lock()
	defer m.agentMu.Unlock()
	a, ok := m.agents[aid]
	return a, ok
}

func (m *AgentManager) addAgent(a *Agent) {
	m.agentMu.Lock()
	defer m.agentMu.Unlock()
	m.agents[a.Id()] = a
}



func (m *AgentManager) Oneway(aid string, btype int32, data []byte, timeout time.Duration) error {

	if a, ok := m.getAgent(aid); ok {
		return a.Oneway(btype, data, timeout)
	} else {
		return errors.New("agent id not found");
	}

}


func (m *AgentManager) Twoway(aid string, btype int32, data []byte, timeout time.Duration) (int32, []byte, error) {
	if a, ok := m.getAgent(aid); ok {
		return a.Twoway(btype, data, timeout)
	} else {
		return 0, nil, errors.New("agent id not found");
	}

}



func (m *AgentManager) accept(
	done chan error,
	tcpAddr net.Addr,

	readto time.Duration,
	heart time.Duration,

) {
	fun := "AgentManager.accept"

	netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String())
	slog.Infof("%s listen:%s", fun, netListen.Addr())
	if error != nil {
		done <-error
		return;
	}
	defer netListen.Close()

	m.addrListen = netListen.Addr()
	done <-nil

	for {
		//slog.Infof("%s Waiting for clients", fun)
		conn, error := netListen.Accept()
		if error != nil {
			slog.Warnf("%s Agent error: ", fun, error)
		} else {
			ag := NewAgent(
				conn,
				readto,
				heart,
				m.callbackOneway,
				m.callbackTwoway,
				m.callbackClose,

			)
			m.addAgent(ag)

			if m.cbNewagent != nil {
				m.cbNewagent(ag)
			}
		}
	}


}

func NewAgentManager(
	addr string,

	readtimeout time.Duration,
	heart time.Duration,


	newagent func(*Agent),
	onenotify FunAgOnewaynotify,
	twonotify FunAgTwowaynotify,
	close FunAgClose,

) (*AgentManager, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	agm := &AgentManager {
		agents: make(map[string]*Agent),

		cbNewagent: newagent,
		cbOnewaynotify: onenotify,
		cbTwowaynotify: twonotify,
		cbClose:close,

	}

	done := make(chan error)

	go agm.accept(done, tcpAddr, readtimeout, heart)
	err = <-done
	if err != nil {
		return nil, err
	}

	return agm, nil

}
