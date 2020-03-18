// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"fmt"
	"time"
)

type Handler interface {
	CommitMsg(ctx context.Context) error
}

type Reader interface {
	FetchMsg(ctx context.Context, value interface{}, ovalue interface{}) (Handler, error)
	ReadMsg(ctx context.Context, value interface{}, ovalue interface{}) error
	SetOffsetAt(ctx context.Context, t time.Time) error
	SetOffset(ctx context.Context, offset int64) error
	Close() error
}

//CommitInterval indicates the interval at which offsets are committed to
// the broker.  If 0, commits will be handled synchronously.
func NewGroupReader(ctx context.Context, topic, groupId string) (Reader, error) {
	config, err := DefaultConfiger.GetConfig(ctx, topic, MQTypeKafka)
	if err != nil {
		return nil, err
	}

	mqType := config.MQType
	switch mqType {
	case MQTypeKafka:
		return NewKafkaReader(config.MQAddr, wrapTopicFromContext(ctx, topic), groupId, 0, 1, 10e6, config.CommitInterval), nil

	default:
		return nil, fmt.Errorf("mqType %d error", mqType)
	}
}

const (
	LastOffset  int64 = -1 // The most recent offset available for a partition.
	FirstOffset       = -2 // The least recent offset available for a partition.
)

func NewPartitionReader(ctx context.Context, topic string, partition int) (Reader, error) {
	config, err := DefaultConfiger.GetConfig(ctx, topic, MQTypeKafka)
	if err != nil {
		return nil, err
	}

	offsetAt := config.OffsetAt
	mqType := config.MQType
	switch mqType {
	case MQTypeKafka:
		reader := NewKafkaReader(config.MQAddr, wrapTopicFromContext(ctx, topic), "", partition, 1, 10e6, 0)
		if len(offsetAt) == 0 {
			return nil, fmt.Errorf("no offsetAt config found")
		}

		t, err := time.Parse("2006-01-02T15:04:05", offsetAt)
		if err != nil {
			return nil, err
		}

		err = reader.SetOffsetAt(ctx, t)
		if err != nil {
			return nil, err
		}

		return reader, err

	default:
		return nil, fmt.Errorf("mqType %d error", mqType)
	}
}
