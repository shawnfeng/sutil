package cache

import (
	"fmt"
	"github.com/shawnfeng/sutil/slog"
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
	expire      int
	redisClient *RedisClient
	prefix      string
}

// redis 地址列表，key前缀，过期时间
func NewCommonCache(serverName, prefix string, poolSize, expire int) (*Cache, error) {
	fun := "NewCommonCache-->"

	redisClient, err := NewCommonRedis(serverName, poolSize)
	if err != nil {
		slog.Errorf("%s NewCommonRedis, serverNam:%s err:%s", fun, serverName, err)
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
		slog.Errorf("%s NewCoreRedis, serverNam:%s err:%s", fun, serverName, err)
	}

	return &Cache{
		redisClient: redisClient,
		expire:      expire,
		prefix:      prefix,
	}, err
}

func (m *Cache) setData(key string, data CacheData) error {
	fun := "Cache.setData -->"

	sdata, merr := data.Marshal()
	if merr != nil {
		slog.Errorf("%s marshal err, cache key:%s err:%s", fun, key, merr)
		sdata = []byte(merr.Error())
	}

	err := m.redisClient.Set(m.fixKey(key), sdata, time.Duration(m.expire)*time.Second).Err()
	if err != nil {
		slog.Errorf("%s set err, cache key:%s err:%s", fun, key, err)
	}

	if merr != nil {
		return merr
	}

	return err
}

func (m *Cache) fixKey(key string) string {
	if len(m.prefix) > 0 {
		return fmt.Sprintf("%s.%s", m.prefix, key)
	}

	return key
}

func (m *Cache) getData(key string, data CacheData) error {
	//fun := "Cache.getData -->"

	sdata, err := m.redisClient.Get(m.fixKey(key)).Bytes()
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
	} else if err.Error() == RedisNil {
		// 空的情况也返回正常
		return nil

	} else if err != nil {
		slog.Warnf("%s cache key:%s err:%s", fun, key, err)
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

	if err != nil && err.Error() != RedisNil {
		slog.Errorf("%s cache key:%s err:%s", fun, key, err)
		return fmt.Errorf("%s cache key:%s err:%s", fun, key, err)
	}

	slog.Infof("%s miss key:%s", fun, key)

	err = data.Load(key)
	if err != nil {
		slog.Warnf("%s load err, cache key:%s err:%s", fun, key, err)
	}

	return m.setData(key, data)

}

func (m *Cache) Set(key string, data CacheData) error {
	//fun := "Cache.Set -->"

	return m.setData(key, data)
}

func (m *Cache) Del(key string) error {
	//fun := "Cache.Del-->"

	err := m.redisClient.Del(m.fixKey(key)).Err()
	if err != nil {
		return fmt.Errorf("del cache key:%s err:%s", key, err.Error())
	}

	return nil
}
