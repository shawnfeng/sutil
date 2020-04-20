package redisext

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	zsetTestKey = "myzset"
)

func TestRedisExt_Set_(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("test/test", "test")
	k, v := "hello", "world"
	_, err := re.Set(ctx, k, v, 10*time.Second)
	assert.NoError(t, err)
}

func TestRedisExt_ZAdd(t *testing.T) {
	ctx := context.Background()

	re := NewRedisExt("base/report", "test")

	members := []Z{
		{1, "one"},
		{2, "two"},
		{3, "three"},
	}
	n, err := re.ZAdd(ctx, zsetTestKey, members)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(members)), n)

	n, err = re.ZAddNX(ctx, zsetTestKey, members)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)

	dn, err := re.Del(ctx, zsetTestKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dn)
}

func TestRedisExt_ZRange(t *testing.T) {
	ctx := context.Background()

	re := NewRedisExt("base/report", "test")

	// prepare
	members := []Z{
		{1, "one"},
		{2, "two"},
		{3, "three"},
	}

	n, err := re.ZAdd(ctx, zsetTestKey, members)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(members)), n)

	// tests
	ss, err := re.ZRange(ctx, zsetTestKey, 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, []string{"one", "two", "three"}, ss)

	rss, err := re.ZRevRange(ctx, zsetTestKey, 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, []string{"three", "two", "one"}, rss)

	zs, err := re.ZRangeWithScores(ctx, zsetTestKey, 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, members, zs)

	rzs, err := re.ZRevRangeWithScores(ctx, zsetTestKey, 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, []Z{
		{3, "three"},
		{2, "two"},
		{1, "one"},
	}, rzs)

	// cleanup
	dn, err := re.Del(ctx, zsetTestKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dn)
}

func TestRedisExt_ZRank(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("base/report", "test")

	// prepare
	members := []Z{
		{1, "one"},
		{2, "two"},
		{3, "three"},
	}

	n, err := re.ZAdd(ctx, zsetTestKey, members)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(members)), n)

	// tests
	n, err = re.ZRank(ctx, zsetTestKey, "one")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)

	n, err = re.ZRevRank(ctx, zsetTestKey, "one")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), n)

	n, err = re.ZRank(ctx, zsetTestKey, "four")
	assert.Error(t, err)
	assert.Equal(t, int64(0), n)

	n, err = re.ZRevRank(ctx, zsetTestKey, "four")
	assert.Error(t, err)
	assert.Equal(t, int64(0), n)

	// cleanup
	dn, err := re.Del(ctx, zsetTestKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dn)
}

func TestRedisExt_ZCount(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("base/report", "test")

	// prepare
	members := []Z{
		{1, "one"},
		{2, "two"},
		{3, "three"},
	}

	n, err := re.ZAdd(ctx, zsetTestKey, members)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(members)), n)

	// tests
	n, err = re.ZCount(ctx, zsetTestKey, "2", "3")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), n)

	// cleanup
	dn, err := re.Del(ctx, zsetTestKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dn)
}

func TestRedisExt_ZScore(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("base/report", "test")

	// prepare
	members := []Z{
		{1, "one"},
		{2, "two"},
		{3, "three"},
	}

	n, err := re.ZAdd(ctx, zsetTestKey, members)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(members)), n)

	// tests
	f, err := re.ZScore(ctx, zsetTestKey, "two")
	assert.NoError(t, err)
	assert.Equal(t, float64(2), f)

	f, err = re.ZScore(ctx, zsetTestKey, "one")
	assert.NoError(t, err)
	assert.Equal(t, float64(1), f)

	f, err = re.ZScore(ctx, zsetTestKey, "four")
	assert.Error(t, err)
	assert.Equal(t, float64(0), f)

	// cleanup
	dn, err := re.Del(ctx, zsetTestKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dn)
}

func TestRedisExt_Expire(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("base/report", "test")

	// prepare
	members := []Z{
		{1, "one"},
		{2, "two"},
		{3, "three"},
	}

	n, err := re.ZAdd(ctx, zsetTestKey, members)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(members)), n)

	expiration := 1 * time.Second
	b, err := re.Expire(ctx, zsetTestKey, expiration)
	assert.NoError(t, err)
	assert.True(t, b)

	time.Sleep(expiration * 2)

	n, err = re.Exists(ctx, zsetTestKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)
}

func TestRedisExt_SetBit(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("base/report", "test")

	n, err := re.SetBit(ctx, "bitoptest", 2, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)
}

func TestRedisExt_GetBit(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("base/report", "test")

	n, err := re.GetBit(ctx, "bitoptest", 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), n)
}

func TestRedisExt_MSet(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("base/report", "test")

	resp, err := re.MSet(ctx, "setkey1", "setvalue1", "setkey2", "setvalue2")
	assert.NoError(t, err)
	assert.Equal(t, "OK", resp)
}

func TestRedisExt_MGet(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("base/report", "test")
	re.MSet(ctx, "getkey1", "getvalue1", "getkey2", "getvalue2")
	resp, err := re.MGet(ctx, []string{"getkey1", "getkey2"}...)
	fmt.Printf("resp:%+v\n", resp)
	assert.NoError(t, err)
	assert.Equal(t, len(resp), 2)
	assert.Contains(t, []string{"getvalue1", "getvalue2"}, resp[0])
	assert.Contains(t, []string{"getvalue1", "getvalue2"}, resp[1])
}

func TestRedisExt_HSetNX(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("base/report", "test")
	b, err := re.HSetNX(ctx, "hkey1", "field1", "value1")
	assert.NoError(t, err)
	assert.Equal(t, true, b)
	b, err = re.HSetNX(ctx, "hkey1", "field1", "value1")
	assert.NoError(t, err)
	assert.Equal(t, false, b)
	re.Del(ctx, "hkey1")
}

func TestRedisExt_TTL(t *testing.T) {
	ctx := context.Background()
	ttl := 10 * time.Second
	re := NewRedisExt("base/report", "test")
	re.Set(ctx,"getttl1", "test", ttl)
	d, err := re.TTL(ctx, "getttl1")
	assert.NoError(t, err)
	assert.Equal(t, ttl, d)
	d, err = re.TTL(ctx, "getttl2")
	assert.NoError(t, err)
	assert.Equal(t, -2 * time.Second, d)
	re.Set(ctx, "getttl3", "test", 0)
	d, err = re.TTL(ctx, "getttl3")
	assert.NoError(t, err)
	assert.Equal(t, -1 * time.Second, d)
	re.Del(ctx, "getttl3")
}
