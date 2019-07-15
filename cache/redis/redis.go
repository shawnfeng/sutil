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
}

func NewClient(ctx context.Context, namespace string) (*Client, error) {
	fun := "NewClient -->"

	config := DefaultConfiger.GetConfig(ctx, namespace)
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
	}, err
}

func (m *Client) fixKey(key string) string {
	return fmt.Sprintf("%s.%s", m.namespace, key)
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

func (m *Client) Close(ctx context.Context) error {
	return m.client.Close()
}
