package mq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKafkaWriter(t *testing.T) {
	w := NewKafkaWriter([]string{"prod.kafka1.ibanyu.com:9092", "prod.kafka2.ibanyu.com:9092", "prod.kafka3.ibanyu.com:9092"}, "palfish.test.test")
	assert.Equal(t, 100, w.config.BatchSize)
}
