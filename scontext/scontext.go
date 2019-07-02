package scontext

import (
	"context"
	"github.com/shawnfeng/sutil/slog"
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
	fun := "GetGroup -->"
	value := ctx.Value(ContextKeyControl)
	if value == nil {
		slog.Infof("%s value is nil", fun)
		return ""
	}

	control, ok := value.(ContextController)
	if ok == false {
		slog.Infof("%s value.(ContextController) is false", fun)
		return ""
	}

	group := control.GetGroup()
	slog.Infof("%s group: %s", fun, group)

	return group
}
