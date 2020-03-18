// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"fmt"
)

type Writer interface {
	WriteMsg(ctx context.Context, key string, value interface{}) error
	WriteMsgs(ctx context.Context, msgs ...Message) error
	Close() error
}

func NewWriter(ctx context.Context, topic string) (Writer, error) {
	config, err := DefaultConfiger.GetConfig(ctx, topic, MQTypeKafka)
	if err != nil {
		return nil, err
	}

	mqType := config.MQType
	switch mqType {
	case MQTypeKafka:
		return NewKafkaWriter(config.MQAddr, wrapTopicFromContext(ctx, topic)), nil

	default:
		return nil, fmt.Errorf("mqType %d error", mqType)
	}
}
