package scontext

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testHead struct {
	Uid     int64  `thrift:"uid,1" json:"uid"`
	Source  int32  `thrift:"source,2" json:"source"`
	Ip      string `thrift:"ip,3" json:"ip"`
	Region  string `thrift:"region,4" json:"region"`
	Dt      int32  `thrift:"dt,5" json:"dt"`
	Unionid string `thrift:"unionid,6" json:"unionid"`
}

func (m testHead) ToKV() map[string]interface{} {
	return map[string]interface{}{
		"uid":     m.Uid,
		"source":  m.Source,
		"ip":      m.Ip,
		"region":  m.Region,
		"dt":      m.Dt,
		"unionid": m.Unionid,
	}
}

func TestGetUid(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyHead, &testHead{
		Uid: int64(100),
	})

	uid, ok := GetUid(ctx)
	assert.True(t, ok)
	assert.Equal(t, int64(100), uid)
}

func TestGetSource(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyHead, &testHead{
		Source: int32(1),
	})

	source, ok := GetSource(ctx)
	assert.True(t, ok)
	assert.Equal(t, int32(1), source)
}

func TestGetIp(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyHead, &testHead{
		Ip: "192.168.0.1",
	})
	ip, ok := GetIp(ctx)
	assert.True(t, ok)
	assert.Equal(t, "192.168.0.1", ip)
}

func TestGetRegion(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyHead, &testHead{
		Region: "asia",
	})

	region, ok := GetRegion(ctx)
	assert.True(t, ok)
	assert.Equal(t, "asia", region)
}

func TestGetDt(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyHead, &testHead{
		Dt: int32(1),
	})

	dt, ok := GetDt(ctx)
	assert.True(t, ok)
	assert.Equal(t, int32(1), dt)
}

func TestGetUnionId(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyHead, &testHead{
		Unionid: "xyz",
	})

	unionId, ok := GetUnionId(ctx)
	assert.True(t, ok)
	assert.Equal(t, "xyz", unionId)
}

type testEmptyHead struct{}

func (m testEmptyHead) ToKV() map[string]interface{} {
	return map[string]interface{}{}
}

func TestGetterNegativeCases(t *testing.T) {
	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()

		ip, ok := GetIp(ctx)
		assert.False(t, ok)
		assert.Equal(t, "", ip)
	})

	t.Run("empty kv", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, ContextKeyHead, &testEmptyHead{})

		ip, ok := GetIp(ctx)
		assert.False(t, ok)
		assert.Equal(t, "", ip)
	})
}
