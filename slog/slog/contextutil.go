package slog

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/uber/jaeger-client-go"
	"strings"
)

var (
	emptyTrace = contextKV{scontext.ContextKeyTraceID: jaeger.TraceID{0, 0}}
	emptyHead  = contextKV{scontext.ContextKeyHeadUid: int64(0)}
)

var errorTraceIDNotFound = errors.New("traceID not found")
var errorHeadKVNotFound = errors.New("valid context head not found")

type contextKV map[string]interface{}

func newContextKV() contextKV {
	return contextKV{}
}

func (ckv contextKV) String() string {
	if v, ok := ckv[scontext.ContextKeyTraceID]; ok {
		return fmt.Sprintf("%v", v)
	}

	var parts []string
	if v, ok := ckv[scontext.ContextKeyHeadUid]; ok {
		if uid, uok := v.(int64); uok {
			parts = append(parts, fmt.Sprintf("%d", uid))
		}
	}

	var restParts []string
	for k, v := range ckv {
		if k != scontext.ContextKeyHeadUid && k != scontext.ContextKeyTraceID {
			restParts = append(restParts, fmt.Sprintf("%s:%v", k, v))
		}
	}
	if len(restParts) > 0 {
		parts = append(parts, strings.Join(restParts, " "))
		return strings.Join(parts, "\t")
	} else {
		return strings.Join(parts, "\t") + "\t"
	}
}

func extractTraceID(ctx context.Context) (error, contextKV) {
	ckv := newContextKV()
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			ckv[scontext.ContextKeyTraceID] = sc.TraceID()
			return nil, ckv
		}
	}
	return errorTraceIDNotFound, nil
}

func extractHead(ctx context.Context, fullHead bool) (error, contextKV) {
	head := ctx.Value(scontext.ContextKeyHead)
	if chd, ok := head.(scontext.ContextHeader); ok {
		kv := chd.ToKV()
		if fullHead {
			return nil, contextKV(chd.ToKV())
		}
		return nil, contextKV(map[string]interface{}{scontext.ContextKeyHeadUid: kv[scontext.ContextKeyHeadUid]})
	}
	return errorHeadKVNotFound, nil
}

func extractContext(ctx context.Context, fullHead bool) (v []interface{}) {
	if ctx == nil {
		return
	}

	if err, ckv := extractTraceID(ctx); err == nil {
		v = append(v, ckv)
	} else {
		v = append(v, emptyTrace)
	}

	if err, ckv := extractHead(ctx, fullHead); err == nil {
		v = append(v, ckv)
	} else {
		v = append(v, emptyHead)
	}

	return
}

func extractContextAsString(ctx context.Context, fullHead bool) (s string) {
	var parts []string
	for _, kv := range extractContext(ctx, fullHead) {
		parts = append(parts, fmt.Sprint(kv))
	}
	return strings.Join(parts, "\t")
}
