// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"fmt"
	// kafka "github.com/segmentio/kafka-go"
	//"github.com/shawnfeng/sutil/slog"
)

type Message struct {
	Key   string
	Value interface{}
}

type Writer interface {
	WriteMsg(ctx context.Context, key string, value interface{}) error
	WriteMsgs(ctx context.Context, msg ...Message) error
	Close() error
}

func NewWriter(topic string) (Writer, error) {
	config := DefaultConfiger.GetConfig(topic)

	mqType := config.MQType
	switch mqType {
	case MQ_TYPE_KAFKA:
		return NewKafkaWriter(config.MQAddr, topic), nil

	default:
		return nil, fmt.Errorf("mqType %s error", mqType)
	}
}
