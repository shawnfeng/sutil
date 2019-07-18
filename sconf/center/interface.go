package center

import (
	"context"
	"fmt"
)

type ConfigCenterType int

const (
	ApolloConfigCenter ConfigCenterType = iota
)

func (c ConfigCenterType) String() string {
	switch c {
	case ApolloConfigCenter:
		return "apollo"
	default:
		return "unknown"
	}
}

func NewConfigCenter(t ConfigCenterType) (ConfigCenter, error) {
	fun := "sconfcenter.NewConfigCenter-->"
	switch t {
	case ApolloConfigCenter:
		return newApolloConfigCenter(), nil
	default:
		return nil, fmt.Errorf("%s unsupported config center type:%v", fun, t)
	}
}

var defaultConfigCenter ConfigCenter

type ConfigObserver interface {
	HandleChangeEvent(event *ChangeEvent)
}

type ConfigCenter interface {
	Init(ctx context.Context, serviceName string, namespaceNames []string) error
	Stop(ctx context.Context) error

	SubscribeNamespaces(ctx context.Context, namespaceNames []string) error

	GetString(ctx context.Context, key string) (string, bool)
	GetStringWithNamespace(ctx context.Context, namespace, key string) (string, bool)
	GetBool(ctx context.Context, key string) (bool, bool)
	GetBoolWithNamespace(ctx context.Context, namespace, key string) (bool, bool)
	GetInt(ctx context.Context, key string) (int, bool)
	GetIntWithNamespace(ctx context.Context, namespace, key string) (int, bool)
	GetAllKeys(ctx context.Context) []string
	GetAllKeysWithNamespace(ctx context.Context, namespace string) []string

	Unmarshal(ctx context.Context, v interface{}) error
	UnmarshalWithNamespace(ctx context.Context, namespace string, v interface{}) error

	StartWatchUpdate(ctx context.Context)
	RegisterObserver(ctx context.Context, observer ConfigObserver) (recall func())
}

func Init(ctx context.Context, serviceName string, namespaceNames []string) error {
	defaultConfigCenter = newApolloConfigCenter()
	return defaultConfigCenter.Init(ctx, serviceName, namespaceNames)
}

func Stop(ctx context.Context) error {
	return defaultConfigCenter.Stop(ctx)
}

func SubscribeNamespaces(ctx context.Context, namespaceNames []string) error {
	return defaultConfigCenter.SubscribeNamespaces(ctx, namespaceNames)
}

func GetString(ctx context.Context, key string) (string, bool) {
	return defaultConfigCenter.GetString(ctx, key)
}

func GetStringWithNamespace(ctx context.Context, namespace, key string) (string, bool) {
	return defaultConfigCenter.GetStringWithNamespace(ctx, namespace, key)
}

func GetBool(ctx context.Context, key string) (bool, bool) {
	return defaultConfigCenter.GetBool(ctx, key)
}

func GetBoolWithNamespace(ctx context.Context, namespace, key string) (bool, bool) {
	return defaultConfigCenter.GetBoolWithNamespace(ctx, namespace, key)
}

func GetInt(ctx context.Context, key string) (int, bool) {
	return defaultConfigCenter.GetInt(ctx, key)
}

func GetIntWithNamespace(ctx context.Context, namespace, key string) (int, bool) {
	return defaultConfigCenter.GetIntWithNamespace(ctx, namespace, key)
}

func StartWatchUpdate(ctx context.Context) {
	defaultConfigCenter.StartWatchUpdate(ctx)
}

func RegisterObserver(ctx context.Context, observer ConfigObserver) (recall func()) {
	return defaultConfigCenter.RegisterObserver(ctx, observer)
}

func Unmarshall(ctx context.Context, v interface{}) error {
	return defaultConfigCenter.Unmarshal(ctx, v)
}

func UnmarshalWithNamespace(ctx context.Context, namespace string, v interface{}) error {
	return defaultConfigCenter.UnmarshalWithNamespace(ctx, namespace, v)
}
