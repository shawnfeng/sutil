// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"gitlab.pri.ibanyu.com/middleware/seaweed/xlog"
)

var (
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

	logs []string
)

func addLogs(log string) {
	slogMutex.Lock()
	defer slogMutex.Unlock()
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

func convertLevel(level string) xlog.Level {
	switch level {
	case "info":
		return xlog.InfoLevel
	case "warn":
		return xlog.WarnLevel
	case "error":
		return xlog.ErrorLevel
	case "fatal":
		return xlog.FatalLevel
	case "panic":
		return xlog.PanicLevel
	default:
		return xlog.InfoLevel
	}
}

// Init init applog in xlog
func Init(logdir string, fileName string, level string) {
	xlog.InitAppLog(logdir, fileName, convertLevel(level))
}

func init() {
	atomic.StoreInt64(&cnStamp, time.Now().Unix())
}

// Traceln xlog.Debug
func Traceln(v ...interface{}) {
	xlog.Debug(context.TODO(), v...)
	atomic.AddInt64(&cnTrace, 1)
}

// Tracef xlog.Debugf
func Tracef(format string, v ...interface{}) {
	xlog.Debugf(context.TODO(), format, v...)
	atomic.AddInt64(&cnTrace, 1)
}

// Debugln xlog.Debug
func Debugln(v ...interface{}) {
	xlog.Debug(context.TODO(), v...)
	atomic.AddInt64(&cnDebug, 1)
}

// Debugf xlog.Debugf
func Debugf(format string, v ...interface{}) {
	xlog.Debugf(context.TODO(), format, v...)
	atomic.AddInt64(&cnDebug, 1)
}

// Infof ...
func Infof(format string, v ...interface{}) {
	xlog.Infof(context.TODO(), format, v...)
	atomic.AddInt64(&cnInfo, 1)
}

// Infoln ...
func Infoln(v ...interface{}) {
	xlog.Info(context.TODO(), v...)
	atomic.AddInt64(&cnInfo, 1)
}

// Warnf ...
func Warnf(format string, v ...interface{}) {
	xlog.Warnf(context.TODO(), format, v...)
	atomic.AddInt64(&cnWarn, 1)
}

// Warnln ...
func Warnln(v ...interface{}) {
	xlog.Warn(context.TODO(), v...)
	atomic.AddInt64(&cnWarn, 1)
}

// Errorf ...
func Errorf(format string, v ...interface{}) {
	xlog.Errorf(context.TODO(), format, v...)
	atomic.AddInt64(&cnError, 1)
	addLogs("ERROR " + fmt.Sprintf(format, v...))
}

// Errorln ...
func Errorln(v ...interface{}) {
	xlog.Error(context.TODO(), v...)
	atomic.AddInt64(&cnError, 1)
	addLogs("ERROR " + fmt.Sprintln(v...))
}

// Fatalf ...
func Fatalf(format string, v ...interface{}) {
	xlog.Fatalf(context.TODO(), format, v...)
	atomic.AddInt64(&cnFatal, 1)
	addLogs("FATAL " + fmt.Sprintf(format, v...))
}

// Fatalln ...
func Fatalln(v ...interface{}) {
	xlog.Fatal(context.TODO(), v...)
	atomic.AddInt64(&cnFatal, 1)
	addLogs("FATAL " + fmt.Sprintln(v...))
}

// Panicf ...
func Panicf(format string, v ...interface{}) {
	xlog.Panicf(context.TODO(), format, v...)
	atomic.AddInt64(&cnPanic, 1)
	addLogs("PANIC " + fmt.Sprintf(format, v...))
}

// Panicln ...
func Panicln(v ...interface{}) {
	xlog.Panic(context.TODO(), v...)
	atomic.AddInt64(&cnPanic, 1)
	addLogs("PANIC " + fmt.Sprintln(v...))
}

// LogStat TODO: 统计日志打印情况，后期都可以去掉
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

type Logger struct {
}

func GetLogger() *Logger {
	return &Logger{}
}

func (m *Logger) Printf(format string, items ...interface{}) {
	Errorf(format, items...)
}
