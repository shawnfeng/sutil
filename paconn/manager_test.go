package paconn

import (
	"testing"
	"fmt"
	"time"

	"github.com/shawnfeng/sutil/slog"
)

func newagentcb(ag *Agent) {
	fun := "newagentcb"

	slog.Infof("%s ag:%s", fun, ag)

	err := ag.Oneway([]byte("Hello Fuck You"), time.Millisecond*100)
	if err != nil {
		slog.Errorln(err)
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


	slog.Infof("%s %s %v", fun, agm.Listenport(), err)

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


	slog.Infoln(ag)


	err = ag.Oneway([]byte("NT"), time.Millisecond*100)
	if err != nil {
		slog.Infoln(err)
		t.Errorf("%s oneway %s", fun, err)
	}

	slog.Infof("%s ^^^^^^^^^^^^^^^^ oneway", fun)
	res, err := ag.Twoway([]byte("NT"), time.Millisecond*100)
	if err != nil {
		slog.Warnln(err)
		t.Errorf("%s twoway %s", fun, err)
	}

	slog.Infof("%s twoway res:%s", fun, res)

	ag.Close()

	time.Sleep(time.Second * time.Duration(5))
}
