// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package paconn


import (
	"net"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"errors"


	"github.com/sdming/gosnow"
	"github.com/golang/protobuf/proto"

	"github.com/shawnfeng/sutil"
	"github.com/shawnfeng/sutil/slog"
	"github.com/shawnfeng/sutil/stime"
	"github.com/shawnfeng/sutil/snetutil"


	"github.com/shawnfeng/sutil/paconn/pb"
)

var msgidgen *gosnow.SnowFlake

const (
	DEFAULT_SEND_TIMEOUT time.Duration = time.Millisecond*200
)


func init() {
	gosnow.Since = stime.Since2014 / 1000
	v, err := gosnow.NewSnowFlake(0)
	if err != nil {
		panic("snowflake init error, msgid can not get!")
	}

	msgidgen = v
}

type ackNotify struct {
	err error
	busstype int32
	result []byte
}

type FunAgOnewaynotify func(*Agent, int32, []byte)
type FunAgTwowaynotify func(*Agent, int32, []byte) (int32, []byte)
type FunAgClose func(*Agent, []byte, error)


type Agent struct {
	id string
	sendLock sync.Mutex
	conn net.Conn
	tuple4 string

	callmsgLock sync.Mutex
	callmsg  map[uint64] chan *ackNotify

	readTimeout time.Duration
	heartIntv time.Duration
	isConn int32


	cbOnewaynotify FunAgOnewaynotify
	cbTwowaynotify FunAgTwowaynotify
	cbClose FunAgClose

}

func (m *Agent) String() string {
	return fmt.Sprintf("%s@%s", m.id, m.tuple4)
}

func (m *Agent) Id() string {
	return m.id
}


func (m *Agent) Close() {
	// 不管错误，即使关闭了，也再关一次，修改isConn状态
	m.conn.Close()
	//err := m.conn.Close()
	//if err != nil {
	//	slog.Warnf("Agent.Close Close err:%s", err)
	//}

	atomic.StoreInt32(&m.isConn, 0)
}

func (m *Agent) Oneway(btype int32, data []byte, timeout time.Duration) error {
	pb := &connproto.ConnProto {
		Type: connproto.ConnProto_CALL.Enum(),
		Busstype: proto.Int32(btype),
		Bussdata: data,

	}


	spb, _ := proto.Marshal(pb)
	
	return m.send(spb, timeout)
}


func (m *Agent) Twoway(btype int32, data []byte, timeout time.Duration) (int32, []byte, error) {
	//fun := "Agent.Twoway"

	msgid, _ := msgidgen.Next()
	pb := &connproto.ConnProto {
		Type: connproto.ConnProto_CALL.Enum(),
		Msgid: proto.Uint64(msgid),
		Busstype: proto.Int32(btype),
		Bussdata: data,
	}


	spb, _ := proto.Marshal(pb)
	
	done := make(chan *ackNotify)

	func () {
		m.callmsgLock.Lock()
		defer m.callmsgLock.Unlock()
		m.callmsg[msgid] = done
	}()

	defer func () {
		m.callmsgLock.Lock()
		defer m.callmsgLock.Unlock()
		delete(m.callmsg, msgid)
	}()

	st := stime.NewTimeStat()
	err := m.send(spb, timeout)
	if err != nil {
		return 0, nil, err
	}

	senduse := st.Duration()

	if senduse >= timeout {
		return 0, nil, errors.New(fmt.Sprintf("call send timetout:%d", senduse))
	}


	select {
	case v := <-done:
		return v.busstype, v.result, v.err
	case <-time.After(timeout-senduse):
		m.Close()
		return 0, nil, errors.New("call ack timetout")
	}

}

func (m *Agent) recvACK(pb *connproto.ConnProto) {
	fun := "Agent.recvACK"
	msgid := pb.GetAckmsgid()
	btype := pb.GetBusstype()


	c, ok := func() (chan *ackNotify, bool) {
		m.callmsgLock.Lock()
		defer m.callmsgLock.Unlock()
		cc, o := m.callmsg[msgid]
		return cc, o
	}()


	if ok {
		an := &ackNotify {
			err: nil,
			busstype: btype,
			result: pb.GetBussdata(),
		}
		select {
		case c <-an:
		default:
			slog.Warnf("%s agent:%s msgid:%d no wait notify", fun, m, msgid)
		}
	} else {
		slog.Warnf("%s agent:%s msgid:%d not found", fun, m, msgid)
	}

}

func (m *Agent) recvCALL(pb *connproto.ConnProto) {
	fun := "Agent.recvCALL"
	data := pb.GetBussdata()
	btype := pb.GetBusstype()
	msgid := pb.GetMsgid()

	if msgid != 0 {
		res := make([]byte, 0)
		var rbtype int32 = 0
		if m.cbTwowaynotify != nil {
			rbtype, res = m.cbTwowaynotify(m, btype, data)
		}

		// 需要回执
		ack := &connproto.ConnProto {
			Type: connproto.ConnProto_ACK.Enum(),
			Ackmsgid: proto.Uint64(msgid),
			Busstype: proto.Int32(rbtype),
			Bussdata: res,
		}

		sdata, _ := proto.Marshal(ack)
		err := m.send(sdata, DEFAULT_SEND_TIMEOUT)
		if err != nil {
			slog.Warnf("%s agent:%s ack error:%s", fun, m, err)
		}

	} else {
		if m.cbOnewaynotify != nil {
			m.cbOnewaynotify(m, btype, data)
		}

	}


}

func (m *Agent) proto(data []byte) {
	fun := "Agent.proto"
	pb := &connproto.ConnProto{}
	err := proto.Unmarshal(data, pb)
	if err != nil {
		m.Close()
		slog.Warnf("%s agent:%s unmarshaling error: %s data:%v sd:%s", fun, m, err, data, data)
		return
	}

	slog.Infof("%s a:%s %s", fun, m, pb)


	pb_type := pb.GetType()
	if pb_type == connproto.ConnProto_ACK {
		m.recvACK(pb)

	} else if pb_type == connproto.ConnProto_CALL {
		m.recvCALL(pb)
	} else if pb_type == connproto.ConnProto_HEART {
		m.recvHEART()
	} else {
		m.Close()
		slog.Warnf("%s agent:%s type error: %s data:%v sd:%s", fun, m, err, data, data)
	}

}

func (m *Agent) send(data []byte, timeout time.Duration) error {
	if atomic.LoadInt32(&m.isConn) == 0 {
		return errors.New("connection is not ok")

	}

	s := snetutil.Packdata(data)

	m.sendLock.Lock()
	defer m.sendLock.Unlock()
	m.conn.SetWriteDeadline(time.Now().Add(timeout))
	a, err := m.conn.Write(s)
	//slog.Infof("%s agent:%s Send Write %d rv %d", fun, m, len(s), a)

	if err != nil {
		m.Close()
		return errors.New(fmt.Sprintf("send write error:%s", err))
	}

	if len(s) != a {
		m.Close()
		return errors.New("send write error:len")
	}

	return nil
}

func (m *Agent) sendHEART() {
	fun := "Agent.sendHEART"
	heart := &connproto.ConnProto {
		Type: connproto.ConnProto_HEART.Enum(),
	}

	//slog.Debugf("%s agent:%s msg:%s", fun, m, heart)

	data, _ := proto.Marshal(heart)
	err := m.send(data, DEFAULT_SEND_TIMEOUT)

	if err != nil {
		slog.Warnf("%s agent:%s error:%s", fun, m, err)
	}



}


func (m *Agent) recvHEART() {
	if m.heartIntv <= 0 {
		// 被动心跳
		m.sendHEART();
	}
}


func (m *Agent) heart() {
	fun := "Agent.heart"
	if m.heartIntv > 0 {
		// 主动心跳
		slog.Infof("%s agent:%s heart:%d", fun, m, m.heartIntv)
		ticker := time.NewTicker(m.heartIntv)
		for {
			select {
			case <-ticker.C:
				if atomic.LoadInt32(&m.isConn) != 0 {
					m.sendHEART();
				}
			}
		}

	} else {
		slog.Infof("%s agent:%s noheart:%d", fun, m, m.heartIntv)
	}

}

func (m *Agent) recv() {

	// 是否是read返回错误socket已经关闭，返回时候没有处理的数据，错误信息
	isclose, data, err := snetutil.PackageSplit(m.conn, m.readTimeout, m.proto)

	if !isclose {
		m.Close()
	}

	if m.cbClose != nil {
		// 所有关闭回调放到这里就好了，关闭时候，其他地方Close会走到这里
		m.cbClose(m, data, err)
	}

}

func NewAgent(
	c net.Conn,
	readto time.Duration,
	heart time.Duration,
	onenotify FunAgOnewaynotify,
	twonotify FunAgTwowaynotify,
	close FunAgClose,
) *Agent {
	fun := "NewAgent"

	if readto <= 0 {
		// 15分
		readto = (1000*60*5) * 3
	}

	aid, err := sutil.GetUniqueMd5()
	if err != nil {
		slog.Errorf("%s new err:%s", fun, err)
		return nil
	}

	a := &Agent {
		id: aid,
		conn: c,
		tuple4: fmt.Sprintf("%s-%s", c.LocalAddr().String(), c.RemoteAddr().String()),
		callmsg: make(map[uint64] chan *ackNotify),
		readTimeout: readto,
		heartIntv: heart,
		isConn: 1,
		cbOnewaynotify: onenotify,
		cbTwowaynotify: twonotify,
		cbClose:close,
	}

	go a.recv()
	go a.heart()

	slog.Infof("%s a:%s", fun, a)

	return a
}

func NewAgentFromAddr(addr string,
	readtimeout time.Duration,
	heart time.Duration,
	onenotify FunAgOnewaynotify,
	twonotify FunAgTwowaynotify,
	close FunAgClose,
) (*Agent, error) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}


	a := NewAgent(
		conn,
		readtimeout,
		heart,
		onenotify,
		twonotify,
		close,
	)

	return a, nil

}

