package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/shawnfeng/sutil/slog"
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

func (m *Client) Get(key string) *redis.StringCmd {
	return m.client.Get(m.fixKey(key))
}

func (m *Client) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return m.client.Set(m.fixKey(key), value, expiration)
}

func (m *Client) Del(keys ...string) *redis.IntCmd {
	var tkeys []string
	for _, key := range keys {
		tkeys = append(tkeys, m.fixKey(key))
	}
	return m.client.Del(tkeys...)
}

func (m *Client) Incr(key string) *redis.IntCmd {
	return m.client.Incr(m.fixKey(key))
}

func (m *Client) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return m.client.SetNX(m.fixKey(key), value, expiration)
}

func (m *Client) Close() error {
	return m.client.Close()
}
