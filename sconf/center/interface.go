package center

import "context"

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
	defaultConfigCenter = newApolloConfigCenter()
	return defaultConfigCenter.Init(ctx, serviceName, namespaceNames)
}

func Stop(ctx context.Context) error {
	return defaultConfigCenter.Stop(ctx)
}

func GetString(ctx context.Context, key string) string {
	return defaultConfigCenter.GetString(ctx, key)
}

func GetStringWithNamespace(ctx context.Context, namespace, key string) string {
	return defaultConfigCenter.GetStringWithNamespace(ctx, namespace, key)
}

func GetBool(ctx context.Context, key string) bool {
	return defaultConfigCenter.GetBool(ctx, key)
}

func GetBoolWithNamespace(ctx context.Context, namespace, key string) bool {
	return defaultConfigCenter.GetBoolWithNamespace(ctx, namespace, key)
}

func GetInt(ctx context.Context, key string) int {
	return defaultConfigCenter.GetInt(ctx, key)
}

func GetIntWithNamespace(ctx context.Context, namespace, key string) int {
	return defaultConfigCenter.GetIntWithNamespace(ctx, namespace, key)
}

func WatchUpdate(ctx context.Context) <-chan *ChangeEvent {
	return defaultConfigCenter.WatchUpdate(ctx)
}
