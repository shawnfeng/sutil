package redisext

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/stime"
)

type Script struct {
	src, hash string
}

// NewScript src: script content
func NewScript(src string) *Script {
	h := sha1.New()
	_, _ = io.WriteString(h, src)
	return &Script{
		src:  src,
		hash: hex.EncodeToString(h.Sum(nil)),
	}
}

// Hash return hash of script
func (s *Script) Hash() string {
	return s.hash
}

// ScriptLoad load script to redis server
func (m *RedisExt) ScriptLoad(ctx context.Context, script *Script) (r string, err error) {
	command := "redisext.ScriptLoad"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
		statReqErr(m.namespace, command, err)
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		r, err = client.ScriptLoad(ctx, script.src).Result()
	}
	return r, err
}

// ScriptExists check if script exists in redis server
func (m *RedisExt) ScriptExists(ctx context.Context, script *Script) (r bool, err error) {
	command := "redisext.ScriptExists"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
		statReqErr(m.namespace, command, err)
	}()
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		result, err := client.ScriptExists(ctx, script.src).Result()
		if err != nil {
			return false, err
		}
		if len(result) > 0 {
			r = result[0]
		}
	}
	return r, err
}

// Eval exec with script
func (m *RedisExt) Eval(ctx context.Context, script *Script, keys []string, args ...interface{}) (r interface{}, err error) {
	command := "redisext.Eval"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
		statReqErr(m.namespace,command, err)
	}()
	for i, key := range keys {
		keys[i] = m.prefixKey(key)
	}
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		r, err = client.Eval(ctx, script.src, keys, args...).Result()
	}
	return r, err
}

// EvalSha exec with script hash
func (m *RedisExt) EvalSha(ctx context.Context, script *Script, keys []string, args ...interface{}) (r interface{}, err error) {
	command := "redisext.EvalSha"
	span, ctx := opentracing.StartSpanFromContext(ctx, command)
	st := stime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(m.namespace, command, st.Millisecond())
		statReqErr(m.namespace,command, err)
	}()
	for i, key := range keys {
		keys[i] = m.prefixKey(key)
	}
	client, err := m.getRedisInstance(ctx)
	if err == nil {
		r, err = client.EvalSha(ctx, script.hash, keys, args...).Result()
	}
	return r, err
}
