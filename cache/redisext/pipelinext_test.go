package redisext

import (
	"context"
	"github.com/shawnfeng/sutil/cache/constants"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPipelineExt_Get(t *testing.T) {
	ctx := context.Background()
	redis := NewRedisExt("base/test", "test")
	redis.Set(ctx, "testPipeline", "success", 15 * time.Second)
	pipe, err := redis.Pipeline(ctx)
	assert.NoError(t, err)

	s := pipe.Get(ctx, "testPipeline")
	str, err := s.Result()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	cmds, err := pipe.Exec(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cmds))
	str, err = s.Result()
	assert.Equal(t, "success", str)
}

func TestPipelineExt_TTL(t *testing.T) {
	ctx := context.Background()
	redis := NewRedisExt("base/test", "test")
	_, err := redis.Set(ctx, "testPipeline", "success", 15 * time.Second)
	re, err := redis.Pipeline(ctx)
	assert.NoError(t, err)

	s := re.Get(ctx, "testPipeline")
	str, err := s.Result()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	tr := re.TTL(ctx, "testPipeline")
	dur, err := tr.Result()
	assert.NoError(t, err)
	assert.Equal(t, "0s", dur.String())

	cmds, err := re.Exec(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cmds))
	str, err = s.Result()
	assert.Equal(t, "success", str)
	dur, err = tr.Result()
	assert.Equal(t, "15s", dur.String())
}

func TestPipelineExt_MGet(t *testing.T) {
	ctx := context.Background()
	redis := NewRedisExt("base/test", "test")
	_, err := redis.MSet(ctx, "testPipeline", "success", "testPipeline2", "success")
	re, err := redis.Pipeline(ctx)
	assert.NoError(t, err)

	r := re.MGet(ctx, "testPipeline", "testPipeline2")
	arr, err := r.Result()
	assert.NoError(t, err)
	assert.Equal(t, []interface {}(nil), arr)

	cmds, err := re.Exec(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cmds))
	arr, err = r.Result()
	assert.Equal(t, "success", arr[0].(string))
	assert.Equal(t, "success", arr[1].(string))
}

func TestPipelineExt_NoPrefix(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.ContextCacheNoPrefix, false)

	redis := NewRedisExt("base/test", "test")
	redis.Set(ctx, "testPipeline", "success", 15 * time.Second)
	pipe, err := redis.Pipeline(ctx)
	assert.NoError(t, err)

	s := pipe.Get(ctx, "testPipeline")
	str, err := s.Result()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	cmds, err := pipe.Exec(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cmds))
	str, err = s.Result()
	assert.Equal(t, "success", str)
}

func TestPipelineExt_HGet(t *testing.T) {
	ctx := context.Background()

	redis := NewRedisExt("base/test", "test")
	redis.HSet(ctx, "testPipeline", "key", "val")
	pipe, err := redis.Pipeline(ctx)
	assert.NoError(t, err)

	strCmd := pipe.HGet(ctx, "testPipeline", "key")
	str, err := strCmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	cmds, err := pipe.Exec(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cmds))
	str, err = strCmd.Result()
	assert.Equal(t, "val", str)
}

func TestPipelineExt_LPop(t *testing.T) {
	ctx := context.Background()

	redis := NewRedisExt("base/test", "test")
	redis.LPush(ctx, "testPipelineLPop", "val1", "val2", "val3")
	pipe, err := redis.Pipeline(ctx)
	assert.NoError(t, err)

	strCmd1 := pipe.LPop(ctx, "testPipelineLPop")
	str, err := strCmd1.Result()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	strCmd2 := pipe.LPop(ctx, "testPipelineLPop")
	str, err = strCmd2.Result()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	strCmd3 := pipe.LPop(ctx, "testPipelineLPop")
	str, err = strCmd3.Result()
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	cmds, err := pipe.Exec(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(cmds))
	str, err = strCmd1.Result()
	assert.Equal(t, "val3", str)
	str, err = strCmd2.Result()
	assert.Equal(t, "val2", str)
	str, err = strCmd3.Result()
	assert.Equal(t, "val1", str)
}
