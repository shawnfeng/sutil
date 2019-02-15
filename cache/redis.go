package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/shawnfeng/sutil/slog"
	"time"
)

var RedisNil = fmt.Sprintf("redis: nil")

func NewCommonRedis(serverName string, poolSize int) (*RedisClient, error) {
	return newRedisClient("common.codis.pri.ibanyu.com:19000", serverName, poolSize)
}

type RedisClient struct {
	client    *redis.Client
	namespace string
}

func newRedisClient(addr, serverName string, poolSize int) (*RedisClient, error) {
	fun := "newRedisClient-->"

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     poolSize,
		PoolTimeout:  4 * time.Second,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		slog.Errorf("%s ping:%s err:%s", fun, pong, err)
	}

	return &RedisClient{
		client:    client,
		namespace: serverName,
	}, err
}

func (m *RedisClient) fixKey(key string) string {
	return fmt.Sprintf("%s.%s", m.namespace, key)
}

func (m *RedisClient) Get(key string) *redis.StringCmd {
	return m.client.Get(m.fixKey(key))
}

func (m *RedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return m.client.Set(m.fixKey(key), value, expiration)
}

func (m *RedisClient) Del(keys ...string) *redis.IntCmd {
	var tkeys []string
	for _, key := range keys {
		tkeys = append(tkeys, m.fixKey(key))
	}
	return m.client.Del(tkeys...)
}

func (m *RedisClient) Incr(key string) *redis.IntCmd {
	return m.client.Incr(m.fixKey(key))
}
