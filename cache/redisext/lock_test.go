package redisext

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLock(t *testing.T) {
	key := "test"
	value := "76c8c07d15c32fa90dd2b89f141246e8"
	ctx := context.Background()
	m := NewRedisExt("test/test", "test")
	r, err := m.Lock(ctx, key, value, time.Second*1)
	assert.Nil(t, err)
	assert.Equal(t, true, r)
	r1, err := m.Lock(ctx, key, value, time.Second*1)
	assert.Nil(t, err)
	assert.Equal(t, false, r1)
	time.Sleep(time.Second * 1)
	r2, err := m.Lock(ctx, key, value, time.Second*1)
	assert.Nil(t, err)
	assert.Equal(t, true, r2)
}

func TestLockWithTimeout(t *testing.T) {
	key := "test"
	value := "76c8c07d15c32fa90dd2b89f141246e8"
	ctx := context.Background()
	m := NewRedisExt("test/test", "test")
	m.Lock(ctx, key, value, time.Second*5)
	r, err := m.LockWithTimeout(ctx, key, value, time.Second, time.Second*1)
	assert.Nil(t, err)
	assert.Equal(t, false, r)
}

func TestUnlock(t *testing.T) {
	key := "test"
	value := "76c8c07d15c32fa90dd2b89f141246e8"
	ctx := context.Background()
	m := NewRedisExt("test/test", "test")
	r, err := m.Lock(ctx, key, value, time.Second*5)
	assert.Nil(t, err)
	assert.Equal(t, true, r)
	r1, err := m.Unlock(ctx, key, value)
	assert.Nil(t, err)
	assert.Equal(t, true, r1)
	r2, err := m.Lock(ctx, key, value, time.Second*1)
	assert.Nil(t, err)
	assert.Equal(t, true, r2)
}

func TestTryAcquire(t *testing.T) {
	orderID := "76c8c07d15c32fa90dd2b89f141246e8"
	ctx := context.Background()
	m := NewRedisExt("test/test", "test")
	m.Del(ctx, orderID)

	// Acquired first time
	canHandle, state, err := m.TryAcquire(ctx, orderID, time.Hour*24, time.Second*2)
	assert.Nil(t, err)
	assert.Equal(t, true, canHandle)
	assert.Equal(t, InitState, state)

	// reject cause by other processor
	canHandle1, state1, err := m.TryAcquire(ctx, orderID, time.Hour*24, time.Second*3)
	assert.Nil(t, err)
	assert.Equal(t, false, canHandle1)
	assert.Equal(t, Doing, state1)

	time.Sleep(time.Second * 2)

	// Acquired cause by other processor timeout
	canHandle2, state2, err := m.TryAcquire(ctx, orderID, time.Hour*24, time.Second*3)
	assert.Nil(t, err)
	assert.Equal(t, true, canHandle2)
	assert.Equal(t, Doing, state2)

	// confirm the record
	err = m.Confirm(ctx, orderID)
	assert.Nil(t, err)

	// idempotent
	canHandle3, state3, err := m.TryAcquire(ctx, orderID, time.Hour*24, time.Second*3)
	assert.Nil(t, err)
	assert.Equal(t, false, canHandle3)
	assert.Equal(t, Done, state3)
}
