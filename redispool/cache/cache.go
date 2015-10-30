package rediscache

import (
	"errors"
	"fmt"

	"github.com/fzzy/radix/redis"
	"github.com/shawnfeng/sutil"
	"github.com/shawnfeng/sutil/slog"
	"github.com/shawnfeng/sutil/redispool"

)

var redisPool  *redispool.RedisPool

var errNilReply = errors.New("nil reply")


func init() {
	redisPool = redispool.NewRedisPool()
}

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
	expire int
	pref string
	addrs []string
}

// redis 地址列表，key前缀，过期时间
func NewCache(addrs []string, pref string, expire int) *Cache {
	return &Cache {
		addrs: addrs,
		pref: pref,
		expire: expire,
	}
}

func (m *Cache) setData(key string, data CacheData) error {
	//fun := "Cache.setData -->"

	sdata, err := data.Marshal()
	if err != nil {
		return fmt.Errorf("data marshal:%s", err)
	}

	rd := sutil.HashV(m.addrs, key)

	// 消息id到data的映射
	rp := redisPool.CmdSingle(
		rd,
		[]interface{}{
			"SETEX",
			fmt.Sprintf("%s.%s", m.pref, key),
			m.expire,
			sdata,
		},
	)

	if rp.Type == redis.ErrorReply {
		return fmt.Errorf("set cache err:%s", rp.String()) 
	}

	return nil

}

func (m *Cache) getData(key string, data CacheData) error {
	//fun := "Cache.getData -->"

	rd := sutil.HashV(m.addrs, key)

	rp := redisPool.CmdSingle(
		rd,
		[]interface{}{
			"GET",
			fmt.Sprintf("%s.%s", m.pref, key),
		},
	)

	if rp.Type == redis.ErrorReply {
		return fmt.Errorf("get cache err:%s", rp.String()) 

	} else if rp.Type == redis.NilReply {
		return errNilReply

	} else {
		sdata, err := rp.Bytes()
		if err != nil {
			return fmt.Errorf("reply bytes err:%s", err) 
		}

		err = data.Unmarshal(sdata)
		if err != nil {
			return fmt.Errorf("reply data unmarshal err:%s", err) 
		}
	}

	return nil

}

func (m *Cache) Get(key string, data CacheData) error {
	fun := "Cache.Get -->"

	err := m.getData(key, data)
	if err == nil {
		return nil
	} else if err == errNilReply {
		slog.Infof("%s miss key:%s", fun, key)

	} else if err != nil {
		slog.Warnf("%s cache key:%s err:%s", fun, key, err)

	}

	// miss load
	err = data.Load(key)
	if err != nil {
		return err
	}

	return m.setData(key, data)

}

// 不会自动load
func (m *Cache) Set(key string, data CacheData) error {
	//fun := "Cache.Set -->"


	return m.setData(key, data)

}


func (m *Cache) Del(key string) error {
	rd := sutil.HashV(m.addrs, key)

	rp := redisPool.CmdSingle(
		rd,
		[]interface{}{
			"DEL",
			fmt.Sprintf("%s.%s", m.pref, key),
		},
	)

	if rp.Type == redis.ErrorReply {
		return fmt.Errorf("del cache err:%s", rp.String())
	}

	//slog.Infof("%s", rp)

	return nil

}
