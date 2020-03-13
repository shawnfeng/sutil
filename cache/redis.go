package cache

import (
	"context"
	"time"

	go_redis "github.com/go-redis/redis"
	"github.com/shawnfeng/sutil/cache/redis"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
)

const (
	defaultTimeout = time.Second * 2
)

func NewCommonRedis(serverName string, poolSize int) (*RedisClient, error) {
	return newRedisClient("common.codis.pri.ibanyu.com:19000", serverName, poolSize)
}

func NewCoreRedis(serverName string, poolSize int) (*RedisClient, error) {
	return newRedisClient("core.codis.pri.ibanyu.com:19000", serverName, poolSize)
}

func NewRedisByNamespace(ctx context.Context, namespace string) (*RedisClient, error) {
	fun := "NewRedisByNamespace -->"
	client, err := redis.DefaultInstanceManager.GetInstance(ctx, getInstanceConf(ctx, namespace))
	if err != nil {
		slog.Errorf(ctx, "%s GetInstance: namespace %s, err: %s", fun, namespace, err.Error())
	}
	return &RedisClient{
		client:    client,
		namespace: namespace,
	}, err
}

func getInstanceConf(ctx context.Context, namespace string) *redis.InstanceConf {
	return &redis.InstanceConf{
		Group:     scontext.GetControlRouteGroupWithDefault(ctx, DefaultRouteGroup),
		Namespace: namespace,
		Wrapper:   WrapperTypeCache,
	}
}

type RedisClient struct {
	client    *redis.Client
	namespace string
}

func newRedisClient(addr, serverName string, poolSize int) (*RedisClient, error) {
	fun := "newRedisClient-->"

	client, err := redis.NewDefaultClient(context.Background(), serverName, addr, WrapperTypeCache, poolSize, false, defaultTimeout)
	if err != nil {
		slog.Errorf(context.TODO(), "%s NewDefaultClient: serverName %s err: %s", fun, serverName, err.Error())
	}
	return &RedisClient{
		client:    client,
		namespace: serverName,
	}, err
}

func (m *RedisClient) Get(key string) *go_redis.StringCmd {
	ctx := context.Background()
	return m.client.Get(ctx, key)
}

func (m *RedisClient) Set(key string, value interface{}, expiration time.Duration) *go_redis.StatusCmd {
	ctx := context.Background()
	return m.client.Set(ctx, key, value, expiration)
}

func (m *RedisClient) Del(keys ...string) *go_redis.IntCmd {
	ctx := context.Background()
	return m.client.Del(ctx, keys...)
}

func (m *RedisClient) Incr(key string) *go_redis.IntCmd {
	ctx := context.Background()
	return m.client.Incr(ctx, key)
}

func (m *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) *go_redis.BoolCmd {
	ctx := context.Background()
	return m.client.SetNX(ctx, key, value, expiration)
}

func (m *RedisClient) TTL(key string) *go_redis.DurationCmd {
	ctx := context.Background()
	return m.client.TTL(ctx, key)
}

func (m *RedisClient) Expire(key string, expiration time.Duration) *go_redis.BoolCmd {
	ctx := context.Background()
	return m.client.Expire(ctx, key, expiration)
}
