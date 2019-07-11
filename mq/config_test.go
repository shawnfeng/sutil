package mq

import (
	"context"
	"fmt"
	"github.com/kaneshin/go-pkg/testing/assert"
	"github.com/shawnfeng/sutil/sconf/center"
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

	conf := NewApolloConfig()

	t.Run("valid topic", func(t *testing.T) {
		topic := defaultTestTopic
		config, err := conf.GetConfig(ctx, topic)
		assert.Equal(t, err, nil)
		assert.Equal(t, config.MQType, MqTypeKafka)
		assert.Equal(t, config.Topic, topic)
		assert.True(t, len(config.MQAddr) > 0)
	})

	t.Run("invalid topic", func(t *testing.T) {
		topic := "topic.never.exist"
		config, err := conf.GetConfig(ctx, topic)
		assert.True(t, config == nil)
		assert.NotEqual(t, err, nil)
	})
}

func TestApolloConfig_buildKey(t *testing.T) {
	ctx := context.TODO()
	conf := NewApolloConfig()

	cases := []struct {
		topic          string
		item           string
		expectedString string
	}{
		{defaultTestTopic, "brokers", fmt.Sprintf("%s.default.kafka.brokers", defaultTestTopic)},
		{"topic", "timeout", "topic.default.kafka.timeout"},
	}

	for _, c := range cases {
		assert.Equal(t, conf.buildKey(ctx, c.topic, c.item), c.expectedString)
	}
}

func TestApolloConfig_ParseKey(t *testing.T) {
	ctx := context.TODO()
	conf := NewApolloConfig()

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
	}

	for _, c := range cases {
		keyParts, err := conf.ParseKey(ctx, c.key)
		assert.Equal(t, c.expectError, err != nil)
		assert.Equal(t, c.expectedKeyParts, keyParts)
	}
}
