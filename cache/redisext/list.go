package redisext

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/stime"
)

func (m *RedisExt) LIndex(ctx context.Context, key string, index int64) (element string, err error) {
	command := "redisext.LIndex"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		element, err = client.LIndex(ctx, m.prefixKey(key), index).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) LInsert(ctx context.Context, key, op string, pivot, value interface{}) (n int64, err error) {
	command := "redisext.LInsert"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.LInsert(ctx, m.prefixKey(key), op, pivot, value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) LLen(ctx context.Context, key string) (n int64, err error) {
	command := "redisext.LLen"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.LLen(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) LPop(ctx context.Context, key string) (element string, err error) {
	command := "redisext.LPop"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		element, err = client.LPop(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) LPush(ctx context.Context, key string, values ...interface{}) (n int64, err error) {
	command := "rdisext.LPush"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.LPush(ctx, m.prefixKey(key), values...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) LPushX(ctx context.Context, key string, value interface{}) (n int64, err error) {
	command := "redisext.LPushX"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.LPushX(ctx, m.prefixKey(key), value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) LRange(ctx context.Context, key string, start, stop int64) (r []string, err error) {
	command := "redisext.LRange"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		r, err = client.LRange(ctx, m.prefixKey(key), start, stop).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) LRem(ctx context.Context, key string, count int64, value interface{}) (n int64, err error) {
	command := "redisext.LRem"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.LRem(ctx, m.prefixKey(key), count, value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) LSet(ctx context.Context, key string, index int64, value interface{}) (r string, err error) {
	command := "redisext.LSet"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		r, err = client.LSet(ctx, m.prefixKey(key), index, value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) LTrim(ctx context.Context, key string, start, stop int64) (r string, err error) {
	command := "redisext.LTrim"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		r, err = client.LTrim(ctx, m.prefixKey(key), start, stop).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) RPop(ctx context.Context, key string) (element string, err error) {
	command := "redisext.RPop"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		element, err = client.RPop(ctx, m.prefixKey(key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) RPush(ctx context.Context, key string, values ...interface{}) (n int64, err error) {
	command := "redisext.RPush"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.RPush(ctx, m.prefixKey(key), values...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (m *RedisExt) RPushX(ctx context.Context, key string, value interface{}) (n int64, err error) {
	command := "redisext.RPushX"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		n, err = client.RPushX(ctx, m.prefixKey(key), value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}
