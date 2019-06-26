package scontext

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