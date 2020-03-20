package redisext

import (
	"context"
	"strconv"
	"time"
)

const (
	retryTimeGap = time.Millisecond * 20

	InitState = "0"
	Doing     = "1"
	Done      = "99"
)

// Lock get a global lock identified by key, if failed, return quickly
func (m *RedisExt) Lock(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	r, err := m.SetNX(ctx, key, value, expiration)
	if err != nil {
		return false, err
	}
	return r, nil
}

// LockWithTimeout will retry get lock in timeout every retryTimeGap
func (m *RedisExt) LockWithTimeout(ctx context.Context, key string, value interface{}, expiration, timeout time.Duration) (r bool, err error) {
	for timeout > 0 {
		r, err = m.SetNX(ctx, key, value, expiration)
		if err != nil || !r {
			time.Sleep(retryTimeGap)
			timeout = timeout - retryTimeGap
		}
	}
	return r, err
}

// Unlock release lock with check lock value
func (m *RedisExt) Unlock(ctx context.Context, key string, value interface{}) (bool, error) {
	src := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	script := NewScript(src)
	r, err := m.Eval(ctx, script, []string{key}, value)
	if err != nil {
		return false, err
	}
	return r.(int64) == 1, nil
}

// TryAcquire 通过特殊的全局锁，保证以key为依据的业务操作是幂等的，可能的情况如下:
//
// 入口参数:
// @timeout: 业务操作的超时时间，比如: 2min
// @expiration: 业务记录过期时间，比如: 24h
//
// 返回参数:
// @canHandle: 是否可以继续执行该笔业务记录，
// @state: 该笔业务记录的状态，InitState：初始状态 Done: 处理完成 Doing：正在处理中
// @err: 异常错误
//
// 特别说明:
// 可以继续执行该笔业务记录的情况有: 记录为初始状态或者记录为Doing状态，且超过了timeout这个执行时间
func (m *RedisExt) TryAcquire(ctx context.Context, key string, expiration, timeout time.Duration) (canHandle bool, state string, err error) {
	value := time.Now().Add(timeout).Second()
	r, err := m.SetNX(ctx, key, value, expiration)
	if err != nil {
		return false, "unknown", err
	}
	// 加锁成功，可以继续业务处理
	if r {
		return true, InitState, nil
	}

	// 加锁失败，需要分情况判断
	// 状态已经被确认，不需要重新处理业务，可以根据状态进行自己所需操作，如直接返回成功
	s, err := m.Get(ctx, key)
	if s == Done {
		return false, Done, nil
	}

	// 当前时间小于业务处理过期时间，应该是有其他的进程/线程在处理，不需要处理业务，返回处理中
	now := time.Now().Second()
	result, _ := strconv.Atoi(s)
	if now < result {
		return false, Doing, nil
	}

	// 锁处于处理中，当前时间大于业务处理过期时间，类似乐观锁
	delta := now - result + int(timeout.Seconds())
	r1, err := m.IncrBy(ctx, key, int64(delta))
	if r1 == int64(result+delta) {
		return true, Doing, nil
	} else {
		m.DecrBy(ctx, key, int64(delta))
	}
	return false, Doing, nil
}

// Confirm 确定业务处理完成
func (m *RedisExt) Confirm(ctx context.Context, key string) error {
	_, err := m.Set(ctx, key, Done, 0)
	return err
}
