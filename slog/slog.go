// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	//"github.com/shawnfeng/sutil/stime"
	"github.com/shawnfeng/lumberjack.v2"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var logger *Logger

func init() {
	Init("", "", "TRACE")
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

type Logger struct {
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
	cnDrop  int64

	isClose int32

	lg *zap.SugaredLogger

	logs    []string
	logChan chan func()

	slogMutex sync.Mutex
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000000"))
}

func CapitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}

func Sync() {
	logger.stop()
	logger.lg.Sync()
}

func Init(logdir string, logpref string, level string) {
	log_level := zap.InfoLevel
	if level == "TRACE" {
		log_level = zap.DebugLevel
	} else if level == "DEBUG" {
		log_level = zap.DebugLevel
	} else if level == "INFO" {
		log_level = zap.InfoLevel
	} else if level == "WARN" {
		log_level = zap.WarnLevel
	} else if level == "ERROR" {
		log_level = zap.ErrorLevel
	} else if level == "FATAL" {
		log_level = zap.FatalLevel
	} else if level == "PANIC" {
		log_level = zap.PanicLevel
	} else {
		log_level = zap.InfoLevel
	}

	logfile := ""
	if logdir != "" && logpref != "" {
		logfile = logdir + "/" + logpref
	}

	var out io.Writer
	if len(logfile) > 0 {
		var ljlogger *lumberjack.Logger
		ljlogger = &lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    1024000,
			MaxBackups: 0,
			MaxAge:     0,
			LocalTime:  true,
		}

		go func() {
			for {
				now := time.Now().Unix()
				duration := 3600 - now%3600
				select {
				case <-time.After(time.Second * time.Duration(duration)):
					//st := stime.NewTimeStat()
					ljlogger.Rotate()
					//dur := st.Duration()
					//Infof("rotate tm:%d", dur)
				}
			}

		}()

		out = ljlogger
	} else {
		out = os.Stdout
	}
	w := zapcore.AddSync(out)

	enconf := zap.NewProductionEncoderConfig()
	enconf.EncodeTime = TimeEncoder
	enconf.CallerKey = "caller"
	enconf.EncodeCaller = zapcore.FullCallerEncoder
	enconf.EncodeLevel = CapitalLevelEncoder
	core := zapcore.NewCore(
		//zapcore.NewJSONEncoder(enconf),
		zapcore.NewConsoleEncoder(enconf),
		w,
		log_level,
	)
	lg := zap.New(core).Sugar()

	if logger != nil {
		logger.stop()
	}

	logger = &Logger{
		lg:      lg,
		cnTrace: 0,
		cnDebug: 0,
		cnInfo:  0,
		cnWarn:  0,
		cnError: 0,
		cnFatal: 0,
		cnPanic: 0,
		isClose: 0,
		cnStamp: time.Now().Unix(),
		logChan: make(chan func(), 1024*128),
	}

	go logger.run()
}

func (m *Logger) stop() {
	atomic.StoreInt32(&m.isClose, 1)
	close(m.logChan)
}

func (m *Logger) run() {
	fun := "Logger.run -->"

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case logFun, ok := <-m.logChan:
			if ok {
				if logFun != nil {
					logFun()
				}
			} else {
				fmt.Printf("%s logChan is close", fun)
				return
			}

		case <-ticker.C:
			cnDrop := atomic.SwapInt64(&m.cnDrop, 0)
			if cnDrop > 0 {
				errorf("%s drop log %d", fun, cnDrop)
			}
		}
	}
}

func (m *Logger) asyncDo(logFun func()) {

	isClose := atomic.LoadInt32(&m.isClose)
	if isClose == 1 {
		return
	}

	select {
	case m.logChan <- logFun:
	default:
		atomic.AddInt64(&m.cnDrop, 1)
	}
}

func (m *Logger) addLogs(log string) {
	m.slogMutex.Lock()
	defer m.slogMutex.Unlock()

	m.logs = append(m.logs, log)
	if len(m.logs) > 10 {
		m.logs = m.logs[len(m.logs)-10:]
	}
}

func (m *Logger) getLogs() []string {
	m.slogMutex.Lock()
	defer m.slogMutex.Unlock()

	tmp := make([]string, len(m.logs))
	copy(tmp, m.logs)
	m.logs = []string{}
	return tmp
}

func Tracef(format string, v ...interface{}) {

	logFun := func(format string, v ...interface{}) func() {
		return func() {
			logger.lg.Debugf(format, v...)
			atomic.AddInt64(&logger.cnTrace, 1)
		}
	}(format, v...)

	logger.asyncDo(logFun)
}

func Traceln(v ...interface{}) {

	logFun := func(v ...interface{}) func() {
		return func() {
			logger.lg.Debug(v...)
			atomic.AddInt64(&logger.cnTrace, 1)
		}
	}(v...)

	logger.asyncDo(logFun)
}

func Debugf(format string, v ...interface{}) {

	logFun := func(format string, v ...interface{}) func() {
		return func() {
			logger.lg.Debugf(format, v...)
			atomic.AddInt64(&logger.cnDebug, 1)
		}
	}(format, v...)

	logger.asyncDo(logFun)
}

func Debugln(v ...interface{}) {

	logFun := func(v ...interface{}) func() {
		return func() {
			logger.lg.Debug(v...)
			atomic.AddInt64(&logger.cnDebug, 1)
		}
	}(v...)

	logger.asyncDo(logFun)
}

func Infof(format string, v ...interface{}) {

	logFun := func(format string, v ...interface{}) func() {
		return func() {
			logger.lg.Infof(format, v...)
			atomic.AddInt64(&logger.cnInfo, 1)
		}
	}(format, v...)

	logger.asyncDo(logFun)
}

func Infoln(v ...interface{}) {

	logFun := func(v ...interface{}) func() {
		return func() {
			logger.lg.Info(v...)
			atomic.AddInt64(&logger.cnInfo, 1)
		}
	}(v...)

	logger.asyncDo(logFun)
}

func Warnf(format string, v ...interface{}) {

	logFun := func(format string, v ...interface{}) func() {
		return func() {
			logger.lg.Warnf(format, v...)
			atomic.AddInt64(&logger.cnWarn, 1)
		}
	}(format, v...)

	logger.asyncDo(logFun)
}

func Warnln(v ...interface{}) {

	logFun := func(v ...interface{}) func() {
		return func() {
			logger.lg.Warn(v...)
			atomic.AddInt64(&logger.cnWarn, 1)
		}
	}(v...)

	logger.asyncDo(logFun)
}

func errorf(format string, v ...interface{}) {
	logger.lg.Errorf(format, v...)
	atomic.AddInt64(&logger.cnError, 1)
	logger.addLogs("ERROR " + fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {

	logFun := func(format string, v ...interface{}) func() {
		return func() {
			errorf(format, v...)
		}
	}(format, v...)

	logger.asyncDo(logFun)
}

func Errorln(v ...interface{}) {

	logFun := func(v ...interface{}) func() {
		return func() {
			logger.lg.Error(v...)
			atomic.AddInt64(&logger.cnError, 1)
			logger.addLogs("ERROR " + fmt.Sprintln(v...))
		}
	}(v...)

	logger.asyncDo(logFun)
}

func Fatalf(format string, v ...interface{}) {

	logFun := func(format string, v ...interface{}) func() {
		return func() {
			logger.lg.Fatalf(format, v...)
			atomic.AddInt64(&logger.cnFatal, 1)
			logger.addLogs("FATAL " + fmt.Sprintf(format, v...))
		}
	}(format, v...)

	logger.asyncDo(logFun)
}

func Fatalln(v ...interface{}) {

	logFun := func(v ...interface{}) func() {
		return func() {
			logger.lg.Fatal(v...)
			atomic.AddInt64(&logger.cnFatal, 1)
			logger.addLogs("FATAL " + fmt.Sprintln(v...))
		}
	}(v...)

	logger.asyncDo(logFun)
}

func Panicf(format string, v ...interface{}) {

	logFun := func(format string, v ...interface{}) func() {
		return func() {
			logger.lg.Panicf(format, v...)
			atomic.AddInt64(&logger.cnPanic, 1)
			logger.addLogs("PANIC " + fmt.Sprintf(format, v...))
		}
	}(format, v...)

	logger.asyncDo(logFun)
}

func Panicln(v ...interface{}) {

	logFun := func(v ...interface{}) func() {
		return func() {
			logger.lg.Panic(v...)
			atomic.AddInt64(&logger.cnPanic, 1)
			logger.addLogs("PANIC " + fmt.Sprintln(v...))
		}
	}(v...)

	logger.asyncDo(logFun)
}

func LogStat() (map[string]int64, []string) {

	st := map[string]int64{
		"TRACE": atomic.SwapInt64(&logger.cnTrace, 0),
		"DEBUG": atomic.SwapInt64(&logger.cnDebug, 0),
		"INFO":  atomic.SwapInt64(&logger.cnInfo, 0),
		"WARN":  atomic.SwapInt64(&logger.cnWarn, 0),
		"ERROR": atomic.SwapInt64(&logger.cnError, 0),
		"FATAL": atomic.SwapInt64(&logger.cnFatal, 0),
		"PANIC": atomic.SwapInt64(&logger.cnPanic, 0),

		"STAMP": atomic.SwapInt64(&logger.cnStamp, time.Now().Unix()),
	}

	return st, logger.getLogs()

}

func GetLogger() *Logger {
	return &Logger{}
}

func (m *Logger) Printf(format string, items ...interface{}) {
	Errorf(format, items...)
}
