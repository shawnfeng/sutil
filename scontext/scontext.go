package scontext

import (
	"context"
	"errors"
)

// 由于请求的上下文信息的 thrift 定义在 util 项目中，本模块主要为了避免循环依赖
const (
	ContextKeyTraceID = "traceID"

	ContextKeyHead        = "Head"
	ContextKeyHeadUid     = "uid"
	ContextKeyHeadSource  = "source"
	ContextKeyHeadIp      = "ip"
	ContextKeyHeadRegion  = "region"
	ContextKeyHeadDt      = "dt"
	ContextKeyHeadUnionId = "unionid"

	ContextKeyControl       = "Control"
)

const DefaultGroup = ""

var ErrInvalidContext = errors.New("invalid context")

type ContextHeader interface {
	ToKV() map[string]interface{}
}

type ContextControlRouter interface {
	GetControlRouteGroup() (string, bool)
	SetControlRouteGroup(string) error
}

type ContextControlCaller interface {
	GetControlCallerServerName() (string, bool)
	SetControlCallerServerName(string) error
	GetControlCallerServerId() (string, bool)
	SetControlCallerServerId(string) error
	GetControlCallerMethod() (string, bool)
	SetControlCallerMethod(string) error
}

func GetControlRouteGroup(ctx context.Context) (group string, ok bool) {
	value := ctx.Value(ContextKeyControl)
	if value == nil {
		ok = false
		return
	}
	control, ok := value.(ContextControlRouter)
	if ok == false {
		return
	}
	return control.GetControlRouteGroup()
}

func SetControlRouteGroup(ctx context.Context, group string) (context.Context, error) {
	value := ctx.Value(ContextKeyControl)
	if value == nil {
		return ctx, ErrInvalidContext
	}
	control, ok := value.(ContextControlRouter)
	if !ok {
		return ctx, ErrInvalidContext
	}

	err := control.SetControlRouteGroup(group)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, ContextKeyControl, control), nil
}

func GetControlRouteGroupWithDefault(ctx context.Context, dv string) string {
	if group, ok := GetControlRouteGroup(ctx); ok {
		return group
	}
	return dv
}

func getHeaderByKey(ctx context.Context, key string) (val interface{}, ok bool) {
	head := ctx.Value(ContextKeyHead)
	if head == nil {
		ok = false
		return
	}

	var header ContextHeader
	if header, ok = head.(ContextHeader); ok {
		val, ok = header.ToKV()[key]
	}
	return
}

func GetUid(ctx context.Context) (uid int64, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadUid)
	if ok {
		uid, ok = val.(int64)
	}
	return
}

func GetSource(ctx context.Context) (source int32, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadSource)
	if ok {
		source, ok = val.(int32)
	}
	return
}

func GetIp(ctx context.Context) (ip string, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadIp)
	if ok {
		ip, ok = val.(string)
	}
	return
}

func GetRegion(ctx context.Context) (region string, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadRegion)
	if ok {
		region, ok = val.(string)
	}
	return
}

func GetDt(ctx context.Context) (dt int32, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadDt)
	if ok {
		dt, ok = val.(int32)
	}
	return
}

func GetUnionId(ctx context.Context) (unionId string, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadUnionId)
	if ok {
		unionId, ok = val.(string)
	}
	return
}

func getControlCaller(ctx context.Context) (ContextControlCaller, error) {
	value := ctx.Value(ContextKeyControl)
	if value == nil {
		return nil, ErrInvalidContext
	}
	caller, ok := value.(ContextControlCaller)
	if !ok {
		return nil, ErrInvalidContext
	}
	return caller, nil
}

func GetControlCallerServerName(ctx context.Context) (serverName string, ok bool) {
	caller, ok := ctx.Value(ContextKeyControl).(ContextControlCaller)
	if !ok {
		return
	}
	return caller.GetControlCallerServerName()
}

func SetControlCallerServerName(ctx context.Context, serverName string) (context.Context, error) {
	caller, err := getControlCaller(ctx)
	if err != nil {
		return ctx, err
	}
	err = caller.SetControlCallerServerName(serverName)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, ContextKeyControl, caller), nil
}

func GetControlCallerServerId(ctx context.Context) (serverId string, ok bool) {
	caller, ok := ctx.Value(ContextKeyControl).(ContextControlCaller)
	if !ok {
		return
	}
	return caller.GetControlCallerServerId()
}

func SetControlCallerServerId(ctx context.Context, serverId string) (context.Context, error) {
	caller, err := getControlCaller(ctx)
	if err != nil {
		return ctx, err
	}
	err = caller.SetControlCallerServerId(serverId)
	return context.WithValue(ctx, ContextKeyControl, caller), nil
}

func GetControlCallerMethod(ctx context.Context) (method string, ok bool) {
	caller, ok := ctx.Value(ContextKeyControl).(ContextControlCaller)
	if !ok {
		return
	}
	return caller.GetControlCallerMethod()
}

func SetControlCallerMethod(ctx context.Context, method string) (context.Context, error) {
	caller, err := getControlCaller(ctx)
	if err != nil {
		return ctx, err
	}
	err = caller.SetControlCallerMethod(method)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, ContextKeyControl, caller), nil
}

