// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	kafka "github.com/segmentio/kafka-go"
	"github.com/shawnfeng/sutil/slog"
)

const (
	defaultBatchSize = 1

	// REBALANCE_IN_PROGRESS
	ErrorMsgRebalanceInProgress = "Rebalance In Progress"
)

type KafkaHandler struct {
	msg    kafka.Message
	reader *kafka.Reader
}

func NewKafkaHandler(reader *kafka.Reader, msg kafka.Message) *KafkaHandler {
	return &KafkaHandler{
		msg:    msg,
		reader: reader,
	}
}

func (m *KafkaHandler) CommitMsg(ctx context.Context) error {
	return m.reader.CommitMessages(ctx, m.msg)
}

type KafkaReader struct {
	*kafka.Reader
}

func NewKafkaReader(brokers []string, topic, groupId string, partition, minBytes, maxBytes int, commitInterval time.Duration) *KafkaReader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupId,
		Partition:      partition,
		MinBytes:       minBytes,
		MaxBytes:       maxBytes,
		CommitInterval: commitInterval,
		StartOffset:    kafka.LastOffset,
		//MaxWait:        30 * time.Second,
		Logger:      slog.GetInfoLogger(),
		ErrorLogger: getErrorLogger(),
	})

	return &KafkaReader{
		Reader: reader,
	}
}

func (m *KafkaReader) logConfigToSpan(span opentracing.Span) {
	config := m.Config()
	span.LogFields(
		log.String(spanLogKeyMQType, fmt.Sprint(MQTypeKafka)),
		log.String(spanLogKeyKafkaBrokers, strings.Join(config.Brokers, apolloBrokersSep)),
		log.String(spanLogKeyKafkaGroupID, config.GroupID),
		log.Int(spanLogKeyKafkaPartition, config.Partition),
	)
}

func (m *KafkaReader) ReadMsg(ctx context.Context, v interface{}, ov interface{}) error {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		m.logConfigToSpan(span)
	}

	msg, err := m.ReadMessage(ctx)
	if err != nil {
		return err
	}

	err = json.Unmarshal(msg.Value, v)
	if err != nil {
		return err
	}

	err = json.Unmarshal(msg.Value, ov)
	if err != nil {
		return err
	}

	return nil
}

func (m *KafkaReader) FetchMsg(ctx context.Context, v interface{}, ov interface{}) (Handler, error) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		m.logConfigToSpan(span)
	}

	msg, err := m.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(msg.Value, v)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(msg.Value, ov)
	if err != nil {
		return nil, err
	}

	return NewKafkaHandler(m.Reader, msg), nil
}

func (m *KafkaReader) Close() error {
	return m.Reader.Close()
}

func (m *KafkaReader) SetOffsetAt(ctx context.Context, t time.Time) error {
	return m.Reader.SetOffsetAt(ctx, t)
}

func (m *KafkaReader) SetOffset(ctx context.Context, offset int64) error {
	return m.Reader.SetOffset(offset)
}

type KafkaWriter struct {
	*kafka.Writer
	// NOTE: KafkaWriter 没有 config 的 getter，故在此保留一份
	config kafka.WriterConfig
}

func NewKafkaWriter(brokers []string, topic string) *KafkaWriter {
	kafkaConfig := kafka.WriterConfig{
		Brokers:   brokers,
		Topic:     topic,
		Balancer:  &kafka.Hash{},
		BatchSize: defaultBatchSize,
		//RequiredAcks: 1,
		//Async:        true,
		Logger:      slog.GetInfoLogger(),
		ErrorLogger: getErrorLogger(),
	}
	// TODO should optimize this, too dumb, double get, reset batchsize
	config, _ := DefaultConfiger.GetConfig(context.TODO(), topic, MQTypeKafka)
	if config != nil {
		if config.BatchSize > defaultBatchSize {
			kafkaConfig.BatchSize = config.BatchSize
		}
		if config.BatchTimeout > 0 {
			kafkaConfig.BatchTimeout = config.BatchTimeout
		}
	}
	writer := kafka.NewWriter(kafkaConfig)

	return &KafkaWriter{
		Writer: writer,
		config: kafkaConfig,
	}
}

func (m *KafkaWriter) logConfigToSpan(span opentracing.Span) {
	config := m.config
	span.LogFields(
		log.String(spanLogKeyMQType, fmt.Sprint(MQTypeKafka)),
		log.String(spanLogKeyKafkaBrokers, strings.Join(config.Brokers, apolloBrokersSep)),
	)
}

func (m *KafkaWriter) WriteMsg(ctx context.Context, k string, v interface{}) error {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		m.logConfigToSpan(span)
	}

	msg, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return m.WriteMessages(ctx, kafka.Message{
		Key:   []byte(k),
		Value: msg,
	})
}

func (m *KafkaWriter) WriteMsgs(ctx context.Context, msgs ...Message) error {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		m.logConfigToSpan(span)
	}

	var kmsgs []kafka.Message
	for _, msg := range msgs {
		body, err := json.Marshal(msg.Value)
		if err != nil {
			return err
		}
		kmsgs = append(kmsgs, kafka.Message{
			Key:   []byte(msg.Key),
			Value: body,
		})
	}

	return m.WriteMessages(ctx, kmsgs...)
}

func (m *KafkaWriter) Close() error {
	return m.Writer.Close()
}

type errorLogger struct {
}

func getErrorLogger() *errorLogger {
	return &errorLogger{}
}

func (m *errorLogger) Printf(format string, items ...interface{}) {
	errMsg := fmt.Sprintf(format, items...)
	if strings.Contains(errMsg, ErrorMsgRebalanceInProgress) {
		slog.Warnf(errMsg)
		return
	}
	slog.Errorf(format, items...)
}
