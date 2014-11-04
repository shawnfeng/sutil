package paconn


import (
	"net"
	"fmt"
	"sync"
	"time"
	"errors"
	"crypto/md5"

	"code.google.com/p/go-uuid/uuid"
	"github.com/sdming/gosnow"
	"code.google.com/p/goprotobuf/proto"

	"sutil/slog"
	"sutil/stime"
	"sutil/snetutil"


	"sutil/paconn/pb"
)

var msgidgen *gosnow.SnowFlake


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
	result []byte
}

type Agent struct {
	id string
	sendLock sync.Mutex
	conn net.Conn
	tuple4 string
	callmsg  map[uint64] chan *ackNotify
	heartIntv int64
	isConn bool


	cbOnewaynotify func(*Agent, []byte)
	cbTwowaynotify func(*Agent, []byte) []byte
	cbClose func(*Agent, []byte, error)

}

func (m *Agent) String() string {
	return fmt.Sprintf("%s@%s", m.id, m.tuple4)
}

func (m *Agent) Close() {
	m.conn.Close()
	m.isConn = false

}

func (m *Agent) Oneway(data []byte, timeout int64) error {
	pb := &connproto.ConnProto {
		Type: connproto.ConnProto_CALL.Enum(),
		Bussdata: data,
	}


	spb, _ := proto.Marshal(pb)
	
	return m.send(spb, timeout)
}


func (m *Agent) Twoway(data []byte, timeout int64) ([]byte, error) {
	fun := "Agent.Twoway"

	msgid, _ := msgidgen.Next()
	pb := &connproto.ConnProto {
		Type: connproto.ConnProto_CALL.Enum(),
		Msgid: proto.Uint64(msgid),
		Bussdata: data,
	}


	spb, _ := proto.Marshal(pb)
	
	done := make(chan *ackNotify)

	m.callmsg[msgid] = done

	defer delete(m.callmsg, msgid)

	st := stime.NewTimeStat(fun)
	err := m.send(spb, timeout)
	if err != nil {
		return nil, err
	}

	senduse := st.Millisecond()

	if senduse >= timeout {
		return nil, errors.New(fmt.Sprintf("call send timetout:%d", senduse))
	}


	select {
	case v := <-done:
		return v.result, v.err
	case <-time.After(time.Millisecond * time.Duration(timeout-senduse)):
		return nil, errors.New("call ack timetout")
	}

}

func (m *Agent) recvACK(pb *connproto.ConnProto) {
	fun := "Agent.recvACK"
	msgid := pb.GetAckmsgid()

	if c, ok := m.callmsg[msgid]; ok {
		an := &ackNotify {
			err: nil,
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
	if m.cbTwowaynotify != nil {
		msgid := pb.GetMsgid()
		if msgid != 0 {
			res := m.cbTwowaynotify(m, data)
			// 需要回执
			ack := &connproto.ConnProto {
				Type: connproto.ConnProto_ACK.Enum(),
				Ackmsgid: proto.Uint64(msgid),
				Bussdata: res,
			}

			sdata, _ := proto.Marshal(ack)
			err := m.send(sdata, 200)
			if err != nil {
				slog.Warnf("%s agent:%s ack error:%s", fun, m, err)
			}

		} else {
			m.cbOnewaynotify(m, data)
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

func (m *Agent) send(data []byte, timeout int64) error {
	if !m.isConn {
		return errors.New("connection is not ok")

	}

	s := snetutil.Packdata(data)

	m.sendLock.Lock()
	defer m.sendLock.Unlock()

	m.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
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
	err := m.send(data, 200)

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
		slog.Infof("%s agent:%s heart:%d ms", fun, m, m.heartIntv)
		ticker := time.NewTicker(time.Millisecond * time.Duration(m.heartIntv))
		for {
			select {
			case <-ticker.C:
				if m.isConn {
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
	isclose, data, err := snetutil.PackageSplit(m.conn, 300 * 3, m.proto)

	if !isclose {
		m.Close()
	}

	if m.cbClose != nil {
		// 所有关闭回调放到这里就好了，关闭时候，其他地方Close会走到这里
		m.cbClose(m, data, err)
	}

}


func NewAgent(c net.Conn,
	heart int64,
	onenotify func(*Agent, []byte),
	twonotify func(*Agent, []byte) []byte,
	close func(*Agent, []byte, error),
) (string, *Agent) {
	uuidgen := uuid.NewUUID()
	aid := uuidgen.String()

	h := md5.Sum([]byte(aid))
	aid = fmt.Sprintf("%x", h)

	a := &Agent {
		id: aid,
		conn: c,
		tuple4: fmt.Sprintf("%s-%s", c.LocalAddr().String(), c.RemoteAddr().String()),
		callmsg: make(map[uint64] chan *ackNotify),
		heartIntv: heart,
		isConn: true,
		cbOnewaynotify: onenotify,
		cbTwowaynotify: twonotify,
		cbClose:close,
	}

	go a.recv()
	go a.heart()

	return aid, a
}

// 是否开启心跳
