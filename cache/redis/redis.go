package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/shawnfeng/sutil/cache/constants"
	"github.com/shawnfeng/sutil/slog/slog"
)

var RedisNil = fmt.Sprintf("redis: nil")

type Client struct {
	client     *redis.Client
	namespace  string
	wrapper    string
	useWrapper bool
}

func NewClient(ctx context.Context, namespace string, wrapper string) (*Client, error) {
	fun := "NewClient -->"

	config, err := DefaultConfiger.GetConfig(ctx, namespace)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.addr,
		DialTimeout:  3 * config.timeout,
		ReadTimeout:  config.timeout,
		WriteTimeout: config.timeout,
		PoolSize:     config.poolSize,
		PoolTimeout:  2 * config.timeout,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		slog.Errorf(ctx, "%s ping:%s err:%s", fun, pong, err)
	}

	return &Client{
		client:     client,
		namespace:  namespace,
		wrapper:    wrapper,
		useWrapper: config.useWrapper,
	}, err
}

func NewDefaultClient(ctx context.Context, namespace, addr, wrapper string, poolSize int, useWrapper bool, timeout time.Duration) (*Client, error) {
	fun := "NewDefaultClient -->"

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  3 * timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		PoolSize:     poolSize,
		PoolTimeout:  2 * timeout,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		slog.Errorf(ctx, "%s Ping: %s err: %s", fun, pong, err)
	}

	return &Client{
		client:     client,
		namespace:  namespace,
		wrapper:    wrapper,
		useWrapper: useWrapper,
	}, err
}

func (m *Client) fixKey(key string) string {
	parts := []string{
		m.namespace,
		m.wrapper,
		key,
	}
	if !m.useWrapper {
		parts = []string{
			m.namespace,
			key,
		}
	}
	return strings.Join(parts, ".")
}

func (m *Client) logSpan(ctx context.Context, op, key string) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.LogFields(
			log.String(constants.SpanLogOp, op),
			log.String(constants.SpanLogKeyKey, key),
			log.String(constants.SpanLogCacheType, fmt.Sprint(constants.CacheTypeRedis)))
	}
}

func (m *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Get", k)
	return m.client.Get(k)
}

func (m *Client) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	var fixKeys = make([]string, len(keys))
	for k, v := range keys {
		key := m.fixKey(v)
		fixKeys[k] = key
	}
	m.logSpan(ctx, "MGet", strings.Join(fixKeys, "||"))
	return m.client.MGet(fixKeys...)
}

func (m *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Set", k)
	return m.client.Set(k, value, expiration)
}

func (m *Client) MSet(ctx context.Context, pairs ...interface{}) *redis.StatusCmd {
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
	m.logSpan(ctx, "MSet", strings.Join(keys, "||"))
	return m.client.MSet(fixPairs...)
}

func (m *Client) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "GetBit", k)
	return m.client.GetBit(k, offset)
}

func (m *Client) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "SetBit", k)
	return m.client.SetBit(k, offset, value)
}

func (m *Client) Exists(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Exists", k)
	return m.client.Exists(k)
}

func (m *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	var tkeys []string
	for _, key := range keys {
		tkeys = append(tkeys, m.fixKey(key))
	}

	m.logSpan(ctx, "Del", strings.Join(tkeys, ","))
	return m.client.Del(tkeys...)
}

func (m *Client) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Expire", k)
	return m.client.Expire(k, expiration)
}

func (m *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Incr", k)
	return m.client.Incr(k)
}

func (m *Client) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "IncrBy", k)
	return m.client.IncrBy(k, value)
}

func (m *Client) Decr(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Decr", k)
	return m.client.Decr(k)
}

func (m *Client) DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "DecrBy", k)
	return m.client.DecrBy(k, value)
}

func (m *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "SetNX", k)
	return m.client.SetNX(k, value, expiration)
}

func (m *Client) HSet(ctx context.Context, key string, field string, value interface{}) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HSet", k)
	return m.client.HSet(k, field, value)
}

func (m *Client) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HDel", k)
	return m.client.HDel(k, fields...)
}

func (m *Client) HExists(ctx context.Context, key string, field string) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HExists", k)
	return m.client.HExists(k, field)
}

func (m *Client) HGet(ctx context.Context, key string, field string) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HGet", k)
	return m.client.HGet(k, field)
}

func (m *Client) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HGetAll", k)
	return m.client.HGetAll(k)
}

func (m *Client) HIncrBy(ctx context.Context, key string, field string, incr int64) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HIncrBy", k)
	return m.client.HIncrBy(k, field, incr)
}

func (m *Client) HIncrByFloat(ctx context.Context, key string, field string, incr float64) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HIncrByFloat", k)
	return m.client.HIncrByFloat(k, field, incr)
}

func (m *Client) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HKeys", k)
	return m.client.HKeys(k)
}

func (m *Client) HLen(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HLen", k)
	return m.client.HLen(k)
}

func (m *Client) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HMGet", k)
	return m.client.HMGet(k, fields...)
}

func (m *Client) HMSet(ctx context.Context, key string, fields map[string]interface{}) *redis.StatusCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HMSet", k)
	return m.client.HMSet(k, fields)
}

func (m *Client) HSetNX(ctx context.Context, key string, field string, val interface{}) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HSetNX", k)
	return m.client.HSetNX(k, field, val)
}

func (m *Client) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "HVals", k)
	return m.client.HVals(k)
}

func (m *Client) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZAdd", k)
	return m.client.ZAdd(k, members...)
}

func (m *Client) ZAddNX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZAddNX", k)
	return m.client.ZAddNX(k, members...)
}

func (m *Client) ZAddNXCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZAddNXCh", k)
	return m.client.ZAddNXCh(k, members...)
}

func (m *Client) ZAddXX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZAddXX", k)
	return m.client.ZAddXX(k, members...)
}

func (m *Client) ZAddXXCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZAddXXCh", k)
	return m.client.ZAddXXCh(k, members...)
}

func (m *Client) ZAddCh(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZAddCh", k)
	return m.client.ZAddCh(k, members...)
}

func (m *Client) ZCard(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZCard", k)
	return m.client.ZCard(k)
}

func (m *Client) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZCount", k)
	return m.client.ZCount(k, min, max)
}

func (m *Client) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZRange", k)
	return m.client.ZRange(k, start, stop)
}

func (m *Client) ZRangeByLex(ctx context.Context, key string, by redis.ZRangeBy) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZRangeByLex", k)
	return m.client.ZRangeByLex(k, by)
}

func (m *Client) ZRangeByScore(ctx context.Context, key string, by redis.ZRangeBy) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZRangeByScore", k)
	return m.client.ZRangeByScore(k, by)
}

func (m *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZRangeWithScores", k)
	return m.client.ZRangeWithScores(k, start, stop)
}

func (m *Client) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZRevRange", k)
	return m.client.ZRevRange(k, start, stop)
}

func (m *Client) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZRevRangeWithScores", k)
	return m.client.ZRevRangeWithScores(k, start, stop)
}

func (m *Client) ZRank(ctx context.Context, key string, member string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZRank", k)
	return m.client.ZRank(k, member)
}

func (m *Client) ZRevRank(ctx context.Context, key string, member string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZRevRank", k)
	return m.client.ZRevRank(k, member)
}

func (m *Client) ZRem(ctx context.Context, key string, members []interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZRem", k)
	return m.client.ZRem(k, members...)
}

func (m *Client) ZIncr(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZIncr", k)
	return m.client.ZIncr(k, member)
}

func (m *Client) ZIncrNX(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZIncrNX", k)
	return m.client.ZIncrNX(k, member)
}

func (m *Client) ZIncrXX(ctx context.Context, key string, member redis.Z) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZIncrXX", k)
	return m.client.ZIncrXX(k, member)
}

func (m *Client) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZIncrBy", k)
	return m.client.ZIncrBy(k, increment, member)
}

func (m *Client) ZScore(ctx context.Context, key string, member string) *redis.FloatCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "ZScore", k)
	return m.client.ZScore(k, member)
}

func (m *Client) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "LIndex", k)
	return m.client.LIndex(k, index)
}

func (m *Client) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "LInsert", k)
	return m.client.LInsert(k, op, pivot, value)
}

func (m *Client) LLen(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "LLen", k)
	return m.client.LLen(k)
}

func (m *Client) LPop(ctx context.Context, key string) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "LPop", k)
	return m.client.LPop(k)
}

func (m *Client) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "LPush", k)
	return m.client.LPush(k, values...)
}

func (m *Client) LPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "LPushX", k)
	return m.client.LPushX(k, value)
}

func (m *Client) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "LRange", k)
	return m.client.LRange(k, start, stop)
}

func (m *Client) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "LRem", k)
	return m.client.LRem(k, count, value)
}

func (m *Client) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "LSet", k)
	return m.client.LSet(k, index, value)
}

func (m *Client) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "LTrim", k)
	return m.client.LTrim(k, start, stop)
}

func (m *Client) RPop(ctx context.Context, key string) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "RPop", k)
	return m.client.RPop(k)
}

func (m *Client) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "RPush", k)
	return m.client.RPush(k, values...)
}

func (m *Client) RPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "RPushX", k)
	return m.client.RPushX(k, value)
}

func (m *Client) TTL(ctx context.Context, key string) *redis.DurationCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "TTL", k)
	return m.client.TTL(k)
}

func (m *Client) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
	m.logSpan(ctx, "ScriptLoad", script)
	return m.client.ScriptLoad(script)
}

func (m *Client) ScriptExists(ctx context.Context, scriptHash string) *redis.BoolSliceCmd {
	m.logSpan(ctx, "ScriptExists", scriptHash)
	return m.client.ScriptExists(scriptHash)
}

func (m *Client) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	m.logSpan(ctx, "Eval", script)
	for i, key := range keys {
		keys[i] = m.fixKey(key)
	}
	return m.client.Eval(script, keys, args...)
}

func (m *Client) EvalSha(ctx context.Context, scriptHash string, keys []string, args ...interface{}) *redis.Cmd {
	m.logSpan(ctx, "EvalSha", scriptHash)
	for i, key := range keys {
		keys[i] = m.fixKey(key)
	}
	return m.client.EvalSha(scriptHash, keys, args...)
}

func (m *Client) Close(ctx context.Context) error {
	return m.client.Close()
}
