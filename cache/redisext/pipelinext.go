package redisext

import (
	"context"
	"fmt"
	go_redis "github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/cache/redis"
	"github.com/shawnfeng/sutil/stime"
	"time"
)

// PipelineExt by RedisExt get pipeline
type PipelineExt struct {
	namespace string
	prefix    string
	pipe      *redis.Pipeline
}

func (m *PipelineExt) prefixKey(key string) string {
	if len(m.prefix) > 0 {
		key = fmt.Sprintf("%s.%s", m.prefix, key)
	}
	return key
}

func (m *PipelineExt) Get(ctx context.Context, key string) (strCmd *go_redis.StringCmd) {
	return m.pipe.Get(ctx, m.prefixKey(key))
}

func (m *PipelineExt) MGet(ctx context.Context, keys ...string) (sliceCmd *go_redis.SliceCmd) {
	var prefixKey = make([]string, len(keys))
	for k, v := range keys {
		prefixKey[k] = m.prefixKey(v)
	}
	return m.pipe.MGet(ctx, prefixKey...)
}

func (m *PipelineExt) Set(ctx context.Context, key string, val interface{}, exp time.Duration) (statusCmd *go_redis.StatusCmd) {
	return m.pipe.Set(ctx, m.prefixKey(key), val, exp)
}

func (m *PipelineExt) MSet(ctx context.Context, pairs ...interface{}) (s *go_redis.StatusCmd) {
	var prefixPairs = make([]interface{}, len(pairs))
	for k, v := range pairs {
		if (k & 1) == 0 {
			prefixPairs[k] = m.prefixKey(v.(string))
		} else {
			prefixPairs[k] = v
		}
	}
	return m.pipe.MSet(ctx, prefixPairs...)
}

func (m *PipelineExt) GetBit(ctx context.Context, key string, offset int64) (n *go_redis.IntCmd) {
	return m.pipe.GetBit(ctx, m.prefixKey(key), offset)
}

func (m *PipelineExt) SetBit(ctx context.Context, key string, offset int64, value int) (n *go_redis.IntCmd) {
	return m.pipe.SetBit(ctx, m.prefixKey(key), offset, value)
}

func (m *PipelineExt) Incr(ctx context.Context, key string) (n *go_redis.IntCmd) {
	return m.pipe.Incr(ctx, m.prefixKey(key))
}

func (m *PipelineExt) IncrBy(ctx context.Context, key string, val int64) (n *go_redis.IntCmd) {
	return m.pipe.IncrBy(ctx, m.prefixKey(key), val)
}

func (m *PipelineExt) Decr(ctx context.Context, key string) (n *go_redis.IntCmd) {
	return m.pipe.Decr(ctx, m.prefixKey(key))
}

func (m *PipelineExt) DecrBy(ctx context.Context, key string, val int64) (n *go_redis.IntCmd) {
	return m.pipe.DecrBy(ctx, m.prefixKey(key), val)
}

func (m *PipelineExt) SetNX(ctx context.Context, key string, val interface{}, exp time.Duration) (b *go_redis.BoolCmd) {
	return m.pipe.SetNX(ctx, m.prefixKey(key), val, exp)
}

func (m *PipelineExt) Exists(ctx context.Context, key string) (n *go_redis.IntCmd) {
	return m.pipe.Exists(ctx, m.prefixKey(key))
}

func (m *PipelineExt) Del(ctx context.Context, key string) (n *go_redis.IntCmd) {
	return m.pipe.Del(ctx, m.prefixKey(key))
}

func (m *PipelineExt) Expire(ctx context.Context, key string, expiration time.Duration) (b *go_redis.BoolCmd) {
	return m.pipe.Expire(ctx, m.prefixKey(key), expiration)
}

// hashes apis
func (m *PipelineExt) HSet(ctx context.Context, key string, field string, value interface{}) (b *go_redis.BoolCmd) {
	return m.pipe.HSet(ctx, m.prefixKey(key), field, value)
}

func (m *PipelineExt) HDel(ctx context.Context, key string, fields ...string) (n *go_redis.IntCmd) {
	return m.pipe.HDel(ctx, m.prefixKey(key), fields...)
}

func (m *PipelineExt) HExists(ctx context.Context, key string, field string) (b *go_redis.BoolCmd) {
	return m.pipe.HExists(ctx, m.prefixKey(key), field)
}

func (m *PipelineExt) HGet(ctx context.Context, key string, field string) (s *go_redis.StringCmd) {
	return m.pipe.HGet(ctx, m.prefixKey(key), field)
}

func (m *PipelineExt) HGetAll(ctx context.Context, key string) (sm *go_redis.StringStringMapCmd) {
	return m.pipe.HGetAll(ctx, m.prefixKey(key))
}

func (m *PipelineExt) HIncrBy(ctx context.Context, key string, field string, incr int64) (n *go_redis.IntCmd) {
	return m.pipe.HIncrBy(ctx, m.prefixKey(key), field, incr)
}

func (m *PipelineExt) HIncrByFloat(ctx context.Context, key string, field string, incr float64) (f *go_redis.FloatCmd) {
	return m.pipe.HIncrByFloat(ctx, m.prefixKey(key), field, incr)
}

func (m *PipelineExt) HKeys(ctx context.Context, key string) (ss *go_redis.StringSliceCmd) {
	return m.pipe.HKeys(ctx, m.prefixKey(key))
}

func (m *PipelineExt) HLen(ctx context.Context, key string) (n *go_redis.IntCmd) {
	return 	m.pipe.HLen(ctx, m.prefixKey(key))
}

func (m *PipelineExt) HMGet(ctx context.Context, key string, fields ...string) (vs *go_redis.SliceCmd) {
	return m.pipe.HMGet(ctx, m.prefixKey(key), fields...)
}

func (m *PipelineExt) HMSet(ctx context.Context, key string, fields map[string]interface{}) (s *go_redis.StatusCmd) {
	return m.pipe.HMSet(ctx, m.prefixKey(key), fields)
}

func (m *PipelineExt) HSetNX(ctx context.Context, key string, field string, val interface{}) (b *go_redis.BoolCmd) {
	return m.pipe.HSetNX(ctx, m.prefixKey(key), field, val)
}

func (m *PipelineExt) HVals(ctx context.Context, key string) (ss *go_redis.StringSliceCmd) {
	return m.pipe.HVals(ctx, m.prefixKey(key))
}

// sorted set apis
func (m *PipelineExt) ZAdd(ctx context.Context, key string, members []Z) (n *go_redis.IntCmd) {
	return m.pipe.ZAdd(ctx, m.prefixKey(key), toRedisZSlice(members)...)
}

func (m *PipelineExt) ZAddNX(ctx context.Context, key string, members []Z) (n *go_redis.IntCmd) {
	return m.pipe.ZAddNX(ctx, m.prefixKey(key), toRedisZSlice(members)...)
}

func (m *PipelineExt) ZAddNXCh(ctx context.Context, key string, members []Z) (n *go_redis.IntCmd) {
	return m.pipe.ZAddNXCh(ctx, m.prefixKey(key), toRedisZSlice(members)...)
}

func (m *PipelineExt) ZAddXX(ctx context.Context, key string, members []Z) (n *go_redis.IntCmd) {
	return m.pipe.ZAddXX(ctx, m.prefixKey(key), toRedisZSlice(members)...)
}

func (m *PipelineExt) ZAddXXCh(ctx context.Context, key string, members []Z) (n *go_redis.IntCmd) {
	return m.pipe.ZAddXXCh(ctx, m.prefixKey(key), toRedisZSlice(members)...)
}

func (m *PipelineExt) ZAddCh(ctx context.Context, key string, members []Z) (n *go_redis.IntCmd) {
	return m.pipe.ZAddCh(ctx, m.prefixKey(key), toRedisZSlice(members)...)
}

func (m *PipelineExt) ZCard(ctx context.Context, key string) (n *go_redis.IntCmd) {
	return  m.pipe.ZCard(ctx, m.prefixKey(key))
}

func (m *PipelineExt) ZCount(ctx context.Context, key, min, max string) (n *go_redis.IntCmd) {
	return m.pipe.ZCount(ctx, m.prefixKey(key), min, max)
}

func (m *PipelineExt) ZRange(ctx context.Context, key string, start, stop int64) (ss *go_redis.StringSliceCmd) {
	return m.pipe.ZRange(ctx, m.prefixKey(key), start, stop)
}

func (m *PipelineExt) ZRangeByLex(ctx context.Context, key string, by ZRangeBy) (ss *go_redis.StringSliceCmd) {
	return m.pipe.ZRangeByLex(ctx, m.prefixKey(key), toRedisZRangeBy(by))
}

func (m *PipelineExt) ZRangeByScore(ctx context.Context, key string, by ZRangeBy) (ss *go_redis.StringSliceCmd) {
	return m.pipe.ZRangeByScore(ctx, m.prefixKey(key), toRedisZRangeBy(by))
}

func (m *PipelineExt) ZRangeWithScores(ctx context.Context, key string, start, stop int64) (zs *go_redis.ZSliceCmd) {
	return m.pipe.ZRangeWithScores(ctx, m.prefixKey(key), start, stop)
}

func (m *PipelineExt) ZRevRange(ctx context.Context, key string, start, stop int64) (ss *go_redis.StringSliceCmd) {
	return m.pipe.ZRevRange(ctx, m.prefixKey(key), start, stop)
}

func (m *PipelineExt) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) (zs *go_redis.ZSliceCmd) {
	return m.pipe.ZRevRangeWithScores(ctx, m.prefixKey(key), start, stop)
}

func (m *PipelineExt) ZRank(ctx context.Context, key string, member string) (n *go_redis.IntCmd) {
	return m.pipe.ZRank(ctx, m.prefixKey(key), member)
}

func (m *PipelineExt) ZRevRank(ctx context.Context, key string, member string) (n *go_redis.IntCmd) {
	return m.pipe.ZRevRank(ctx, m.prefixKey(key), member)
}

func (m *PipelineExt) ZRem(ctx context.Context, key string, members []interface{}) (n *go_redis.IntCmd) {
	return m.pipe.ZRem(ctx, m.prefixKey(key), members)
}

func (m *PipelineExt) ZIncr(ctx context.Context, key string, member Z) (f *go_redis.FloatCmd) {
	return m.pipe.ZIncr(ctx, m.prefixKey(key), member.toRedisZ())
}

func (m *PipelineExt) ZIncrNX(ctx context.Context, key string, member Z) (f *go_redis.FloatCmd) {
	return m.pipe.ZIncrNX(ctx, m.prefixKey(key), member.toRedisZ())
}

func (m *PipelineExt) ZIncrXX(ctx context.Context, key string, member Z) (f *go_redis.FloatCmd) {
	return m.pipe.ZIncrXX(ctx, m.prefixKey(key), member.toRedisZ())
}

func (m *PipelineExt) ZIncrBy(ctx context.Context, key string, increment float64, member string) (f *go_redis.FloatCmd) {
	return 	m.pipe.ZIncrBy(ctx, m.prefixKey(key), increment, member)
}

func (m *PipelineExt) ZScore(ctx context.Context, key string, member string) (f *go_redis.FloatCmd) {
	return m.pipe.ZScore(ctx, m.prefixKey(key), member)
}

func (m *PipelineExt) TTL(ctx context.Context, key string) (d *go_redis.DurationCmd) {
	return m.pipe.TTL(ctx, m.prefixKey(key))
}

func (m *PipelineExt) Exec(ctx context.Context) (cmds []go_redis.Cmder, err error) {
	command := "PipelineExt.Exec"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()

	cmds, err = m.pipe.Exec(ctx)
	statReqErr(m.namespace, command, err)
	return
}

func (m *PipelineExt) Discard(ctx context.Context) (err error) {
	command := "PipelineExt.Discard"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	err = m.pipe.Discard(ctx)
	statReqErr(m.namespace, command, err)
	return
}

func (m *PipelineExt) Close(ctx context.Context) error {
	return m.pipe.Close()
}
