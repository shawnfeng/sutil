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
	writer, err := mq.NewWriter("palfish.test.test")
	if err != nil {
		slog.Errorf("err: %s", err)
		return
	}
	defer writer.Close()

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
		}
		err := writer.WriteMsgs(ctx, msgs...)
		slog.Infof("in msgs: %v, err:%v", msgs, err)
	}()

	reader, err := mq.NewGroupReader("palfish.test.test", "testid", time.Second)
	if err != nil {
		slog.Errorf("err: %s", err)
		return
	}
	defer reader.Close()

	ctx1 := context.Background()
	go func() {
		for i := 0; i < 10000; i++ {
			var msg Msg
			err := reader.ReadMsg(ctx1, &msg)
			slog.Infof("out msg: %v, err:%v", msg, err)
		}
	}()

	time.Sleep(100 * time.Second)
}
