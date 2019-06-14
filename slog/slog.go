// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"context"
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
	lg        *zap.SugaredLogger

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

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000000"))
}

func CapitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}

func Sync() {
	lg.Sync()
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
		ljlogger = lumberjack.NewLogger(logfile, 10240000, 0, 0, true, false)

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
	logger := zap.New(core)
	lg = logger.Sugar()
}

func init() {
	Init("", "", "TRACE")

	atomic.StoreInt64(&cnStamp, time.Now().Unix())
}

func formatFromContext(ctx context.Context, includeHead bool, format string) string {
	if cs := extractContextAsString(ctx, includeHead); cs != "" {
		return fmt.Sprintf("%s %s", cs, format)
	}

	return format
}

func vFromContext(ctx context.Context, includeHead bool, v ...interface{}) []interface{} {
	if vv := extractContext(ctx, includeHead); len(vv) > 0 {
		return append(vv, append([]interface{}{" "}, v...)...)
	}
	return v
}

func Tracef(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	lg.Debugf(format, v...)
	atomic.AddInt64(&cnTrace, 1)
}

func Traceln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	lg.Debug(v...)
	atomic.AddInt64(&cnTrace, 1)
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	lg.Debugf(format, v...)
	atomic.AddInt64(&cnDebug, 1)
}

func Debugln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	lg.Debug(v...)
	atomic.AddInt64(&cnDebug, 1)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	lg.Infof(format, v...)
	atomic.AddInt64(&cnInfo, 1)
}

func Infoln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	lg.Info(v...)
	atomic.AddInt64(&cnInfo, 1)
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	lg.Warnf(format, v...)
	atomic.AddInt64(&cnWarn, 1)
}

func Warnln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	lg.Warn(v...)
	atomic.AddInt64(&cnWarn, 1)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format)
	lg.Errorf(format, v...)
	atomic.AddInt64(&cnError, 1)
	addLogs("ERROR " + fmt.Sprintf(format, v...))
}

func Errorln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, true, v...)
	lg.Error(v...)
	atomic.AddInt64(&cnError, 1)
	addLogs("ERROR " + fmt.Sprintln(v...))
}

func Fatalf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format)
	lg.Fatalf(format, v...)
	atomic.AddInt64(&cnFatal, 1)
	addLogs("FATAL " + fmt.Sprintf(format, v...))
}

func Fatalln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, true, v...)
	lg.Fatal(v...)
	atomic.AddInt64(&cnFatal, 1)
	addLogs("FATAL " + fmt.Sprintln(v...))
}

func Panicf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format)
	lg.Panicf(format, v...)
	atomic.AddInt64(&cnPanic, 1)
	addLogs("PANIC " + fmt.Sprintf(format, v...))
}

func Panicln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, true, v...)
	lg.Panic(v...)
	atomic.AddInt64(&cnPanic, 1)
	addLogs("PANIC " + fmt.Sprintln(v...))
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

type Logger struct {
}

func GetLogger() *Logger {
	return &Logger{}
}

func (m *Logger) Printf(ctx context.Context, format string, items ...interface{}) {
	Errorf(ctx, format, items...)
}
