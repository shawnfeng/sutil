package redisext

import (
	"context"
	"fmt"
	redis2 "github.com/go-redis/redis"
	"github.com/shawnfeng/sutil/cache"
	"github.com/shawnfeng/sutil/cache/redis"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"time"
)

type RedisExt struct {
	namespace string
	prefix    string
}

func NewRedisExt(namespace, prefix string) *RedisExt {
	return &RedisExt{namespace, prefix}
}

type Z struct {
	Score  float64
	Member interface{}
}

func (z Z) toRedisZ() redis2.Z {
	return redis2.Z{
		Score:  z.Score,
		Member: z.Member,
	}
}

func fromRedisZ(rz redis2.Z) Z {
	return Z{
		Score:  rz.Score,
		Member: rz.Member,
	}
}

func toRedisZSlice(zs []Z) (rzs []redis2.Z) {
	for _, z := range zs {
		rzs = append(rzs, z.toRedisZ())
	}
	return
}

func fromRedisZSlice(rzs []redis2.Z) (zs []Z) {
	for _, rz := range rzs {
		zs = append(zs, fromRedisZ(rz))
	}
	return
}

func (m *RedisExt) prefixKey(key string) string {
	if len(m.prefix) > 0 {
		key = fmt.Sprintf("%s.%s", m.prefix, key)
	}
	return key
}

func (m *RedisExt) getRedisInstance(ctx context.Context) (client *redis.Client, err error) {
	conf := m.getInstanceConf(ctx)
	return redis.DefaultInstanceManager.GetInstance(ctx, conf)
}

func (m *RedisExt) getInstanceConf(ctx context.Context) *redis.InstanceConf {
	return &redis.InstanceConf{
		Group:     scontext.GetGroupWithDefault(ctx, cache.DefaultRouteGroup),
		Namespace: m.namespace,
		Wrapper:   cache.WrapperTypeRedisExt,
	}
}

func (m *RedisExt) Get(ctx context.Context, key string) (s string, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		s, err = client.Get(ctx, m.prefixKey(key)).Result()
	}
	return
}

func (m *RedisExt) Set(ctx context.Context, key string, val interface{}, exp time.Duration) (s string, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		s, err = client.Set(ctx, m.prefixKey(key), val, exp).Result()
	}
	return
}

func (m *RedisExt) Exists(ctx context.Context, key string) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.Exists(ctx, m.prefixKey(key)).Result()
	}
	return
}

func (m *RedisExt) Del(ctx context.Context, key string) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.Del(ctx, m.prefixKey(key)).Result()
	}
	return
}

func (m *RedisExt) Expire(ctx context.Context, key string, expiration time.Duration) (b bool, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		b, err = client.Expire(ctx, m.prefixKey(key), expiration).Result()
	}
	return
}

func (m *RedisExt) ZAdd(ctx context.Context, key string, members []Z) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAdd(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	return
}

func (m *RedisExt) ZAddNX(ctx context.Context, key string, members []Z) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAddNX(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	return
}

func (m *RedisExt) ZAddNXCh(ctx context.Context, key string, members []Z) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAddNXCh(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	return
}

func (m *RedisExt) ZAddXX(ctx context.Context, key string, members []Z) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAddXX(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	return
}

func (m *RedisExt) ZAddXXCh(ctx context.Context, key string, members []Z) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAddXXCh(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	return
}

func (m *RedisExt) ZAddCh(ctx context.Context, key string, members []Z) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAddCh(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	return
}

func (m *RedisExt) ZCard(ctx context.Context, key string) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZCard(ctx, m.prefixKey(key)).Result()
	}
	return
}

func (m *RedisExt) ZCount(ctx context.Context, key, min, max string) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZCount(ctx, m.prefixKey(key), min, max).Result()
	}
	return
}

func (m *RedisExt) ZRange(ctx context.Context, key string, start, stop int64) (ss []string, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		ss, err = client.ZRange(ctx, m.prefixKey(key), start, stop).Result()
	}
	return
}

func (m *RedisExt) ZRangeWithScores(ctx context.Context, key string, start, stop int64) (zs []Z, err error) {
	var rzs []redis2.Z
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		rzs, err = client.ZRangeWithScores(ctx, m.prefixKey(key), start, stop).Result()
		zs = fromRedisZSlice(rzs)
	}
	return
}

func (m *RedisExt) ZRevRange(ctx context.Context, key string, start, stop int64) (ss []string, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		ss, err = client.ZRevRange(ctx, m.prefixKey(key), start, stop).Result()
	}
	return
}

func (m *RedisExt) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) (zs []Z, err error) {
	var rzs []redis2.Z
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		rzs, err = client.ZRevRangeWithScores(ctx, m.prefixKey(key), start, stop).Result()
		zs = fromRedisZSlice(rzs)
	}
	return
}

func (m *RedisExt) ZRank(ctx context.Context, key string, member string) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZRank(ctx, m.prefixKey(key), member).Result()
	}
	return
}

func (m *RedisExt) ZRevRank(ctx context.Context, key string, member string) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZRevRank(ctx, m.prefixKey(key), member).Result()
	}
	return
}

func (m *RedisExt) ZRem(ctx context.Context, key string, members []interface{}) (n int64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZRem(ctx, m.prefixKey(key), members).Result()
	}
	return
}

func (m *RedisExt) ZIncr(ctx context.Context, key string, member Z) (f float64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.ZIncr(ctx, m.prefixKey(key), member.toRedisZ()).Result()
	}
	return
}

func (m *RedisExt) ZIncrNX(ctx context.Context, key string, member Z) (f float64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.ZIncrNX(ctx, m.prefixKey(key), member.toRedisZ()).Result()
	}
	return
}

func (m *RedisExt) ZIncrXX(ctx context.Context, key string, member Z) (f float64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.ZIncrXX(ctx, m.prefixKey(key), member.toRedisZ()).Result()
	}
	return
}

func (m *RedisExt) ZIncrBy(ctx context.Context, key string, increment float64, member string) (f float64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.ZIncrBy(ctx, m.prefixKey(key), increment, member).Result()
	}
	return
}

func (m *RedisExt) ZScore(ctx context.Context, key string, member string) (f float64, err error) {
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.ZScore(ctx, m.prefixKey(key), member).Result()
	}
	return
}

func SetConfiger(ctx context.Context, configerType cache.ConfigerType) error {
	fun := "Cache.SetConfiger-->"
	configer, err := redis.NewConfiger(configerType)
	if err != nil {
		slog.Errorf(ctx, "%s create configer err:%v", fun, err)
		return err
	}
	slog.Infof(ctx, "%s %v configer created", fun, configerType)
	redis.DefaultConfiger = configer
	return redis.DefaultConfiger.Init(ctx)
}

func WatchUpdate(ctx context.Context) {
	go redis.DefaultInstanceManager.Watch(ctx)
}

func init() {
	fun := "redisext.init -->"
	ctx := context.Background()
	err := SetConfiger(ctx, cache.ConfigerTypeApollo)
	if err != nil {
		slog.Errorf(ctx, "%s set redisext configer:%v err:%v", fun, cache.ConfigerTypeApollo, err)
	} else {
		slog.Infof(ctx, "%s redisext configer:%v been set", fun, cache.ConfigerTypeApollo)
	}
}
