package redisext

import (
	"context"
	"fmt"
	"time"

	redis2 "github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/cache"
	"github.com/shawnfeng/sutil/cache/constants"
	"github.com/shawnfeng/sutil/cache/redis"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"github.com/shawnfeng/sutil/stime"
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

type ZRangeBy struct {
	Min, Max      string
	Offset, Count int64
}

func toRedisZRangeBy(by ZRangeBy) redis2.ZRangeBy {
	return redis2.ZRangeBy{
		Min:    by.Min,
		Max:    by.Max,
		Offset: by.Offset,
		Count:  by.Count,
	}
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
		Group:     scontext.GetControlRouteGroupWithDefault(ctx, constants.DefaultRouteGroup),
		Namespace: m.namespace,
		Wrapper:   cache.WrapperTypeRedisExt,
	}
}

func (m *RedisExt) Get(ctx context.Context, key string) (s string, err error) {
	command := "redisext.Get"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		s, err = client.Get(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) MGet(ctx context.Context, keys ...string) (v []interface{}, err error) {
	command := "redisext.MGet"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		var prefixKey = make([]string, len(keys))
		for k, v := range keys {
			prefixKey[k] = m.prefixKey(v)
		}
		v, err = client.MGet(ctx, prefixKey...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) Set(ctx context.Context, key string, val interface{}, exp time.Duration) (s string, err error) {
	command := "redisext.Set"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		s, err = client.Set(ctx, m.prefixKey(key), val, exp).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) MSet(ctx context.Context, pairs ...interface{}) (s string, err error) {
	command := "redisext.MSet"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		var prefixPairs = make([]interface{}, len(pairs))
		for k, v := range pairs {
			if (k & 1) == 0 {
				prefixPairs[k] = m.prefixKey(v.(string))
			} else {
				prefixPairs[k] = v
			}
		}
		s, err = client.MSet(ctx, prefixPairs...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) GetBit(ctx context.Context, key string, offset int64) (n int64, err error) {
	command := "redisext.GetBit"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.GetBit(ctx, m.prefixKey(key), offset).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) SetBit(ctx context.Context, key string, offset int64, value int) (n int64, err error) {
	command := "redisext.SetBit"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.SetBit(ctx, m.prefixKey(key), offset, value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) Incr(ctx context.Context, key string) (n int64, err error) {
	command := "redisext.Incr"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.Incr(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) IncrBy(ctx context.Context, key string, val int64) (n int64, err error) {
	command := "redisext.IncrBy"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.IncrBy(ctx, m.prefixKey(key), val).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) Decr(ctx context.Context, key string) (n int64, err error) {
	command := "redisext.Decr"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.Decr(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) DecrBy(ctx context.Context, key string, val int64) (n int64, err error) {
	command := "redisext.DecrBy"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.DecrBy(ctx, m.prefixKey(key), val).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) SetNX(ctx context.Context, key string, val interface{}, exp time.Duration) (b bool, err error) {
	command := "redisext.SetNX"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		b, err = client.SetNX(ctx, m.prefixKey(key), val, exp).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) Exists(ctx context.Context, key string) (n int64, err error) {
	command := "redisext.Exists"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.Exists(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) Del(ctx context.Context, key string) (n int64, err error) {
	command := "redisext.Del"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.Del(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) Expire(ctx context.Context, key string, expiration time.Duration) (b bool, err error) {
	command := "redisext.Expire"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		b, err = client.Expire(ctx, m.prefixKey(key), expiration).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

// hashes apis
func (m *RedisExt) HSet(ctx context.Context, key string, field string, value interface{}) (b bool, err error) {
	command := "redisext.HSet"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		b, err = client.HSet(ctx, m.prefixKey(key), field, value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HDel(ctx context.Context, key string, fields ...string) (n int64, err error) {
	command := "redisext.HDel"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.HDel(ctx, m.prefixKey(key), fields...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HExists(ctx context.Context, key string, field string) (b bool, err error) {
	command := "redisext.HExists"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		b, err = client.HExists(ctx, m.prefixKey(key), field).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HGet(ctx context.Context, key string, field string) (s string, err error) {
	command := "redisext.HGet"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		s, err = client.HGet(ctx, m.prefixKey(key), field).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HGetAll(ctx context.Context, key string) (sm map[string]string, err error) {
	command := "redisext.HGetAll"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		sm, err = client.HGetAll(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HIncrBy(ctx context.Context, key string, field string, incr int64) (n int64, err error) {
	command := "redisext.HIncrBy"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.HIncrBy(ctx, m.prefixKey(key), field, incr).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HIncrByFloat(ctx context.Context, key string, field string, incr float64) (f float64, err error) {
	command := "redisext.HIncrByFloat"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.HIncrByFloat(ctx, m.prefixKey(key), field, incr).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HKeys(ctx context.Context, key string) (ss []string, err error) {
	command := "redisext.HKeys"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		ss, err = client.HKeys(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HLen(ctx context.Context, key string) (n int64, err error) {
	command := "redisext.HLen"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.HLen(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HMGet(ctx context.Context, key string, fields ...string) (vs []interface{}, err error) {
	command := "redisext.HMGet"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		vs, err = client.HMGet(ctx, m.prefixKey(key), fields...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HMSet(ctx context.Context, key string, fields map[string]interface{}) (s string, err error) {
	command := "redisext.HMSet"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		s, err = client.HMSet(ctx, m.prefixKey(key), fields).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HSetNX(ctx context.Context, key string, field string, val interface{}) (b bool, err error) {
	command := "redisext.HSetNX"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		b, err = client.HSet(ctx, m.prefixKey(key), field, val).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) HVals(ctx context.Context, key string) (ss []string, err error) {
	command := "redisext.HVals"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		ss, err = client.HVals(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

// sorted set apis
func (m *RedisExt) ZAdd(ctx context.Context, key string, members []Z) (n int64, err error) {
	command := "redisext.ZAdd"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAdd(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZAddNX(ctx context.Context, key string, members []Z) (n int64, err error) {
	command := "redisext.ZAddNX"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAddNX(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZAddNXCh(ctx context.Context, key string, members []Z) (n int64, err error) {
	command := "redisext.ZAddNXCh"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAddNXCh(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZAddXX(ctx context.Context, key string, members []Z) (n int64, err error) {
	command := "redisext.ZAddXX"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAddXX(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZAddXXCh(ctx context.Context, key string, members []Z) (n int64, err error) {
	command := "redisext.ZAddXXCh"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAddXXCh(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZAddCh(ctx context.Context, key string, members []Z) (n int64, err error) {
	command := "redisext.ZAddCh"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZAddCh(ctx, m.prefixKey(key), toRedisZSlice(members)...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZCard(ctx context.Context, key string) (n int64, err error) {
	command := "redisext.ZCard"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZCard(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZCount(ctx context.Context, key, min, max string) (n int64, err error) {
	command := "redisext.ZCount"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZCount(ctx, m.prefixKey(key), min, max).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZRange(ctx context.Context, key string, start, stop int64) (ss []string, err error) {
	command := "redisext.ZRange"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		ss, err = client.ZRange(ctx, m.prefixKey(key), start, stop).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZRangeByLex(ctx context.Context, key string, by ZRangeBy) (ss []string, err error) {
	command := "redisext.ZRangeByLex"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		ss, err = client.ZRangeByLex(ctx, m.prefixKey(key), toRedisZRangeBy(by)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZRangeByScore(ctx context.Context, key string, by ZRangeBy) (ss []string, err error) {
	command := "redisext.ZRangeByScore"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		ss, err = client.ZRangeByScore(ctx, m.prefixKey(key), toRedisZRangeBy(by)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZRangeWithScores(ctx context.Context, key string, start, stop int64) (zs []Z, err error) {
	command := "redisext.ZRangeWithScores"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	var rzs []redis2.Z
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		rzs, err = client.ZRangeWithScores(ctx, m.prefixKey(key), start, stop).Result()
		zs = fromRedisZSlice(rzs)
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZRevRange(ctx context.Context, key string, start, stop int64) (ss []string, err error) {
	command := "redisext.ZRevRange"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		ss, err = client.ZRevRange(ctx, m.prefixKey(key), start, stop).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) (zs []Z, err error) {
	command := "redisext.ZRevRangeWithScores"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	var rzs []redis2.Z
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		rzs, err = client.ZRevRangeWithScores(ctx, m.prefixKey(key), start, stop).Result()
		zs = fromRedisZSlice(rzs)
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZRank(ctx context.Context, key string, member string) (n int64, err error) {
	command := "redisext.ZRank"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZRank(ctx, m.prefixKey(key), member).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZRevRank(ctx context.Context, key string, member string) (n int64, err error) {
	command := "redisext.ZRevRank"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZRevRank(ctx, m.prefixKey(key), member).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZRem(ctx context.Context, key string, members []interface{}) (n int64, err error) {
	command := "redisext.ZRem"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.ZRem(ctx, m.prefixKey(key), members).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZIncr(ctx context.Context, key string, member Z) (f float64, err error) {
	command := "redisext.ZIncr"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.ZIncr(ctx, m.prefixKey(key), member.toRedisZ()).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZIncrNX(ctx context.Context, key string, member Z) (f float64, err error) {
	command := "redisext.ZIncrNX"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.ZIncrNX(ctx, m.prefixKey(key), member.toRedisZ()).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZIncrXX(ctx context.Context, key string, member Z) (f float64, err error) {
	command := "redisext.ZIncrXX"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.ZIncrXX(ctx, m.prefixKey(key), member.toRedisZ()).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZIncrBy(ctx context.Context, key string, increment float64, member string) (f float64, err error) {
	command := "redisext.ZIncrBy"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.ZIncrBy(ctx, m.prefixKey(key), increment, member).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) ZScore(ctx context.Context, key string, member string) (f float64, err error) {
	command := "redisext.ZScore"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		f, err = client.ZScore(ctx, m.prefixKey(key), member).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func SetConfiger(ctx context.Context, configerType constants.ConfigerType) error {
	fun := "Cache.SetConfiger-->"
	configer, err := redis.NewConfiger(configerType)
	if err != nil {
		slog.Errorf(ctx, "%s create configer err:%v", fun, err)
		return err
	}
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
	fun := "redisext.init -->"
	ctx := context.Background()
	err := SetConfiger(ctx, constants.ConfigerTypeApollo)
	if err != nil {
		slog.Errorf(ctx, "%s set redisext configer:%v err:%v", fun, constants.ConfigerTypeApollo, err)
	} else {
		slog.Infof(ctx, "%s redisext configer:%v been set", fun, constants.ConfigerTypeApollo)
	}
	WatchUpdate(ctx)
}
