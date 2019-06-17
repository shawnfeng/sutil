// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/shawnfeng/sutil/trace"
	"testing"
)

type testHead struct {
	uid     int64
	source  int32
	ip      string
	region  string
	dt      int32
	unionid string
}

func (th *testHead) toKV() map[string]interface{} {
	return map[string]interface{}{
		"uid":     th.uid,
		"source":  th.source,
		"ip":      th.ip,
		"region":  th.region,
		"dt":      th.dt,
		"unionid": th.unionid,
	}
}

var ctx, lctx context.Context

func init() {
	_ = trace.InitDefaultTracer("slog test")
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("testlog")
	ctx = context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = context.WithValue(ctx, "Head", &testHead{
		uid:     1234,
		source:  5678,
		ip:      "192.168.0.1",
		region:  "asia",
		dt:      1560499340,
		unionid: "7494ab07987ba112bd5c4f9857ccfb3f",
	})

	// ctx with large uid
	lspan := tracer.StartSpan("ltestlog")
	lctx = context.Background()
	lctx = opentracing.ContextWithSpan(lctx, lspan)
	lctx = context.WithValue(lctx, "Head", &testHead{
		uid:     1234567890,
		source:  5678,
		ip:      "192.168.0.1",
		region:  "asia",
		dt:      1560499340,
		unionid: "7494ab07987ba112bd5c4f9857ccfb3f",
	})
}

func TestLog(t *testing.T) {
	startTestSuite("To Console", t)
	fmt.Println()
}

func startTestSuite(name string, t *testing.T) {
	printBar(name)
	lnCases := []struct {
		ctx context.Context
		v   []interface{}
	}{
		{context.TODO(), []interface{}{"user not found"}},
		{context.TODO(), []interface{}{"user not found", 1, false}},
		{ctx, []interface{}{"user not found"}},
		{ctx, []interface{}{"user not found", 1, false}},
	}

	fCases := []struct {
		ctx    context.Context
		format string
		v      []interface{}
	}{
		{context.TODO(), "dbcluster not found", []interface{}{}},
		{context.TODO(), "%s err:%v", []interface{}{"CheckAuth-->", errors.New("key not found")}},
		{ctx, "not param", []interface{}{}},
		{ctx, "%s err:%v", []interface{}{"CheckAuth-->", errors.New("key not found")}},
	}

	t.Run("Traceln", func(t *testing.T) {
		printBar("Traceln")
		for _, c := range lnCases {
			Traceln(c.ctx, c.v...)
		}
	})

	t.Run("Tracef", func(t *testing.T) {
		printBar("Tracef")
		for _, c := range fCases {
			Tracef(c.ctx, c.format, c.v...)
		}
	})

	t.Run("Debugln", func(t *testing.T) {
		printBar("Debugln")
		for _, c := range lnCases {
			Debugln(c.ctx, c.v...)
		}
	})

	t.Run("Debugf", func(t *testing.T) {
		printBar("Debugf")
		for _, c := range fCases {
			Debugf(c.ctx, c.format, c.v...)
		}
	})

	t.Run("Infoln", func(t *testing.T) {
		printBar("Infoln")
		for _, c := range lnCases {
			Infoln(c.ctx, c.v...)
		}
	})

	t.Run("Infof", func(t *testing.T) {
		printBar("Infof")
		for _, c := range fCases {
			Infof(c.ctx, c.format, c.v...)
		}
	})

	t.Run("Warnln", func(t *testing.T) {
		printBar("Warnln")
		for _, c := range lnCases {
			Warnln(c.ctx, c.v...)
		}
	})

	t.Run("Warnf", func(t *testing.T) {
		printBar("Warnf")
		for _, c := range fCases {
			Warnf(c.ctx, c.format, c.v...)
		}
	})

	t.Run("Errorln", func(t *testing.T) {
		printBar("Errorln")
		for _, c := range lnCases {
			Errorln(c.ctx, c.v...)
		}
	})

	t.Run("Errorf", func(t *testing.T) {
		printBar("Errorf")
		for _, c := range fCases {
			Errorf(c.ctx, c.format, c.v...)
		}
	})

	t.Run("test large uid", func(t *testing.T) {
		printBar("LargeUid")
		for _, c := range lnCases {
			Errorln(lctx, c.v...)
		}

		for _, c := range fCases {
			Errorf(lctx, c.format, c.v...)
		}
	})

	//t.Run("Fatalln", func(t *testing.T) {
	//	printBar("Fatalln")
	//	for _, c := range lnCases {
	//		Fatalln(c.ctx, c.v...)
	//	}
	//})
	//
	//t.Run("Fatalf", func(t *testing.T) {
	//	printBar("Fatalf")
	//	for _, c := range fCases {
	//		Fatalf(c.ctx, c.format, c.v...)
	//	}
	//})

	//t.Run("Panicln", func(t *testing.T) {
	//	printBar("Panicln")
	//	for _, c := range lnCases {
	//		Panicln(c.ctx, c.v...)
	//	}
	//})
	//
	//t.Run("Panicf", func(t *testing.T) {
	//	printBar("Panicf")
	//	for _, c := range fCases {
	//		Panicf(c.ctx, c.format, c.v...)
	//	}
	//})
}

func printBar(title string) {
	fmt.Printf("================= %s =================\n", title)
}
