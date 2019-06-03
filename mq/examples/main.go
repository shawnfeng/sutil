// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/mq"
	"github.com/shawnfeng/sutil/slog"
	"time"
)

type Msg struct {
	Id   int
	Body string
}

func main() {
	topic := "palfish.test.test"

	ctx := context.Background()
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

	ctx1 := context.Background()
	go func() {
		for i := 0; i < 10000; i++ {
			var msg Msg
			handler, err := mq.FetchMsgByGroup(ctx1, topic, "group2", &msg)
			slog.Infof("out msg: %v, err:%v", msg, err)
			handler.CommitMsg(ctx1)
		}
	}()

	defer mq.Close()

	time.Sleep(3 * time.Second)
}
