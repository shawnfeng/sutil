package redisext

import (
	"context"
	go_redis "github.com/go-redis/redis"

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
		element, err = client.LIndex(ctx, m.prefixKeyWithContext(ctx, key), index).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) LIndex(ctx context.Context, key string, index int64) *go_redis.StringCmd {
	return p.pipe.LIndex(ctx, p.prefixKeyWithContext(ctx, key), index)
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
		n, err = client.LInsert(ctx, m.prefixKeyWithContext(ctx, key), op, pivot, value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *go_redis.IntCmd {
	return p.pipe.LInsert(ctx, p.prefixKeyWithContext(ctx, key), op, pivot, value)
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
		n, err = client.LLen(ctx, m.prefixKeyWithContext(ctx, key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) LLen(ctx context.Context, key string) *go_redis.IntCmd {
	return p.pipe.LLen(ctx, p.prefixKeyWithContext(ctx, key))
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
		element, err = client.LPop(ctx, m.prefixKeyWithContext(ctx, key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) LPop(ctx context.Context, key string) *go_redis.StringCmd {
	return p.pipe.LPop(ctx, p.prefixKeyWithContext(ctx, key))
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
		n, err = client.LPush(ctx, m.prefixKeyWithContext(ctx, key), values...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) LPush(ctx context.Context, key string, values ...interface{}) *go_redis.IntCmd {
	return p.pipe.LPush(ctx, p.prefixKeyWithContext(ctx, key), values)
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
		n, err = client.LPushX(ctx, m.prefixKeyWithContext(ctx, key), value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) LPushX(ctx context.Context, key string, value interface{}) *go_redis.IntCmd {
	return p.pipe.LPushX(ctx, p.prefixKeyWithContext(ctx, key), value)
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
		r, err = client.LRange(ctx, m.prefixKeyWithContext(ctx, key), start, stop).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) LRange(ctx context.Context, key string, start, stop int64) *go_redis.StringSliceCmd {
	return p.pipe.LRange(ctx, p.prefixKeyWithContext(ctx, key), start, stop)
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
		n, err = client.LRem(ctx, m.prefixKeyWithContext(ctx, key), count, value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) LRem(ctx context.Context, key string, count int64, value interface{}) *go_redis.IntCmd {
	return p.pipe.LRem(ctx, p.prefixKeyWithContext(ctx, key), count, value)
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
		r, err = client.LSet(ctx, m.prefixKeyWithContext(ctx, key), index, value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) LSet(ctx context.Context, key string, count int64, value interface{}) *go_redis.StatusCmd {
	return p.pipe.LSet(ctx, p.prefixKeyWithContext(ctx, key), count, value)
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
		r, err = client.LTrim(ctx, m.prefixKeyWithContext(ctx, key), start, stop).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) LTrim(ctx context.Context, key string, start, stop int64) *go_redis.StatusCmd {
	return p.pipe.LTrim(ctx, p.prefixKeyWithContext(ctx, key), start, stop)
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
		element, err = client.RPop(ctx, m.prefixKeyWithContext(ctx, key)).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) RPop(ctx context.Context, key string) *go_redis.StringCmd {
	return p.pipe.RPop(ctx, p.prefixKeyWithContext(ctx, key))
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
		n, err = client.RPush(ctx, m.prefixKeyWithContext(ctx, key), values...).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) RPush(ctx context.Context, key string, values ...interface{}) *go_redis.IntCmd {
	return p.pipe.RPush(ctx, p.prefixKeyWithContext(ctx, key), values)
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
		n, err = client.RPushX(ctx, m.prefixKeyWithContext(ctx, key), value).Result()
	}
	statReqErr(m.namespace, command, err)
	return
}

func (p *PipelineExt) RPushX(ctx context.Context, key string, value interface{}) *go_redis.IntCmd {
	return p.pipe.RPushX(ctx, p.prefixKeyWithContext(ctx, key), value)
}
