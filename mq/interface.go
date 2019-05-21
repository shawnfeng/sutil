// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/slog"
	// kafka "github.com/segmentio/kafka-go"
)

type Message struct {
	Key   string
	Value interface{}
}

func WriteMsg(ctx context.Context, topic string, key string, value interface{}) error {
	fun := "WriteMsg -->"

	//todo flag
	writer := DefaultInstanceManager.getWriter("", ROLE_TYPE_WRITER, topic, "", 0)
	if writer == nil {
		slog.Errorf("%s getWriter err, topic: %s", fun, topic)
		return fmt.Errorf("%s, getWriter err, topic: %s", fun, topic)
	}

	return writer.WriteMsg(ctx, key, value)
}

func WriteMsgs(ctx context.Context, topic string, msgs ...Message) error {
	fun := "WriteMsgs -->"

	//todo flag
	writer := DefaultInstanceManager.getWriter("", ROLE_TYPE_WRITER, topic, "", 0)
	if writer == nil {
		slog.Errorf("%s getWriter err, topic: %s", fun, topic)
		return fmt.Errorf("%s, getWriter err, topic: %s", fun, topic)
	}

	return writer.WriteMsgs(ctx, msgs...)
}

func ReadMsgByGroup(ctx context.Context, topic, groupId string, value interface{}) error {
	fun := "ReadMsgByGroup -->"

	//todo flag
	reader := DefaultInstanceManager.getReader("", ROLE_TYPE_READER, topic, groupId, 0)
	if reader == nil {
		slog.Errorf("%s getReader err, topic: %s", fun, topic)
		return fmt.Errorf("%s, getReader err, topic: %s", fun, topic)
	}

	return reader.ReadMsg(ctx, value)
}

func ReadMsgByPartition(ctx context.Context, topic string, partition int, value interface{}) error {
	fun := "ReadMsgByPartition -->"

	//todo flag
	reader := DefaultInstanceManager.getReader("", ROLE_TYPE_READER, topic, "", partition)
	if reader == nil {
		slog.Errorf("%s getReader err, topic: %s", fun, topic)
		return fmt.Errorf("%s, getReader err, topic: %s", fun, topic)
	}

	return reader.ReadMsg(ctx, value)
}

func Close() {
	DefaultInstanceManager.Close()
}
