package scontext

import (
	"context"
)

// 由于请求的上下文信息的 thrift 定义在 util 项目中，本模块主要为了避免循环依赖
const (
	ContextKeyTraceID = "traceID"
	ContextKeyControl = "Control"

	ContextKeyHead        = "Head"
	ContextKeyHeadUid     = "uid"
	ContextKeyHeadSource  = "source"
	ContextKeyHeadIp      = "ip"
	ContextKeyHeadRegion  = "region"
	ContextKeyHeadDt      = "dt"
	ContextKeyHeadUnionId = "unionid"
)

type ContextHeader interface {
	ToKV() map[string]interface{}
}

type ContextController interface {
	GetGroup() string
}

func GetGroup(ctx context.Context) string {
	value := ctx.Value(ContextKeyControl)
	if value == nil {
		return ""
	}

	control, ok := value.(ContextController)
	if ok == false {
		return ""
	}

	return control.GetGroup()
}

func GetGroupWithDefault(ctx context.Context, dv string) string {
	if group := GetGroup(ctx); group != "" {
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