package stime

import (
	"time"
	"fmt"

)


var (
	Since2014 int64 = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1000
)


func Timestamp2014() uint64 {
	return uint64(time.Now().UnixNano()/1000 - Since2014)

}


type runTimeStat struct {
	logkey string
	stamp int64
}

func (m *runTimeStat) StatLog() string {
	return fmt.Sprintf("%s RUNTIME:%d", m.logkey, m.Nanosecond())
}

func (m *runTimeStat) Millisecond() int64 {
	return m.Microsecond() / 1000
}

func (m *runTimeStat) Microsecond() int64 {
	return m.Nanosecond() / 1000

}

func (m *runTimeStat) Nanosecond() int64 {
	return time.Now().UnixNano()-m.stamp
}




func NewTimeStat(key string) *runTimeStat {
	return &runTimeStat {
		logkey: key,
		stamp: time.Now().UnixNano(),
	}
}
