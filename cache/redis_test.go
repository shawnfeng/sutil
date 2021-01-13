package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var client, _ = NewCommonRedis("test/test", 1024)

func TestTtl(t *testing.T) {
	key := "aaa"
	val := "bbb"
	exp := 2 * time.Hour
	setcmd := client.Set(key, val, exp)
	assert.True(t, setcmd.Err() == nil)
	ttl := client.TTL(key)
	assert.True(t, ttl.Val() > time.Hour && ttl.Val() <= exp)

	expire := client.Expire(key, time.Hour)
	assert.True(t, expire.Val())
	ttl = client.TTL(key)
	assert.True(t, ttl.Val() > 0 && ttl.Val() <= time.Hour)
}

func TestGet(t *testing.T) {
	key := "aaa"
	ttl := client.TTL(key)
	t.Log(ttl.Val().Seconds())
	get := client.Get(key)
	t.Log(get.Val())
}

func TestNewRedisByAddr(t *testing.T) {
	// 对比和同样参数的common是否一致
	key := "aaa"
	val := "bbb"
	addrClient, err := NewRedisByAddr("common.codis.pri.ibanyu.com:19000", "test/test", 1024)
	assert.NoError(t, err)
	addrClient.Set(key, val, 10*time.Second)

	strCmd := client.Get(key)
	str, err := strCmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, str, val)
}
