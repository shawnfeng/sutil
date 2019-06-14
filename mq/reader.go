// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"fmt"
	//"time"
	// kafka "github.com/segmentio/kafka-go"
	//"github.com/shawnfeng/sutil/slog"
)

type Handler interface {
	CommitMsg(ctx context.Context) error
}

type Reader interface {
	FetchMsg(ctx context.Context, value interface{}) (Handler, error)
	ReadMsg(ctx context.Context, value interface{}) error
	Close() error
}

//CommitInterval indicates the interval at which offsets are committed to
// the broker.  If 0, commits will be handled synchronously.
func NewGroupReader(topic, groupId string) (Reader, error) {
	config := DefaultConfiger.GetConfig(topic)

	mqType := config.MQType
	switch mqType {
	case MQ_TYPE_KAFKA:
		return NewKafkaReader(config.MQAddr, topic, groupId, 0, 1, 10e6, config.CommitInterval), nil

	default:
		return nil, fmt.Errorf("mqType %d error", mqType)
	}
}

const (
	LastOffset  int64 = -1 // The most recent offset available for a partition.
	FirstOffset       = -2 // The least recent offset available for a partition.
)

func NewPartitionReader(topic string, partition int) (Reader, error) {
	config := DefaultConfiger.GetConfig(topic)

	offset := config.Offset
	mqType := config.MQType
	switch mqType {
	case MQ_TYPE_KAFKA:
		reader := NewKafkaReader(config.MQAddr, topic, "", partition, 1, 10e6, 0)
		err := reader.SetOffset(offset)
		if err != nil {
			return nil, err
		}

		return reader, err

	default:
		return nil, fmt.Errorf("mqType %d error", mqType)
	}
}
