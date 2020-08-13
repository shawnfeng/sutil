// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"context"

	"gitlab.pri.ibanyu.com/middleware/seaweed/xlog"
)

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

// Sync sync of app logger
func Sync() {
	xlog.AppLogSync()
}

// Traceln xlog.Debug
func Traceln(v ...interface{}) {
	xlog.Debug(context.TODO(), v...)
}

// Tracef xlog.Debugf
func Tracef(format string, v ...interface{}) {
	xlog.Debugf(context.TODO(), format, v...)
}

// Debugln xlog.Debug
func Debugln(v ...interface{}) {
	xlog.Debug(context.TODO(), v...)
}

// Debugf xlog.Debugf
func Debugf(format string, v ...interface{}) {
	xlog.Debugf(context.TODO(), format, v...)
}

// Infof ...
func Infof(format string, v ...interface{}) {
	xlog.Infof(context.TODO(), format, v...)
}

// Infoln ...
func Infoln(v ...interface{}) {
	xlog.Info(context.TODO(), v...)
}

// Warnf ...
func Warnf(format string, v ...interface{}) {
	xlog.Warnf(context.TODO(), format, v...)
}

// Warnln ...
func Warnln(v ...interface{}) {
	xlog.Warn(context.TODO(), v...)
}

// Errorf ...
func Errorf(format string, v ...interface{}) {
	xlog.Errorf(context.TODO(), format, v...)
}

// Errorln ...
func Errorln(v ...interface{}) {
	xlog.Error(context.TODO(), v...)
}

// Fatalf ...
func Fatalf(format string, v ...interface{}) {
	xlog.Fatalf(context.TODO(), format, v...)
}

// Fatalln ...
func Fatalln(v ...interface{}) {
	xlog.Fatal(context.TODO(), v...)
}

// Panicf ...
func Panicf(format string, v ...interface{}) {
	xlog.Panicf(context.TODO(), format, v...)
}

// Panicln ...
func Panicln(v ...interface{}) {
	xlog.Panic(context.TODO(), v...)
}

// LogStat TODO: 统计日志打印情况，后期都可以去掉
func LogStat() (map[string]int64, []string) {
	return xlog.LogStat()
}

// Logger 注入其他基础库的日志句柄
type Logger struct {
}

func GetLogger() *Logger {
	return &Logger{}
}

func (m *Logger) Printf(format string, items ...interface{}) {
	Errorf(format, items...)
}

type InfoLogger struct {
}

func GetInfoLogger() *InfoLogger {
	return &InfoLogger{}
}

func (m *InfoLogger) Printf(format string, items ...interface{}) {
	Infof(format, items...)
}
