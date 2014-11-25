package stime

import (
	"time"
//	"fmt"

)


var (
	Since2014 int64 = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1000
)


func Timestamp2014() uint64 {
	return uint64(time.Now().UnixNano()/1000 - Since2014)

}


type runTimeStat struct {
	//logkey string
	since time.Time
}

//func (m *runTimeStat) StatLog() string {
//	return fmt.Sprintf("%s RUNTIME:%d", m.logkey, m.Duration())
//}

func (m *runTimeStat) Millisecond() int64 {
	return m.Microsecond() / 1000
}

func (m *runTimeStat) Microsecond() int64 {
	return m.Duration().Nanoseconds() / 1000

}

func (m *runTimeStat) Nanosecond() int64 {
	return m.Duration().Nanoseconds()
}

func (m *runTimeStat) Duration() time.Duration {
	return time.Since(m.since)
}


func (m *runTimeStat) Reset() {
	m.since = time.Now()
}



//func NewTimeStat(key string) *runTimeStat {
func NewTimeStat() *runTimeStat {
	return &runTimeStat {
		//logkey: key,
		since: time.Now(),
	}
}
