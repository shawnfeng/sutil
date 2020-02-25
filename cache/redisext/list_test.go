package redisext

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	listName = "unittest_list"
)

func TestRedisExt_LIndex(t *testing.T) {
	ctx := context.Background()
	client := NewRedisExt("base/report", "test")
	client.Del(ctx, listName)
	n, err := client.LPush(ctx, listName, "one", "two")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), n)
	r, err := client.LIndex(ctx, listName, 1)
	assert.NoError(t, err)
	assert.Equal(t, "one", r)
	client.Del(ctx, listName)
}

func TestRedisExt_LInsert(t *testing.T) {
	ctx := context.Background()
	client := NewRedisExt("base/report", "test")
	client.Del(ctx, listName)
	_, err := client.RPush(ctx, listName, "Hello", "World")
	assert.NoError(t, err)
	_, err = client.LInsert(ctx, listName, "BEFORE", "World", "Three")
	assert.NoError(t, err)
	r, err := client.LRange(ctx, listName, 0, -1)
	assert.Equal(t, []string{"Hello", "Three", "World"}, r)
	client.Del(ctx, listName)
}

func TestRedisExt_LLen(t *testing.T) {
	ctx := context.Background()
	client := NewRedisExt("base/report", "test")
	client.Del(ctx, listName)
	_, err := client.LPush(ctx, listName, "one", "two")
	assert.NoError(t, err)
	n, err := client.LLen(ctx, listName)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), n)
	client.Del(ctx, listName)
}

func TestRedisExt_LPop(t *testing.T) {
	ctx := context.Background()
	client := NewRedisExt("base/report", "test")
	client.Del(ctx, listName)
	_, err := client.RPush(ctx, listName, "one", "two")
	assert.NoError(t, err)
	r, err := client.LPop(ctx, listName)
	assert.NoError(t, err)
	assert.Equal(t, "one", r)
	client.Del(ctx, listName)
}

func TestRedisExt_LPush(t *testing.T) {
	ctx := context.Background()
	client := NewRedisExt("base/report", "test")
	client.Del(ctx, listName)
	_, err := client.LPush(ctx, listName, "one", "two")
	assert.NoError(t, err)
	r, err := client.LRange(ctx, listName, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"two", "one"}, r)
	client.Del(ctx, listName)
}

func TestRedisExt_LRem(t *testing.T) {
	ctx := context.Background()
	client := NewRedisExt("base/report", "test")
	client.Del(ctx, listName)
	_, err := client.RPush(ctx, listName, "one", "two", "one", "three", "one")
	assert.NoError(t, err)
	n, err := client.LRem(ctx, listName, 2, "one")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), n)
	r, err := client.LRange(ctx, listName, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"two", "three", "one"}, r)
	client.Del(ctx, listName)
}

func TestRedisExt_LSet(t *testing.T) {
	ctx := context.Background()
	client := NewRedisExt("base/report", "test")
	client.Del(ctx, listName)
	_, err := client.RPush(ctx, listName, "one", "two", "three")
	assert.NoError(t, err)
	r1, err := client.LSet(ctx, listName, 0, "four")
	assert.NoError(t, err)
	assert.Equal(t, "OK", r1)
	r2, err := client.LRange(ctx, listName, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"four", "two", "three"}, r2)
	client.Del(ctx, listName)
}

func TestRedisExt_LTrim(t *testing.T) {
	ctx := context.Background()
	client := NewRedisExt("base/report", "test")
	client.Del(ctx, listName)
	_, err := client.RPush(ctx, listName, "one", "two", "three")
	assert.NoError(t, err)
	r1, err := client.LTrim(ctx, listName, 1, -1)
	assert.NoError(t, err)
	assert.Equal(t, "OK", r1)
	r2, err := client.LRange(ctx, listName, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"two", "three"}, r2)
	client.Del(ctx, listName)
}

func TestRedisExt_RPop(t *testing.T) {
	ctx := context.Background()
	client := NewRedisExt("base/report", "test")
	client.Del(ctx, listName)
	_, err := client.RPush(ctx, listName, "one", "two", "three")
	assert.NoError(t, err)
	r1, err := client.RPop(ctx, listName)
	assert.NoError(t, err)
	assert.Equal(t, "three", r1)
	r2, err := client.LRange(ctx, listName, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"one", "two"}, r2)
	client.Del(ctx, listName)
}
