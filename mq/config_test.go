package mq

import (
	"context"
	"fmt"
	"github.com/kaneshin/go-pkg/testing/assert"
	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"os"
	"testing"
)

// NOTE: 跑测试时，需要配置 /etc/hosts
// 10.111.209.211 apollo-meta.ibanyu.com

const (
	defaultTestTopic = "palfish.test.test"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	_ = center.Init(context.TODO(), "test/test", []string{"infra.mq"})
}

func teardown() {
	_ = center.Stop(context.TODO())
}

func TestApolloConfig_GetConfig(t *testing.T) {
	ctx := context.TODO()

	conf := NewApolloConfiger()
	conf.Init(ctx)
	t.Run("valid topic", func(t *testing.T) {
		topic := defaultTestTopic
		config, err := conf.GetConfig(ctx, topic, MQTypeKafka)
		assert.Equal(t, err, nil)
		assert.Equal(t, config.MQType, MQTypeKafka)
		assert.Equal(t, config.Topic, topic)
		assert.True(t, len(config.MQAddr) > 0)
	})

	t.Run("invalid topic", func(t *testing.T) {
		topic := "topic.never.exist"
		config, err := conf.GetConfig(ctx, topic, MQTypeKafka)
		assert.True(t, config == nil)
		assert.NotEqual(t, err, nil)
	})

	t.Run("delay topic", func(t *testing.T) {
		topic := defaultTestTopic
		config, err := conf.GetConfig(ctx, topic, MQTypeDelay)
		assert.Equal(t, err, nil)
		assert.Equal(t, config.MQType, MQTypeDelay)
		assert.Equal(t, config.Topic, topic)
		assert.True(t, len(config.MQAddr) > 0)
	})
}

func TestApolloConfig_buildKey(t *testing.T) {
	ctx := context.TODO()
	conf := NewApolloConfiger()

	cases := []struct {
		topic          string
		item           string
		expectedString string
		mqType         MQType
	}{
		{defaultTestTopic, "brokers", fmt.Sprintf("%s.default.kafka.brokers", defaultTestTopic), MQTypeKafka},
		{"topic", "timeout", "topic.default.kafka.timeout", MQTypeKafka},
		{defaultTestTopic, "brokers", fmt.Sprintf("%s.default.delay.brokers", defaultTestTopic), MQTypeDelay},
	}

	for _, c := range cases {
		assert.Equal(t, conf.buildKey(ctx, c.topic, c.item, c.mqType), c.expectedString)
	}
}

func TestApolloConfig_ParseKey(t *testing.T) {
	ctx := context.TODO()
	conf := NewApolloConfiger()

	cases := []struct {
		key              string
		expectError      bool
		expectedKeyParts *KeyParts
	}{
		{
			"topic.default.kafka.brokers",
			false,
			&KeyParts{"topic", "default"},
		},
		{
			"palfish.test.test.default.kafka.brokers",
			false,
			&KeyParts{"palfish.test.test", "default"},
		},
		{
			"a.b.c.d",
			false,
			&KeyParts{"a", "b"},
		},
		{
			"key",
			true,
			nil,
		},
		{
			"base.changeboard.event.default.delay.brokers",
			false,
			&KeyParts{"base.changeboard.event", "default"},
		},
	}

	for _, c := range cases {
		keyParts, err := conf.ParseKey(ctx, c.key)
		assert.Equal(t, c.expectError, err != nil)
		assert.Equal(t, c.expectedKeyParts, keyParts)
	}
}

func TestApolloConfig_getConfigItemWithFallback(t *testing.T) {
	t.Run("empty ctx should get default value", func(t *testing.T) {
		ctx := context.TODO()
		conf := NewApolloConfiger()

		brokersVal, ok := conf.getConfigItemWithFallback(ctx, defaultTestTopic, apolloBrokersKey, MQTypeKafka)
		assert.True(t, ok)
		assert.True(t, len(brokersVal) > 0, "got brokers:", brokersVal)
		slog.Infof(ctx, "got brokers:%s", brokersVal)
	})

	t.Run("ctx with unknown group should get default value", func(t *testing.T) {
		ctx := context.TODO()
		ctx = context.WithValue(ctx, scontext.ContextKeyControl, simpleContextControlRouter{"unknown"})

		conf := NewApolloConfiger()

		brokersVal, ok := conf.getConfigItemWithFallback(ctx, defaultTestTopic, apolloBrokersKey, MQTypeKafka)
		assert.True(t, ok)
		assert.True(t, len(brokersVal) > 0, "got brokers:", brokersVal)
		slog.Infof(ctx, "got brokers:%s", brokersVal)
	})

	t.Run("ctx with known group should get its value", func(t *testing.T) {
		ctx := context.TODO()
		ctx = context.WithValue(ctx, scontext.ContextKeyControl, simpleContextControlRouter{"testgroup"})

		conf := NewApolloConfiger()

		brokersVal, ok := conf.getConfigItemWithFallback(ctx, defaultTestTopic, apolloBrokersKey, MQTypeKafka)
		assert.True(t, ok)
		assert.True(t, len(brokersVal) > 0, "got brokers:", brokersVal)
		slog.Infof(ctx, "got brokers:%s", brokersVal)
	})

	t.Run("empty ctx should get default delay value", func(t *testing.T) {
		ctx := context.TODO()
		conf := NewApolloConfiger()
		conf.Init(ctx)

		brokersVal, ok := conf.getConfigItemWithFallback(ctx, defaultTestTopic, apolloBrokersKey, MQTypeDelay)
		assert.True(t, ok)
		assert.True(t, len(brokersVal) > 0, "got brokers:", brokersVal)
		slog.Infof(ctx, "got brokers:%s", brokersVal)
	})
}
