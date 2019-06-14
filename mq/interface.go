// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/slog"
	// kafka "github.com/segmentio/kafka-go"
)

type Message struct {
	Key   string
	Value interface{}
}

func WriteMsg(ctx context.Context, topic string, key string, value interface{}) error {
	fun := "WriteMsg -->"

	span, _ := opentracing.StartSpanFromContext(ctx, "mq.WriteMsg")
	if span != nil {
		defer span.Finish()
	}

	//todo flag
	writer := DefaultInstanceManager.getWriter("", ROLE_TYPE_WRITER, topic, "", 0)
	if writer == nil {
		slog.Errorf(ctx, "%s getWriter err, topic: %s", fun, topic)
		return fmt.Errorf("%s, getWriter err, topic: %s", fun, topic)
	}

	return writer.WriteMsg(ctx, key, value)
}

func WriteMsgs(ctx context.Context, topic string, msgs ...Message) error {
	fun := "WriteMsgs -->"

	span, _ := opentracing.StartSpanFromContext(ctx, "mq.WriteMsgs")
	if span != nil {
		defer span.Finish()
	}

	//todo flag
	writer := DefaultInstanceManager.getWriter("", ROLE_TYPE_WRITER, topic, "", 0)
	if writer == nil {
		slog.Errorf(ctx, "%s getWriter err, topic: %s", fun, topic)
		return fmt.Errorf("%s, getWriter err, topic: %s", fun, topic)
	}

	return writer.WriteMsgs(ctx, msgs...)
}

// 读完消息后会自动提交offset
func ReadMsgByGroup(ctx context.Context, topic, groupId string, value interface{}) error {
	fun := "ReadMsgByGroup -->"

	span, _ := opentracing.StartSpanFromContext(ctx, "mq.ReadMsgByGroup")
	if span != nil {
		defer span.Finish()
	}

	//todo flag
	reader := DefaultInstanceManager.getReader("", ROLE_TYPE_READER, topic, groupId, 0)
	if reader == nil {
		slog.Errorf(ctx, "%s getReader err, topic: %s", fun, topic)
		return fmt.Errorf("%s, getReader err, topic: %s", fun, topic)
	}

	return reader.ReadMsg(ctx, value)
}

//
func ReadMsgByPartition(ctx context.Context, topic string, partition int, value interface{}) error {
	fun := "ReadMsgByPartition -->"

	span, _ := opentracing.StartSpanFromContext(ctx, "mq.ReadMsgByPartition")
	if span != nil {
		defer span.Finish()
	}

	//todo flag
	reader := DefaultInstanceManager.getReader("", ROLE_TYPE_READER, topic, "", partition)
	if reader == nil {
		slog.Errorf(ctx, "%s getReader err, topic: %s", fun, topic)
		return fmt.Errorf("%s, getReader err, topic: %s", fun, topic)
	}

	return reader.ReadMsg(ctx, value)
}

// 读完消息后不会自动提交offset,需要手动调用Handle.CommitMsg方法来提交offset
func FetchMsgByGroup(ctx context.Context, topic, groupId string, value interface{}) (Handler, error) {
	fun := "FetchMsgByGroup -->"

	span, _ := opentracing.StartSpanFromContext(ctx, "mq.FetchMsgByGroup")
	if span != nil {
		defer span.Finish()
	}

	//todo flag
	reader := DefaultInstanceManager.getReader("", ROLE_TYPE_READER, topic, groupId, 0)
	if reader == nil {
		slog.Errorf(ctx, "%s getReader err, topic: %s", fun, topic)
		return nil, fmt.Errorf("%s, getReader err, topic: %s", fun, topic)
	}

	return reader.FetchMsg(ctx, value)
}

func Close() {
	DefaultInstanceManager.Close()
}
