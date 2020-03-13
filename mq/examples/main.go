// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	//"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/mq"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"github.com/shawnfeng/sutil/trace"
	"time"
)

type Msg struct {
	Id   int
	Body string
}

func main() {

	_ = trace.InitDefaultTracer("mq.test")
	topic := "palfish.test.test"

	ctx := context.Background()
	ctx = context.WithValue(ctx, scontext.ContextKeyHead, "hahahaha")
	ctx = context.WithValue(ctx, scontext.ContextKeyControl, "{\"group\": \"g1\"}")
	span, ctx := opentracing.StartSpanFromContext(ctx, "main")
	if span != nil {
		defer span.Finish()
	}

	//_ = mq.SetConfiger(ctx, mq.ConfigerTypeApollo)
	//mq.WatchUpdate(ctx)

	/*
		go func() {
			var msgs []mq.Message
			for i := 0; i < 3; i++ {
				value := &Msg{
					Id:   1,
					Body: fmt.Sprintf("%d", i),
				}

				msgs = append(msgs, mq.Message{
					Key:   value.Body,
					Value: value,
				})
				err := mq.WriteMsg(ctx, topic, value.Body, value)
				slog.Infof(ctx, "in msg: %v, err:%v", value, err)
			}
			err := mq.WriteMsgs(ctx, topic, msgs...)
			slog.Infof(ctx, "in msgs: %v, err:%v", msgs, err)
		}()
	*/

	go func() {
		msg := &Msg{
			Id:   1,
			Body: "test",
		}
		jobID, err := mq.WriteDelayMsg(ctx, topic, msg, 5)
		slog.Infof(ctx, "write delay msg, jobID = %s, err = %v", jobID, err)
	}()

	ctx1 := context.Background()
	/*
		//err := mq.SetOffsetAt(ctx1, topic, 1, time.Date(2019, time.December, 4, 0, 0, 0, 0, time.UTC))
		//err := mq.SetOffset(ctx1, topic, 1, -2)
		if err != nil {
			slog.Infof(ctx, "2222222222222222,err:%v", err)
		}
	*/
	//go func() {
	//	for i := 0; i < 10000; i++ {
	//		var msg Msg
	//		ctx, err := mq.ReadMsgByGroup(ctx1, topic, "group3", &msg)
	//		slog.Infof(ctx, "1111111111111111out msg: %v, ctx:%v, err:%v", msg, ctx, err)
	//	}
	//}()

	go func() {
		for i := 0; i < 10; i ++ {
			var msg Msg
			ctx, ack, err := mq.FetchDelayMsg(ctx1, topic, &msg)
			slog.Infof(ctx, "1111111111111111out msg: %v, ctx:%v, err:%v", msg, ctx, err)
			err = ack.Ack(ctx)
			slog.Infof(ctx, "2222222222222222out msg: %v, ctx:%v, err:%v", msg, ctx, err)
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(15 * time.Second)
	defer mq.Close()

}
