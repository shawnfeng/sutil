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

type simpleContextControlRouter struct {
	Group string
}

func (m *simpleContextControlRouter) GetControlRouteGroup() (string, bool) {
	return m.Group, true
}

func (m *simpleContextControlRouter) SetControlRouteGroup(group string) error {
	m.Group = group
	return nil
}

func TestGetRouteGroup(t *testing.T) {
	expectedGroup := "lane1"
	sc := &simpleContextControlRouter{Group: expectedGroup}
	ctx := context.WithValue(context.Background(), ContextKeyControl, sc)

	var group string
	var ok bool
	group, ok = GetControlRouteGroup(ctx)
	assert.True(t, ok)
	assert.Equal(t, expectedGroup, group)

	group, ok = GetControlRouteGroup(context.Background())
	assert.False(t, ok)
	assert.Equal(t, "", group)
}

func TestGetRouteGroupWithDefault(t *testing.T) {
	groupLane1 := "lane1"
	groupLane2 := "lane2"
	sc := &simpleContextControlRouter{Group: groupLane2}
	ctx := context.WithValue(context.Background(), ContextKeyControl, sc)

	assert.Equal(t, groupLane1, GetControlRouteGroupWithDefault(context.Background(), groupLane1))
	assert.Equal(t, groupLane2, GetControlRouteGroupWithDefault(ctx, groupLane1))
}

type simpleContextControlCaller struct {
	ServerName string
	ServerId   string
	Method     string
}

func (m *simpleContextControlCaller) GetControlCallerServerName() (string, bool) {
	return m.ServerName, true
}

func (m *simpleContextControlCaller) SetControlCallerServerName(serverName string) error {
	m.ServerName = serverName
	return nil
}

func (m *simpleContextControlCaller) GetControlCallerServerId() (string, bool) {
	return m.ServerId, true
}

func (m *simpleContextControlCaller) SetControlCallerServerId(serverId string) error {
	m.ServerId = serverId
	return nil
}

func (m *simpleContextControlCaller) GetControlCallerMethod() (string, bool) {
	return m.Method, true
}

func (m *simpleContextControlCaller) SetControlCallerMethod(method string) error {
	m.Method = method
	return nil
}

func TestContextControlCaller(t *testing.T) {
	simpleCaller := &simpleContextControlCaller{"report", "0", "GetRtcRuntimeLog"}
	var caller ContextControlCaller
	caller = simpleCaller

	serverName, ok := caller.GetControlCallerServerName()
	assert.True(t, ok)
	assert.Equal(t, "report", serverName)

	serverId, ok := caller.GetControlCallerServerId()
	assert.True(t, ok)
	assert.Equal(t, "0", serverId)

	method, ok := caller.GetControlCallerMethod()
	assert.True(t, ok)
	assert.Equal(t, "GetRtcRuntimeLog", method)
}
