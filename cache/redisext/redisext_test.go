package redisext

import (
	"context"
	"github.com/shawnfeng/sutil/cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	zsetTestKey = "myzset"
)

func TestRedisExt_ZAdd(t *testing.T) {
	ctx := context.Background()

	re := NewRedisExt("base/report", "test")
	_ = SetConfiger(ctx, cache.ConfigerTypeApollo)

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
	_ = SetConfiger(ctx, cache.ConfigerTypeApollo)

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
	_ = SetConfiger(ctx, cache.ConfigerTypeApollo)

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

	// cleanup
	dn, err := re.Del(ctx, zsetTestKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dn)
}

func TestRedisExt_ZCount(t *testing.T) {
	ctx := context.Background()
	re := NewRedisExt("base/report", "test")
	_ = SetConfiger(ctx, cache.ConfigerTypeApollo)

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
	_ = SetConfiger(ctx, cache.ConfigerTypeApollo)

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

	// cleanup
	dn, err := re.Del(ctx, zsetTestKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dn)
}
