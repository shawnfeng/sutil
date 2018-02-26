// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package stime

import (
	"time"
	//"fmt"

)

func DayBeginStamp(now int64) int64 {

	_, offset := time.Now().Zone()
	//fmt.Println(zone, offset)
	return now - (now+int64(offset)) % int64(3600 * 24)
	//return (now + int64(offset))/int64(3600 * 24) * int64(3600 * 24) - int64(offset)

}


func HourBeginStamp(now int64) int64 {

	_, offset := time.Now().Zone()
	//fmt.Println(zone, offset)
	return now - (now+int64(offset)) % int64(3600)
	//return (now + int64(offset))/int64(3600 * 24) * int64(3600 * 24) - int64(offset)

}


// 获取指定天的时间范围
// 天格式 2006-01-02
// 为空时候返回当天的
func DayBeginStampFromStr(day string) (int64, error) {
	nowt := time.Now()
	now := nowt.Unix()

	var begin int64
	if len(day) > 0 {
		tm, err := time.ParseInLocation("2006-01-02", day, nowt.Location())
		if err != nil {
			return 0, err
		}

		begin = tm.Unix()

	} else {
		begin = DayBeginStamp(now)

	}

	return begin, nil

}

func WeekScope(stamp int64) (int64, int64) {
    now := time.Unix(stamp, 0)
    weekday := time.Duration(now.Weekday())
    if weekday == 0 {
        weekday = 7
    }
    year, month, day := now.Date()
	currentZeroDay:= time.Date(year, month, day, 0, 0, 0, 0, time.Local)
    begin := currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour).Unix()
    return begin, begin+24*3600*7-1
}


func MonthScope(stamp int64) (int64, int64) {
	now := time.Unix(stamp, 0)
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return firstOfMonth.Unix(), lastOfMonth.Unix()+3600*24-1
}



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
