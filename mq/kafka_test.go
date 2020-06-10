package mq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewKafkaWriter(t *testing.T) {
	w := NewKafkaWriter([]string{"prod.kafka1.ibanyu.com:9092", "prod.kafka2.ibanyu.com:9092", "prod.kafka3.ibanyu.com:9092"}, "palfish.test.test")
	assert.Equal(t, 100, w.config.BatchSize)
	assert.Equal(t, time.Duration(5000000), w.config.BatchTimeout)
}
