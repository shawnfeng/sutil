package mq

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWriter(t *testing.T) {
	wi, err := NewWriter(context.TODO(), "palfish.test.test")
	assert.Nil(t, err)
	w, ok := wi.(*KafkaWriter)
	assert.Equal(t, true, ok)
	assert.Equal(t, 100, w.config.BatchSize)
}
