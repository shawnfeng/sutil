package redisext

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewScript(t *testing.T) {
	script := `
	if redis.call("GET", KEYS[1]) ~= false then
			return redis.call("INCRBY", KEYS[1], ARGV[1])
	end
	return false`
	s := NewScript(script)
	assert.NotNil(t, s)
	assert.Equal(t, s.Hash(), "0b2cd31fec150908e9e0304c8189bc7168c0b441")
}

func TestEval(t *testing.T) {
	script := `
	if redis.call("GET", KEYS[1]) ~= false then
			return redis.call("INCRBY", KEYS[1], ARGV[1])
	end
	return false`
	s := NewScript(script)
	ctx := context.Background()
	m := NewRedisExt("test/test", "test")
	m.Set(ctx, "key1", 100, 1*time.Second)
	r, err := m.Eval(ctx, s, []string{"key1"}, 100)
	r1, _ := m.Get(ctx, "key1")
	assert.Nil(t, err)
	assert.Equal(t, r, int64(200))
	assert.Equal(t, r1, "200")
}
