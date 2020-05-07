package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/shawnfeng/sutil/cache/constants"
	"strings"
	"time"
)

type Pipeline struct {
	namespace  string
	pipeline   redis.Pipeliner
	opts       *options
}

func (m *Pipeline) fixKey(key string) string {
	if m.opts.noFixKey {
		return key
	}
	parts := []string{
		m.namespace,
		m.opts.wrapper,
		key,
	}
	if !m.opts.useWrapper {
		parts = []string{
			m.namespace,
			key,
		}
	}
	return strings.Join(parts, ".")
}

func (m *Pipeline) logSpan(ctx context.Context, op, key string) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.LogFields(
			log.String(constants.SpanLogOp, op),
			log.String(constants.SpanLogKeyKey, key),
			log.String(constants.SpanLogCacheType, fmt.Sprint(constants.CacheTypeRedis)))
	}
}


func (m *Pipeline) Get(ctx context.Context, key string) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.Get", k)
	return m.pipeline.Get(k)
}

func (m *Pipeline) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	var fixKeys = make([]string, len(keys))
	for k, v := range keys {
		key := m.fixKey(v)
		fixKeys[k] = key
	}
	m.logSpan(ctx, "Pipeline.MGet", strings.Join(fixKeys, "||"))
	return m.pipeline.MGet(fixKeys...)
}

func (m *Pipeline) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.Set", k)
	return m.pipeline.Set(k, value, expiration)
}

func (m *Pipeline) MSet(ctx context.Context, pairs ...interface{}) *redis.StatusCmd {
	var fixPairs = make([]interface{}, len(pairs))
	var keys []string
	for k, v := range pairs {
		if (k & 1) == 0 {
			key := m.fixKey(v.(string))
			keys = append(keys, key)
			fixPairs[k] = key
		} else {
			fixPairs[k] = v
		}
	}
	m.logSpan(ctx, "Pipeline.MSet", strings.Join(keys, "||"))
	return m.pipeline.MSet(fixPairs...)
}

func (m *Pipeline) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.GetBit", k)
	return m.pipeline.GetBit(k, offset)
}

func (m *Pipeline) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.SetBit", k)
	return m.pipeline.SetBit(k, offset, value)
}

func (m *Pipeline) Exists(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Exists", k)
	return m.pipeline.Exists(k)
}

func (m *Pipeline) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	var tkeys []string
	for _, key := range keys {
		tkeys = append(tkeys, m.fixKey(key))
	}

	m.logSpan(ctx, "Pipeline.Del", strings.Join(tkeys, ","))
	return m.pipeline.Del(tkeys...)
}

func (m *Pipeline) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.Expire", k)
	return m.pipeline.Expire(k, expiration)
}

func (m *Pipeline) Incr(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.Incr", k)
	return m.pipeline.Incr(k)
}

func (m *Pipeline) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.IncrBy", k)
	return m.pipeline.IncrBy(k, value)
}

func (m *Pipeline) Decr(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.Decr", k)
	return m.pipeline.Decr(k)
}

func (m *Pipeline) DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.DecrBy", k)
	return m.pipeline.DecrBy(k, value)
}

func (m *Pipeline) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.SetNX", k)
	return m.pipeline.SetNX(k, value, expiration)
}

func (m *Pipeline) HSet(ctx context.Context, key string, field string, value interface{}) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HSet", k)
	return m.pipeline.HSet(k, field, value)
}

func (m *Pipeline) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HDel", k)
	return m.pipeline.HDel(k, fields...)
}

func (m *Pipeline) HExists(ctx context.Context, key string, field string) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HExists", k)
	return m.pipeline.HExists(k, field)
}

func (m *Pipeline) HGet(ctx context.Context, key string, field string) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HGet", k)
	return m.pipeline.HGet(k, field)
}

func (m *Pipeline) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HGetAll", k)
	return m.pipeline.HGetAll(k)
}

func (m *Pipeline) HIncrBy(ctx context.Context, key string, field string, incr int64) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HIncrBy", k)
	return m.pipeline.HIncrBy(k, field, incr)
}

func (m *Pipeline) HIncrByFloat(ctx context.Context, key string, field string, incr float64) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HIncrByFloat", k)
	return m.pipeline.HIncrByFloat(k, field, incr)
}

func (m *Pipeline) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HKeys", k)
	return m.pipeline.HKeys(k)
}

func (m *Pipeline) HLen(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HLen", k)
	return m.pipeline.HLen(k)
}

func (m *Pipeline) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HMGet", k)
	return m.pipeline.HMGet(k, fields...)
}

func (m *Pipeline) HMSet(ctx context.Context, key string, fields map[string]interface{}) *redis.StatusCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HMSet", k)
	return m.pipeline.HMSet(k, fields)
}

func (m *Pipeline) HSetNX(ctx context.Context, key string, field string, val interface{}) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HSetNX", k)
	return m.pipeline.HSetNX(k, field, val)
}

func (m *Pipeline) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.HVals", k)
	return m.pipeline.HVals(k)
}

func (m *Pipeline) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZAdd", k)
	return m.pipeline.ZAdd(k, members...)
}

func (m *Pipeline) ZAddNX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZAddNX", k)
	return m.pipeline.ZAddNX(k, members...)
}

func (m *Pipeline) ZAddNXCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZAddNXCh", k)
	return m.pipeline.ZAddNXCh(k, members...)
}

func (m *Pipeline) ZAddXX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZAddXX", k)
	return m.pipeline.ZAddXX(k, members...)
}

func (m *Pipeline) ZAddXXCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZAddXXCh", k)
	return m.pipeline.ZAddXXCh(k, members...)
}

func (m *Pipeline) ZAddCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZAddCh", k)
	return m.pipeline.ZAddCh(k, members...)
}

func (m *Pipeline) ZCard(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZCard", k)
	return m.pipeline.ZCard(k)
}

func (m *Pipeline) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZCount", k)
	return m.pipeline.ZCount(k, min, max)
}

func (m *Pipeline) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZRange", k)
	return m.pipeline.ZRange(k, start, stop)
}

func (m *Pipeline) ZRangeByLex(ctx context.Context, key string, by redis.ZRangeBy) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZRangeByLex", k)
	return m.pipeline.ZRangeByLex(k, by)
}

func (m *Pipeline) ZRangeByScore(ctx context.Context, key string, by redis.ZRangeBy) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZRangeByScore", k)
	return m.pipeline.ZRangeByScore(k, by)
}

func (m *Pipeline) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZRangeWithScores", k)
	return m.pipeline.ZRangeWithScores(k, start, stop)
}

func (m *Pipeline) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZRevRange", k)
	return m.pipeline.ZRevRange(k, start, stop)
}

func (m *Pipeline) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZRevRangeWithScores", k)
	return m.pipeline.ZRevRangeWithScores(k, start, stop)
}

func (m *Pipeline) ZRank(ctx context.Context, key string, member string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZRank", k)
	return m.pipeline.ZRank(k, member)
}

func (m *Pipeline) ZRevRank(ctx context.Context, key string, member string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZRevRank", k)
	return m.pipeline.ZRevRank(k, member)
}

func (m *Pipeline) ZRem(ctx context.Context, key string, members []interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZRem", k)
	return m.pipeline.ZRem(k, members...)
}

func (m *Pipeline) ZIncr(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZIncr", k)
	return m.pipeline.ZIncr(k, member)
}

func (m *Pipeline) ZIncrNX(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZIncrNX", k)
	return m.pipeline.ZIncrNX(k, member)
}

func (m *Pipeline) ZIncrXX(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZIncrXX", k)
	return m.pipeline.ZIncrXX(k, member)
}

func (m *Pipeline) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZIncrBy", k)
	return m.pipeline.ZIncrBy(k, increment, member)
}

func (m *Pipeline) ZScore(ctx context.Context, key string, member string) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.ZScore", k)
	return m.pipeline.ZScore(k, member)
}

func (m *Pipeline) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.LIndex", k)
	return m.pipeline.LIndex(k, index)
}

func (m *Pipeline) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.LInsert", k)
	return m.pipeline.LInsert(k, op, pivot, value)
}

func (m *Pipeline) LLen(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.LLen", k)
	return m.pipeline.LLen(k)
}

func (m *Pipeline) LPop(ctx context.Context, key string) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.LPop", k)
	return m.pipeline.LPop(k)
}

func (m *Pipeline) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.LPush", k)
	return m.pipeline.LPush(k, values...)
}

func (m *Pipeline) LPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.LPushX", k)
	return m.pipeline.LPushX(k, value)
}

func (m *Pipeline) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.LRange", k)
	return m.pipeline.LRange(k, start, stop)
}

func (m *Pipeline) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.LRem", k)
	return m.pipeline.LRem(k, count, value)
}

func (m *Pipeline) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.LSet", k)
	return m.pipeline.LSet(k, index, value)
}

func (m *Pipeline) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.LTrim", k)
	return m.pipeline.LTrim(k, start, stop)
}

func (m *Pipeline) RPop(ctx context.Context, key string) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.RPop", k)
	return m.pipeline.RPop(k)
}

func (m *Pipeline) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.RPush", k)
	return m.pipeline.RPush(k, values...)
}

func (m *Pipeline) RPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.RPushX", k)
	return m.pipeline.RPushX(k, value)
}

func (m *Pipeline) TTL(ctx context.Context, key string) *redis.DurationCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Pipeline.TTL", k)
	return m.pipeline.TTL(k)
}

func (m *Pipeline) Exec(ctx context.Context) ([]redis.Cmder, error){
	return m.pipeline.Exec()
}

func (m *Pipeline) Discard(ctx context.Context) error {
	return m.pipeline.Discard()
}

func (m *Pipeline) Close() error {
	return m.pipeline.Close()
}


