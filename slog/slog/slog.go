// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/slog"
)

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
	slog.Tracef(format, v...)
}

func Traceln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	slog.Traceln(v...)
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	slog.Debugf(format, v...)
}

func Debugln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	slog.Debugln(v...)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	slog.Infof(format, v...)
}

func Infoln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	slog.Infoln(v...)
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	slog.Warnf(format, v...)
}

func Warnln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	slog.Warnln(v...)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format)
	slog.Errorf(format, v...)
}

func Errorln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, true, v...)
	slog.Errorln(v...)
}

func Fatalf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format)
	slog.Fatalf(format, v...)
}

func Fatalln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, true, v...)
	slog.Fatalln(v...)
}

func Panicf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format)
	slog.Panicf(format, v...)
}

func Panicln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, true, v...)
	slog.Panicln(v...)
}
