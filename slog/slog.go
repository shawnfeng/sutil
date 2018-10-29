// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"log"
	//	"io"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type logger struct {
	hourUse int64

	logpref string

	logfp *os.File
	per   *log.Logger
}

func (self *logger) resetOutput(logpref string) {
	self.logpref = logpref
	now := time.Now()
	self.hourOutput(&now)

}

func (self *logger) hourOutput(now *time.Time) {
	if self.logfp != nil {
		self.logfp.Close()
		self.logfp = nil
	}

	nx := now.Unix() + 3600
	self.hourUse = time.Unix(nx-nx%3600, 0).Unix()

	if self.logpref == "" {
		self.per = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	} else {
		hour := now.Format("2006-01-02-15")
		logFile := fmt.Sprintf("%s.%s.log", self.logpref, hour)

		logf, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err)
			return
		}

		self.per = log.New(logf, "", log.Ldate|log.Ltime|log.Lmicroseconds)
		self.logfp = logf

	}

}

func (self *logger) setOutput() {
	now := time.Now()
	if now.Unix() >= self.hourUse {
		self.hourOutput(&now)
	}

}

func (self *logger) Printf(format string, v ...interface{}) {
	self.setOutput()
	if self.per == nil {
		log.Println("slog nil")
		return
	}
	self.per.Printf(format, v...)
}

func (self *logger) Panicf(format string, v ...interface{}) {
	self.setOutput()
	if self.per == nil {
		log.Println("slog nil")
		return
	}

	self.per.Panicf(format, v...)
}

func (self *logger) Println(v ...interface{}) {
	self.setOutput()
	if self.per == nil {
		log.Println("slog nil")
		return
	}

	self.per.Println(v...)
}

func (self *logger) Panicln(v ...interface{}) {
	self.setOutput()
	if self.per == nil {
		log.Println("slog nil")
		return
	}

	self.per.Panicln(v...)
}

// log 级别
const (
	LV_TRACE int = 0
	LV_DEBUG int = 1
	LV_INFO  int = 2
	LV_WARN  int = 3
	LV_ERROR int = 4
	LV_FATAL int = 5
	LV_PANIC int = 6
)

var (
	headTrace string
	headDebug string
	headInfo  string
	headWarn  string
	headError string
	headFatal string
	headPanic string

	headFmtTrace string
	headFmtDebug string
	headFmtInfo  string
	headFmtWarn  string
	headFmtError string
	headFmtFatal string
	headFmtPanic string

	// log count
	cnTrace int64
	cnDebug int64
	cnInfo  int64
	cnWarn  int64
	cnError int64
	cnFatal int64
	cnPanic int64
	// log count stat stamp
	cnStamp int64

	slogMutex sync.Mutex
	log_level int
	lg        *logger

	logs []string
)

func addLogs(log string) {
	logs = append(logs, log)
	if len(logs) > 10 {
		logs = logs[len(logs)-10:]
	}
}

func getLogs() []string {

	slogMutex.Lock()
	defer slogMutex.Unlock()

	tmp := make([]string, len(logs))
	copy(tmp, logs)
	logs = []string{}
	return tmp
}

func Init(logdir string, logpref string, level string) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if level == "TRACE" {
		log_level = LV_TRACE
	} else if level == "DEBUG" {
		log_level = LV_DEBUG
	} else if level == "INFO" {
		log_level = LV_INFO
	} else if level == "WARN" {
		log_level = LV_WARN
	} else if level == "ERROR" {
		log_level = LV_ERROR
	} else if level == "FATAL" {
		log_level = LV_FATAL
	} else if level == "PANIC" {
		log_level = LV_PANIC
	} else {
		log_level = LV_INFO
	}

	if logdir != "" {
		err := os.MkdirAll(logdir, 0777)
		if err != nil {
			log.Fatalln("slog mkdir ", logdir, " err:", err)
		}
	}

	logfile := ""
	if logdir != "" && logpref != "" {
		logfile = logdir + "/" + logpref
	}

	if lg == nil {
		lg = &logger{logpref: logfile, logfp: nil, per: nil}
	} else {
		lg.resetOutput(logfile)
	}

}

func init() {
	pid := os.Getpid()
	headTrace = fmt.Sprintf("[TRACE] [%d]", pid)
	headDebug = fmt.Sprintf("[DEBUG] [%d]", pid)
	headInfo = fmt.Sprintf("[INFO] [%d]", pid)
	headWarn = fmt.Sprintf("[WARN] [%d]", pid)
	headError = fmt.Sprintf("[ERROR] [%d]", pid)
	headFatal = fmt.Sprintf("[FATAL] [%d]", pid)
	headPanic = fmt.Sprintf("[PANIC] [%d]", pid)

	headFmtTrace = fmt.Sprintf("%s ", headTrace)
	headFmtDebug = fmt.Sprintf("%s ", headDebug)
	headFmtInfo = fmt.Sprintf("%s ", headInfo)
	headFmtWarn = fmt.Sprintf("%s ", headWarn)
	headFmtError = fmt.Sprintf("%s ", headError)
	headFmtFatal = fmt.Sprintf("%s ", headFatal)
	headFmtPanic = fmt.Sprintf("%s ", headPanic)

	Init("", "", "TRACE")

	atomic.StoreInt64(&cnStamp, time.Now().Unix())
}

func Tracef(format string, v ...interface{}) {
	if LV_TRACE < log_level {
		return
	}

	slogMutex.Lock()
	defer slogMutex.Unlock()
	lg.Printf(headFmtTrace+format, v...)
	atomic.AddInt64(&cnTrace, 1)
}

func Traceln(v ...interface{}) {
	if LV_TRACE < log_level {
		return
	}

	slogMutex.Lock()
	defer slogMutex.Unlock()
	lg.Println(append([]interface{}{headTrace}, v...)...)
	atomic.AddInt64(&cnTrace, 1)
}

func Debugf(format string, v ...interface{}) {
	if LV_DEBUG < log_level {
		return
	}

	slogMutex.Lock()
	defer slogMutex.Unlock()
	lg.Printf(headFmtDebug+format, v...)
	atomic.AddInt64(&cnDebug, 1)
}

func Debugln(v ...interface{}) {
	if LV_DEBUG < log_level {
		return
	}

	slogMutex.Lock()
	defer slogMutex.Unlock()
	lg.Println(append([]interface{}{headDebug}, v...)...)
	atomic.AddInt64(&cnDebug, 1)
}

func Infof(format string, v ...interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if LV_INFO >= log_level {
		lg.Printf(headFmtInfo+format, v...)
		atomic.AddInt64(&cnInfo, 1)
	}
}

func Infoln(v ...interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if LV_INFO >= log_level {
		lg.Println(append([]interface{}{headInfo}, v...)...)
		atomic.AddInt64(&cnInfo, 1)
	}
}

func Warnf(format string, v ...interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if LV_WARN >= log_level {
		lg.Printf(headFmtWarn+format, v...)
		atomic.AddInt64(&cnWarn, 1)
	}
}

func Warnln(v ...interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if LV_WARN >= log_level {
		lg.Println(append([]interface{}{headWarn}, v...)...)
		atomic.AddInt64(&cnWarn, 1)
	}
}

func Errorf(format string, v ...interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if LV_ERROR >= log_level {
		lg.Printf(headFmtError+format, v...)
		atomic.AddInt64(&cnError, 1)
		addLogs(fmt.Sprintf(headFmtError+format, v...))
	}
}

func Errorln(v ...interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if LV_ERROR >= log_level {
		lg.Println(append([]interface{}{headError}, v...)...)
		atomic.AddInt64(&cnError, 1)
	}
}

func Fatalf(format string, v ...interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if LV_FATAL >= log_level {
		lg.Printf(headFmtFatal+format, v...)
		atomic.AddInt64(&cnFatal, 1)
		addLogs(fmt.Sprintf(headFmtError+format, v...))
	}
}

func Fatalln(v ...interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if LV_FATAL >= log_level {
		lg.Println(append([]interface{}{headFatal}, v...)...)
		atomic.AddInt64(&cnFatal, 1)
	}
}

func Panicf(format string, v ...interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if LV_PANIC >= log_level {
		lg.Panicf(headFmtPanic+format, v...)
		atomic.AddInt64(&cnPanic, 1)
		addLogs(fmt.Sprintf(headFmtError+format, v...))
	}
}

func Panicln(v ...interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	if LV_PANIC >= log_level {
		lg.Panicln(append([]interface{}{headPanic}, v...)...)
		atomic.AddInt64(&cnPanic, 1)
	}
}

func LogStat() (map[string]int64, []string) {

	st := map[string]int64{
		"TRACE": atomic.SwapInt64(&cnTrace, 0),
		"DEBUG": atomic.SwapInt64(&cnDebug, 0),
		"INFO":  atomic.SwapInt64(&cnInfo, 0),
		"WARN":  atomic.SwapInt64(&cnWarn, 0),
		"ERROR": atomic.SwapInt64(&cnError, 0),
		"FATAL": atomic.SwapInt64(&cnFatal, 0),
		"PANIC": atomic.SwapInt64(&cnPanic, 0),

		"STAMP": atomic.SwapInt64(&cnStamp, time.Now().Unix()),
	}

	return st, getLogs()

}
