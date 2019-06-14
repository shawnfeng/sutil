// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/mq"
	"github.com/shawnfeng/sutil/slog"
	"github.com/shawnfeng/sutil/trace"
	"time"
)

type Msg struct {
	Id   int
	Body string
}

func main() {

	trace.InitDefaultTracer("mq.test")
	topic := "palfish.test.test"

	ctx := context.Background()
	ctx = context.WithValue(ctx, "Head", "hahahaha")
	span, ctx := opentracing.StartSpanFromContext(ctx, "main")
	if span != nil {
		defer span.Finish()
	}
	go func() {
		var msgs []mq.Message
		for i := 0; i < 10; i++ {
			value := &Msg{
				Id:   1,
				Body: fmt.Sprintf("%d", i),
			}

			msgs = append(msgs, mq.Message{
				Key:   value.Body,
				Value: value,
			})
			err := mq.WriteMsg(ctx, topic, value.Body, value)
			slog.Infof("in msg: %v, err:%v", value, err)
		}
		err := mq.WriteMsgs(ctx, topic, msgs...)
		slog.Infof("in msgs: %v, err:%v", msgs, err)
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			var msg Msg
			ctx1 := context.Background()
			ctx, err := mq.ReadMsgByGroup(ctx1, topic, "group2", &msg)
			slog.Infof("out msg: %v, ctx:%v, err:%v", msg, ctx, err)
		}
	}()

	defer mq.Close()

	time.Sleep(3 * time.Second)
}
