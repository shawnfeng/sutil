package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/shawnfeng/sutil/cache"
	"github.com/shawnfeng/sutil/slog/slog"
	"strings"
	"time"
)

var RedisNil = fmt.Sprintf("redis: nil")

type Client struct {
	client    *redis.Client
	namespace string
	wrapper   string
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
		client:    client,
		namespace: namespace,
		wrapper:   wrapper,
	}, err
}

func (m *Client) fixKey(key string) string {
	return strings.Join([]string{
		m.namespace,
		m.wrapper,
		key,
	}, ".")
}

func (m *Client) logSpan(ctx context.Context, op, key string) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.LogFields(
			log.String(cache.SpanLogOp, op),
			log.String(cache.SpanLogKeyKey, key),
			log.String(cache.SpanLogCacheType, fmt.Sprint(cache.CacheTypeRedis)))
	}
}

func (m *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Get", k)
	return m.client.Get(k)
}

func (m *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Set", k)
	return m.client.Set(k, value, expiration)
}

func (m *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	var tkeys []string
	for _, key := range keys {
		tkeys = append(tkeys, m.fixKey(key))
	}

	m.logSpan(ctx, "Del", strings.Join(tkeys, ","))
	return m.client.Del(tkeys...)
}

func (m *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "Incr", k)
	return m.client.Incr(k)
}

func (m *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	k := m.fixKey(key)
	m.logSpan(ctx, "SetNX", k)
	return m.client.SetNX(k, value, expiration)
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

func (m *Client) Close(ctx context.Context) error {
	return m.client.Close()
}
