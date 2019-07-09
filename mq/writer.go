// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"fmt"
	// kafka "github.com/segmentio/kafka-go"
	//"github.com/shawnfeng/sutil/slog/slog"
)

type Writer interface {
	WriteMsg(ctx context.Context, key string, value interface{}) error
	WriteMsgs(ctx context.Context, msgs ...Message) error
	Close() error
}

func NewWriter(topic string) (Writer, error) {
	config := DefaultConfiger.GetConfig(topic)

	mqType := config.MQType
	switch mqType {
	case MqTypeKafka:
		return NewKafkaWriter(config.MQAddr, topic), nil

	default:
		return nil, fmt.Errorf("mqType %d error", mqType)
	}
}
