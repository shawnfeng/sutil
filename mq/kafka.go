// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"encoding/json"
	kafka "github.com/segmentio/kafka-go"
	//"github.com/shawnfeng/sutil/slog"
	"time"
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
		//MaxWait:        30 * time.Second,
	})

	return &KafkaReader{
		Reader: reader,
	}
}

func (m *KafkaReader) ReadMsg(ctx context.Context, v interface{}) error {
	msg, err := m.ReadMessage(ctx)
	if err != nil {
		return err
	}

	err = json.Unmarshal(msg.Value, v)
	if err != nil {
		return err
	}

	return nil
}

func (m *KafkaReader) FetchMsg(ctx context.Context, v interface{}) (Handler, error) {
	msg, err := m.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(msg.Value, v)
	if err != nil {
		return nil, err
	}

	return NewKafkaHandler(m.Reader, msg), nil
}

func (m *KafkaReader) Close() error {
	return m.Reader.Close()
}

type KafkaWriter struct {
	*kafka.Writer
}

func NewKafkaWriter(brokers []string, topic string) *KafkaWriter {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   brokers,
		Topic:     topic,
		Balancer:  &kafka.Hash{},
		BatchSize: 1,
		//RequiredAcks: 1,
		//Async:        true,
	})

	return &KafkaWriter{
		Writer: writer,
	}
}

func (m *KafkaWriter) WriteMsg(ctx context.Context, k string, v interface{}) error {
	msg, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = m.WriteMessages(ctx, kafka.Message{
		Key:   []byte(k),
		Value: msg,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *KafkaWriter) WriteMsgs(ctx context.Context, msgs ...Message) error {
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

	err := m.WriteMessages(ctx, kmsgs...)
	if err != nil {
		return err
	}

	return nil
}

func (m *KafkaWriter) Close() error {
	return m.Writer.Close()
}
