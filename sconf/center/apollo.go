package center

import (
	"context"
	"github.com/ZhengHe-MD/agollo"
	"github.com/opentracing/opentracing-go"
	"github.com/shawnfeng/sutil/slog/slog"
	"os"
	"sync"
)

const (
	envApolloCluster  = "APOLLO_CLUSTER"
	envApolloHostPort = "APOLLO_HOST_PORT"

	defaultCluster = "default"
	// TODO: change it after hostname registered
	defaultHostPort             = "10.111.203.142:30002"
	defaultCacheDir             = "/tmp/sconfcenter"
	defaultNamespaceApplication = "application"

	defaultChangeEventSize = 32
)

type apolloConfigCenter struct {
	conf            *agollo.Conf
	watchUpdateOnce sync.Once
	changeEventChan chan *ChangeEvent
}

func newApolloConfigCenter() *apolloConfigCenter {
	return &apolloConfigCenter{
		changeEventChan: make(chan *ChangeEvent, defaultChangeEventSize),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func confFromEnv() *agollo.Conf {
	cluster := getEnvWithDefault(envApolloCluster, defaultCluster)
	hostport := getEnvWithDefault(envApolloHostPort, defaultHostPort)

	return &agollo.Conf{
		Cluster:  cluster,
		CacheDir: defaultCacheDir,
		IP:       hostport,
	}
}

func (ap *apolloConfigCenter) Init(ctx context.Context, serviceName string, namespaceNames []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.Init")
	defer span.Finish()

	conf := confFromEnv()
	conf.AppID = serviceName

	if len(namespaceNames) > 0 {
		conf.NameSpaceNames = namespaceNames
	} else {
		conf.NameSpaceNames = []string{defaultNamespaceApplication}
	}

	ap.conf = conf
	return agollo.StartWithConf(ap.conf)
}

func (ap *apolloConfigCenter) Stop(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.Stop")
	defer span.Finish()

	return agollo.Stop()
}

func (ap *apolloConfigCenter) GetString(ctx context.Context, key string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetString")
	defer span.Finish()

	return agollo.GetString(key, "")
}

func (ap *apolloConfigCenter) GetStringWithNamespace(ctx context.Context, namespace, key string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetStringWithNamespace")
	defer span.Finish()

	return agollo.GetStringWithNamespace(namespace, key, "")
}

func (ap *apolloConfigCenter) GetBool(ctx context.Context, key string) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetBool")
	defer span.Finish()

	return agollo.GetBool(key, false)
}

func (ap *apolloConfigCenter) GetBoolWithNamespace(ctx context.Context, namespace, key string) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetBoolWithNamespace")
	defer span.Finish()

	return agollo.GetBoolWithNamespace(namespace, key, false)
}

func (ap *apolloConfigCenter) GetInt(ctx context.Context, key string) int {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetInt")
	defer span.Finish()

	return agollo.GetInt(key, 0)
}

func (ap *apolloConfigCenter) GetIntWithNamespace(ctx context.Context, namespace, key string) int {
	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.GetIntWithNamespace")
	defer span.Finish()

	return agollo.GetIntWithNamespace(namespace, key, 0)
}

func (ap *apolloConfigCenter) WatchUpdate(ctx context.Context) <-chan *ChangeEvent {
	fun := "sconfcenter.WatchUpdate-->"

	span, _ := opentracing.StartSpanFromContext(ctx, "sconfcenter.WatchUpdate")
	defer span.Finish()

	ap.watchUpdateOnce.Do(func() {
		agolloChangeEventChan := agollo.WatchUpdate()
		go func() {
		WatchUpdateLoop:
			for {
				select {
				case <-ctx.Done():
					close(ap.changeEventChan)
					break WatchUpdateLoop
				case ace := <-agolloChangeEventChan:
					ce := fromAgolloChangeEvent(ace)
					slog.Infof(ctx, "%s receive change event:%v", fun, ce)
					ap.changeEventChan <- ce
				}
			}
		}()
	})

	return ap.changeEventChan
}
