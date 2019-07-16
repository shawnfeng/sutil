package mq

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestInstanceConfFromString(t *testing.T) {
	cases := []struct {
		s            string
		expectErr    bool
		expectedConf *instanceConf
	}{
		{
			fmt.Sprintf("default-0-%s-g1-0", defaultTestTopic),
			false,
			&instanceConf{
				group:     "default",
				role:      RoleTypeReader,
				topic:     defaultTestTopic,
				groupId:   "g1",
				partition: 0,
			},
		},
		{
			"test-1-topic-g1-1",
			false,
			&instanceConf{
				group:     "test",
				role:      RoleTypeWriter,
				topic:     "topic",
				groupId:   "g1",
				partition: 1,
			},
		},
		{
			"test-test-test-test-test",
			true,
			nil,
		},
		{
			"test-3-topic-test-0",
			true,
			nil,
		},
		{
			"test-1-topic-test-test",
			true,
			nil,
		},
	}

	for _, c := range cases {
		conf, err := instanceConfFromString(c.s)
		assert.Equal(t, c.expectErr, err != nil)
		assert.Equal(t, c.expectedConf, conf)
	}
}

func getSyncMapSizeUnSafe(m sync.Map) (ret int) {
	m.Range(func(k, v interface{}) bool {
		ret += 1
		return true
	})
	return
}

func TestInstanceManager_applyChangeEvent_apolloConfig(t *testing.T) {
	ctx := context.TODO()
	_ = SetConfiger(ctx, ConfigerTypeApollo)
	apolloConfig, _ := DefaultConfiger.(*ApolloConfig)

	t.Run("ignore add event", func(t *testing.T) {
		m := NewInstanceManager()

		conf := &instanceConf{
			group:     "default",
			role:      0,
			topic:     defaultTestTopic,
			groupId:   "g1",
			partition: 0,
		}

		ce := &center.ChangeEvent{
			Source:    center.Apollo,
			Namespace: center.DefaultApolloMQNamespace,
			Changes: map[string]*center.Change{
				m.buildKey(conf): {
					NewValue:   "localhost:9092",
					ChangeType: center.ADD,
				},
			},
		}

		assert.Equal(t, 0, getSyncMapSizeUnSafe(m.instances))
		m.applyChangeEvent(ctx, ce)
	})

	t.Run("modify/delete instance", func(t *testing.T) {
		m := NewInstanceManager()

		conf := &instanceConf{
			group:     "default",
			role:      0,
			topic:     defaultTestTopic,
			groupId:   "g1",
			partition: 0,
		}

		in, err := m.newInstance(ctx, conf)
		assert.Equal(t, nil, err)
		m.add(conf, in)

		ce := &center.ChangeEvent{
			Source:    center.Apollo,
			Namespace: center.DefaultApolloMQNamespace,
			Changes: map[string]*center.Change{
				apolloConfig.buildKey(ctx, defaultTestTopic, "brokers"): {
					ChangeType: center.MODIFY,
				},
			},
		}

		assert.Equal(t, 1, getSyncMapSizeUnSafe(m.instances))
		m.applyChangeEvent(ctx, ce)
		assert.Equal(t, 1, getSyncMapSizeUnSafe(m.instances))
		modifiedIn := m.get(ctx, conf)
		assert.NotEqual(t, modifiedIn, in)
	})

	t.Run("modify/delete default config", func(t *testing.T) {
		m := NewInstanceManager()

		conf := &instanceConf{
			group:     "unknown",
			role:      0,
			topic:     defaultTestTopic,
			groupId:   "g1",
			partition: 0,
		}

		in, err := m.newInstance(ctx, conf)
		assert.Equal(t, nil, err)
		m.add(conf, in)

		ce := &center.ChangeEvent{
			Source:    center.Apollo,
			Namespace: center.DefaultApolloMQNamespace,
			Changes: map[string]*center.Change{
				apolloConfig.buildKey(ctx, defaultTestTopic, "brokers"): {
					ChangeType: center.MODIFY,
				},
			},
		}

		assert.Equal(t, 1, getSyncMapSizeUnSafe(m.instances))
		m.applyChangeEvent(ctx, ce)
		assert.Equal(t, 1, getSyncMapSizeUnSafe(m.instances))
		modifiedIn := m.get(ctx, conf)
		assert.NotEqual(t, modifiedIn, in)
	})
}
