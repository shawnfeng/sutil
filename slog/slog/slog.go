// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"context"
	"fmt"
	"gitlab.pri.ibanyu.com/middleware/seaweed/xlog"
)

func formatFromContext(ctx context.Context, includeHead bool, format string) string {
	if cs := extractContextAsString(ctx, includeHead); cs != "" {
		return fmt.Sprintf("%s%s", cs, format)
	}

	return format
}

func vFromContext(ctx context.Context, includeHead bool, v ...interface{}) []interface{} {
	if cs := extractContextAsString(ctx, includeHead); len(cs) > 0 {
		return append([]interface{}{cs}, v...)
	}
	return v
}

func Tracef(ctx context.Context, format string, v ...interface{}) {
	//format = formatFromContext(ctx, false, format)
	xlog.Debugf(ctx, format, v...)
}

func Traceln(ctx context.Context, v ...interface{}) {
	//v = vFromContext(ctx, false, v...)
	xlog.Debug(ctx, v...)
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	//format = formatFromContext(ctx, false, format)
	xlog.Debugf(ctx, format, v...)
}

func Debugln(ctx context.Context, v ...interface{}) {
	//v = vFromContext(ctx, false, v...)
	xlog.Debug(ctx, v...)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	//format = formatFromContext(ctx, false, format)
	xlog.Infof(ctx, format, v...)
}

func Infoln(ctx context.Context, v ...interface{}) {
	//v = vFromContext(ctx, false, v...)
	xlog.Info(ctx, v...)
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	//format = formatFromContext(ctx, false, format)
	xlog.Warnf(ctx, format, v...)
}

func Warnln(ctx context.Context, v ...interface{}) {
	//v = vFromContext(ctx, false, v...)
	xlog.Warn(ctx, v...)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	//format = formatFromContext(ctx, true, format)
	xlog.Errorf(ctx, format, v...)
}

func Errorln(ctx context.Context, v ...interface{}) {
	//v = vFromContext(ctx, true, v...)
	xlog.Error(ctx, v...)
}

func Fatalf(ctx context.Context, format string, v ...interface{}) {
	//format = formatFromContext(ctx, true, format)
	xlog.Fatalf(ctx, format, v...)
}

func Fatalln(ctx context.Context, v ...interface{}) {
	//v = vFromContext(ctx, true, v...)
	xlog.Fatal(ctx, v...)
}

func Panicf(ctx context.Context, format string, v ...interface{}) {
	//format = formatFromContext(ctx, true, format)
	xlog.Panicf(ctx, format, v...)
}

func Panicln(ctx context.Context, v ...interface{}) {
	//v = vFromContext(ctx, true, v...)
	xlog.Panic(ctx, v...)
}

type Logger struct {}

func GetLogger() *Logger {
	return &Logger{}
}

func (m *Logger) Printf(format string, v ...interface{}) {
	Infof(context.Background(), format, v...)
}
