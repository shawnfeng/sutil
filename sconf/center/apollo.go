package center

import (
	"context"
	"github.com/ZhengHe-MD/agollo"
	"github.com/shawnfeng/sutil/slog/slog"
	"os"
	"strings"
	"sync"
)

const (
	envApolloCluster            = "APOLLO_CLUSTER"
	envApolloHostPort           = "APOLLO_HOST_PORT"
	defaultCluster              = "default"
	defaultHostPort             = "apollo-meta.ibanyu.com:30002"
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

// NOTE: apollo 不支持在项目名称中使用 '/'，因此规定用 '.' 代替 '/'
//       base/authapi => base.authapi
func normalizeServiceName(serviceName string) string {
	return strings.Replace(serviceName, "/", ".", -1)
}

func (ap *apolloConfigCenter) Init(ctx context.Context, serviceName string, namespaceNames []string) error {
	fun := "apollo.Init-->"

	agollo.SetLogger(slog.GetLogger())

	conf := confFromEnv()
	conf.AppID = normalizeServiceName(serviceName)

	if len(namespaceNames) > 0 {
		conf.NameSpaceNames = namespaceNames
	} else {
		conf.NameSpaceNames = []string{defaultNamespaceApplication}
	}

	ap.conf = conf

	slog.Infof(ctx, "%s start agollo with conf:%v", fun, ap.conf)

	if err := agollo.StartWithConf(ap.conf); err != nil {
		slog.Errorf(ctx, "%s agollo starts err:%v", fun, err)
	} else {
		slog.Infof(ctx, "%s agollo starts succeed:%v", fun, err)
	}

	return nil
}

func (ap *apolloConfigCenter) Stop(ctx context.Context) error {
	return agollo.Stop()
}

func (ap *apolloConfigCenter) GetString(ctx context.Context, key string) string {
	return agollo.GetString(key, "")
}

func (ap *apolloConfigCenter) GetStringWithNamespace(ctx context.Context, namespace, key string) string {
	return agollo.GetStringWithNamespace(namespace, key, "")
}

func (ap *apolloConfigCenter) GetBool(ctx context.Context, key string) bool {
	return agollo.GetBool(key, false)
}

func (ap *apolloConfigCenter) GetBoolWithNamespace(ctx context.Context, namespace, key string) bool {
	return agollo.GetBoolWithNamespace(namespace, key, false)
}

func (ap *apolloConfigCenter) GetInt(ctx context.Context, key string) int {
	return agollo.GetInt(key, 0)
}

func (ap *apolloConfigCenter) GetIntWithNamespace(ctx context.Context, namespace, key string) int {
	return agollo.GetIntWithNamespace(namespace, key, 0)
}

func (ap *apolloConfigCenter) StartWatchUpdate(ctx context.Context) {
	agollo.StartWatchUpdate()
}

type agolloObserver struct {
	observer ConfigObserver
}

func (o *agolloObserver) HandleChangeEvent(ce *agollo.ChangeEvent) {
	o.observer.HandleChangeEvent(fromAgolloChangeEvent(ce))
}

func (ap *apolloConfigCenter) RegisterObserver(ctx context.Context, observer ConfigObserver) func() {
	return agollo.RegisterObserver(&agolloObserver{observer})
}
