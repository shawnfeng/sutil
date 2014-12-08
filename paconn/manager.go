package paconn

import (
	"net"
	"time"
	"errors"

	"github.com/shawnfeng/sutil/snetutil"
	"github.com/shawnfeng/sutil/slog"
)

type AgentManager struct {
	agents map[string]*Agent

	cbNewagent func(*Agent)

	cbOnewaynotify FunAgOnewaynotify
	cbTwowaynotify FunAgTwowaynotify
	cbClose FunAgClose
	

	addrListen net.Addr

}

func (m *AgentManager) Agents () map[string]*Agent {
	return m.agents
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


	if _, ok := m.agents[a.Id()]; ok {
		delete(m.agents, a.Id())
		if m.cbClose != nil {
			m.cbClose(a, pack, err)
		}
	} else {
		slog.Errorf("%s delete not find:%s", fun, a)
	}

}


func (m *AgentManager) Oneway(aid string, btype int32, data []byte, timeout time.Duration) error {

	if a, ok := m.agents[aid]; ok {
		return a.Oneway(btype, data, timeout)
	} else {
		return errors.New("agent id not found");
	}

}


func (m *AgentManager) Twoway(aid string, btype int32, data []byte, timeout time.Duration) (int32, []byte, error) {
	if a, ok := m.agents[aid]; ok {
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
			m.agents[ag.Id()] = ag

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
