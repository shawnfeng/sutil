package cache

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/cache/constants"
	"github.com/shawnfeng/sutil/cache/redis"
	"github.com/shawnfeng/sutil/slog/slog"
	"time"
)

type CacheData interface {
	// 序列化接口
	Marshal() ([]byte, error)
	// 反序列化接口
	Unmarshal([]byte) error
	// cache miss load数据接口
	Load(key string) error
}

// 采用json进行序列化的的cache
type Cache struct {
	expire        int
	redisClient   *RedisClient
	prefix        string
	withNamespace bool
	namespace     string
}

// redis 地址列表，key前缀，过期时间
func NewCommonCache(serverName, prefix string, poolSize, expire int) (*Cache, error) {
	fun := "NewCommonCache-->"

	redisClient, err := NewCommonRedis(serverName, poolSize)
	if err != nil {
		slog.Errorf(context.TODO(), "%s NewCommonRedis, serverNam:%s err:%s", fun, serverName, err)
	}

	return &Cache{
		redisClient: redisClient,
		expire:      expire,
		prefix:      prefix,
	}, err
}

func NewCoreCache(serverName, prefix string, poolSize, expire int) (*Cache, error) {
	fun := "NewCoreCache-->"

	redisClient, err := NewCoreRedis(serverName, poolSize)
	if err != nil {
		slog.Errorf(context.TODO(), "%s NewCoreRedis, serverNam:%s err:%s", fun, serverName, err)
	}

	return &Cache{
		redisClient: redisClient,
		expire:      expire,
		prefix:      prefix,
	}, err
}

func NewCacheByNamespace(ctx context.Context, namespace, prefix string, expire int) (*Cache, error) {
	fun := "NewCacheByNamespace"

	client, err := NewRedisByNamespace(ctx, namespace)
	if err != nil {
		slog.Errorf(ctx, "%s GetConfig, namespace: %s err: %s", fun, namespace, err.Error())
	}
	return &Cache{
		expire:        expire,
		redisClient:   client,
		prefix:        prefix,
		withNamespace: true,
		namespace:     namespace,
	}, err
}

func (m *Cache) setData(key string, data CacheData) error {
	fun := "Cache.setData -->"
	expire := time.Duration(m.expire) * time.Second
	sdata, merr := data.Marshal()
	if merr != nil {
		sdata = []byte(merr.Error())
		merr = fmt.Errorf("%s marshal err, cache key:%s err:%s", fun, key, merr)
		expire = constants.CacheDirtyExpireTime
	}

	client, err := m.getRedisClient()
	if err != nil {
		return fmt.Errorf("%s get redis client err:%s", fun, err.Error())
	}

	err = client.Set(m.fixKey(key), sdata, expire).Err()
	if err != nil {
		return fmt.Errorf("%s set err, cache key:%s err:%s", fun, key, err)
	}

	if merr != nil {
		return merr
	}

	return nil
}

func (m *Cache) fixKey(key string) string {
	if len(m.prefix) > 0 {
		return fmt.Sprintf("%s.%s", m.prefix, key)
	}

	return key
}

func (m *Cache) getData(key string, data CacheData) error {
	fun := "Cache.getData -->"
	client, err := m.getRedisClient()
	if err != nil {
		return fmt.Errorf("%s get redis client err:%s", fun, err.Error())
	}
	sdata, err := client.Get(m.fixKey(key)).Bytes()
	if err != nil {
		return err
	}

	err = data.Unmarshal(sdata)
	if err != nil {
		return fmt.Errorf("reply data unmarshal err:%s", err)
	}

	return nil
}

func (m *Cache) GetCache(key string, data CacheData) error {
	fun := "Cache.GetCache -->"

	err := m.getData(key, data)
	if err == nil {
		return nil
	} else if err.Error() == constants.RedisNil {
		// 空的情况也返回正常
		return nil

	} else if err != nil {
		slog.Warnf(context.TODO(), "%s cache key:%s err:%s", fun, key, err)
		return err
	}

	return nil
}

func (m *Cache) Get(key string, data CacheData) error {
	fun := "Cache.Get -->"

	err := m.getData(key, data)
	if err == nil {
		return nil
	}

	if err != nil && err.Error() != constants.RedisNil {
		slog.Errorf(context.TODO(), "%s cache key:%s err:%s", fun, key, err)
		return fmt.Errorf("%s cache key:%s err:%s", fun, key, err)
	}

	slog.Infof(context.TODO(), "%s miss key:%s", fun, key)

	err = data.Load(key)
	if err != nil {
		slog.Warnf(context.TODO(), "%s load err, cache key:%s err:%s", fun, key, err)
		return err
	}

	err = m.setData(key, data)
	if err != nil {
		slog.Warnf(context.TODO(), "%s setData err, key:%s, err:%s", fun, key, err)
	}
	return err
}

func (m *Cache) Set(key string, data CacheData) error {
	fun := "Cache.Set -->"
	err := m.setData(key, data)
	if err != nil {
		slog.Errorf(context.TODO(), "%s setData err, key:%s, err:%s", fun, key, err)
	}
	return err
}

func (m *Cache) Del(key string) error {
	fun := "Cache.Del-->"

	client, err := m.getRedisClient()
	if err != nil {
		slog.Errorf(context.TODO(), "%s get redis client err:%s", fun, err.Error())
		return fmt.Errorf("%s get redis client err:%s", fun, err.Error())
	}
	err = client.Del(m.fixKey(key)).Err()
	if err != nil {
		slog.Errorf(context.TODO(), "del cache key:%s err:%s", key, err.Error())
		return fmt.Errorf("del cache key:%s err:%s", key, err.Error())
	}

	return nil
}

// get namespace update client
func (m *Cache) getRedisClient() (*RedisClient, error) {
	if m.withNamespace {
		return NewRedisByNamespace(context.Background(), m.namespace)
	}
	return m.redisClient, nil
}

func SetConfiger(ctx context.Context, configerType constants.ConfigerType) error {
	fun := "Cache.SetConfiger-->"
	configer, err := redis.NewConfiger(configerType)
	if err != nil {
		slog.Errorf(ctx, "%s create configer err:%v", fun, err)
		return err
	}
	slog.Infof(ctx, "%s %v configer created", fun, configerType)
	err = configer.Init(ctx)
	if err != nil {
		slog.Errorf(ctx, "%s init configer err:%v", fun, err)
	}
	redis.DefaultConfiger = configer
	return err
}

func WatchUpdate(ctx context.Context) {
	go redis.DefaultInstanceManager.Watch(ctx)
}

func init() {
	fun := "cache.init -->"
	ctx := context.Background()
	err := SetConfiger(ctx, constants.ConfigerTypeApollo)
	if err != nil {
		slog.Errorf(ctx, "%s set cache configer:%v err:%v", fun, constants.ConfigerTypeApollo, err)
	} else {
		slog.Infof(ctx, "%s cache configer:%v been set", fun, constants.ConfigerTypeApollo)
	}
	WatchUpdate(ctx)
}
