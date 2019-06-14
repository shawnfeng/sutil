package slog

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	"strings"
)

var ErrorTraceIDNotFound = errors.New("traceID not found")
var ErrorHeadKVNotFound = errors.New("valid context head not found")

type contextKV map[string]interface{}

func newContextKV() contextKV {
	return contextKV{}
}

func (ckv contextKV) String() string {
	var parts []string
	for k, v := range ckv {
		parts = append(parts, fmt.Sprintf("%s:%v", k, v))
	}
	return strings.Join(parts, " ")
}

func extractTraceID(ctx context.Context) (error, contextKV) {
	ckv := newContextKV()
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			ckv["traceID"] = sc.TraceID()
			return nil, ckv
		}
	}
	return ErrorTraceIDNotFound, nil
}

func extractHead(ctx context.Context) (error, contextKV) {
	head := ctx.Value("Head")
	if chd, ok := head.(contextHeader); ok {
		return nil, contextKV(chd.toKV())
	}
	return ErrorHeadKVNotFound, nil
}

type contextHeader interface {
	toKV() map[string]interface{}
}

func extractContext(ctx context.Context, includeHead bool) []interface{} {
	var v []interface{}

	if err, ckv := extractTraceID(ctx); err == nil {
		v = append(v, ckv)
	}

	if includeHead {
		if err, ckv := extractHead(ctx); err == nil {
			v = append(v, ckv)
		}
	}
	return v
}

func extractContextAsString(ctx context.Context, includeHead bool) (s string) {
	var parts []string
	for _, kv := range extractContext(ctx, includeHead) {
		parts = append(parts, fmt.Sprint(kv))
	}
	return strings.Join(parts, " ")
}