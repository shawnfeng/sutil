package center

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

var defaultConfigCenter ConfigCenter

type ConfigCenter interface {
	Init(ctx context.Context, serviceName string, namespaceNames []string) error
	Stop(ctx context.Context) error

	GetString(ctx context.Context, key string) string
	GetStringWithNamespace(ctx context.Context, namespace, key string) string
	GetBool(ctx context.Context, key string) bool
	GetBoolWithNamespace(ctx context.Context, namespace, key string) bool
	GetInt(ctx context.Context, key string) int
	GetIntWithNamespace(ctx context.Context, namespace, key string) int

	WatchUpdate(ctx context.Context) <-chan *ChangeEvent
}

func Init(ctx context.Context, serviceName string, namespaceNames []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.Init")
	defer span.Finish()

	defaultConfigCenter = newApolloConfigCenter()
	return defaultConfigCenter.Init(ctx, serviceName, namespaceNames)
}

func Stop(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.Stop")
	defer span.Finish()

	return defaultConfigCenter.Stop(ctx)
}

func GetString(ctx context.Context, key string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetString")
	defer span.Finish()

	return defaultConfigCenter.GetString(ctx, key)
}

func GetStringWithNamespace(ctx context.Context, namespace, key string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetStringWithNamespace")
	defer span.Finish()

	return defaultConfigCenter.GetStringWithNamespace(ctx, namespace, key)
}

func GetBool(ctx context.Context, key string) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetBool")
	defer span.Finish()

	return defaultConfigCenter.GetBool(ctx, key)
}

func GetBoolWithNamespace(ctx context.Context, namespace, key string) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetBoolWithNamespace")
	defer span.Finish()

	return defaultConfigCenter.GetBoolWithNamespace(ctx, namespace, key)
}

func GetInt(ctx context.Context, key string) int {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetInt")
	defer span.Finish()

	return defaultConfigCenter.GetInt(ctx, key)
}

func GetIntWithNamespace(ctx context.Context, namespace, key string) int {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetIntWithNamespace")
	defer span.Finish()

	return defaultConfigCenter.GetIntWithNamespace(ctx, namespace, key)
}

func WatchUpdate(ctx context.Context) <-chan *ChangeEvent {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.WatchUpdate")
	defer span.Finish()

	return defaultConfigCenter.WatchUpdate(ctx)
}
