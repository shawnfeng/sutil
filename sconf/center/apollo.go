package center

import (
	"context"
	"github.com/ZhengHe-MD/agollo/v4"
	"github.com/ZhengHe-MD/properties"
	"github.com/opentracing/opentracing-go"
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
	defaultChangeEventSize      = 32
)

type apolloConfigCenter struct {
	conf            *agollo.Conf
	ag              *agollo.Agollo
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

func (ap *apolloConfigCenter) Init(ctx context.Context, serviceName string, namespaceNames []string) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.Init")
	defer span.Finish()

	fun := "apolloConfigCenter.Init-->"

	agollo.SetLogger(slog.GetLogger())

	conf := confFromEnv()
	conf.AppID = normalizeServiceName(serviceName)

	if len(namespaceNames) > 0 {
		conf.NameSpaceNames = namespaceNames
	} else {
		conf.NameSpaceNames = []string{defaultNamespaceApplication}
	}

	for i, namespace := range conf.NameSpaceNames {
		conf.NameSpaceNames[i] = normalizeServiceName(namespace)
	}

	ap.conf = conf
	ap.ag = agollo.NewAgollo(conf)

	slog.Infof(ctx, "%s start agollo with conf:%v", fun, ap.conf)

	if err = ap.ag.Start(); err != nil {
		slog.Errorf(ctx, "%s agollo starts err:%v", fun, err)
	} else {
		slog.Infof(ctx, "%s agollo starts succeed:%v", fun, err)
	}

	return
}

func (ap *apolloConfigCenter) Stop(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.Stop")
	defer span.Finish()
	return ap.ag.Stop()
}

func (ap *apolloConfigCenter) SubscribeNamespaces(ctx context.Context, namespaceNames []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.SubscribeNamespaces")
	defer span.Finish()
	return ap.ag.SubscribeToNamespaces(namespaceNames...)
}

func (ap *apolloConfigCenter) GetString(ctx context.Context, key string) (string, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetString")
	defer span.Finish()
	return ap.ag.GetString(key)
}

func (ap *apolloConfigCenter) GetStringWithNamespace(ctx context.Context, namespace, key string) (string, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetStringWithNamespace")
	defer span.Finish()

	return ap.ag.GetStringWithNamespace(namespace, key)
}

func (ap *apolloConfigCenter) GetBool(ctx context.Context, key string) (bool, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetBool")
	defer span.Finish()

	return ap.ag.GetBool(key)
}

func (ap *apolloConfigCenter) GetBoolWithNamespace(ctx context.Context, namespace, key string) (bool, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetBoolWithNamespace")
	defer span.Finish()

	return ap.ag.GetBoolWithNamespace(namespace, key)
}

func (ap *apolloConfigCenter) GetInt(ctx context.Context, key string) (int, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetInt")
	defer span.Finish()

	return ap.ag.GetInt(key)
}

func (ap *apolloConfigCenter) GetIntWithNamespace(ctx context.Context, namespace, key string) (int, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetIntWithNamespace")
	defer span.Finish()

	return ap.ag.GetIntWithNamespace(namespace, key)
}

func (ap *apolloConfigCenter) GetAllKeys(ctx context.Context) []string {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetAllKeys")
	defer span.Finish()

	return ap.ag.GetAllKeys("application")
}

func (ap *apolloConfigCenter) GetAllKeysWithNamespace(ctx context.Context, namespace string) []string {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetAllKeysWithNamespace")
	defer span.Finish()

	return ap.ag.GetAllKeys(namespace)
}

func (ap *apolloConfigCenter) StartWatchUpdate(ctx context.Context) {
	ap.ag.StartWatchUpdate()
}

type agolloObserver struct {
	observer ConfigObserver
}

func (o *agolloObserver) HandleChangeEvent(ce *agollo.ChangeEvent) {
	o.observer.HandleChangeEvent(fromAgolloChangeEvent(ce))
}

func (ap *apolloConfigCenter) RegisterObserver(ctx context.Context, observer ConfigObserver) func() {
	return ap.ag.RegisterObserver(&agolloObserver{observer})
}

func (ap *apolloConfigCenter) Unmarshal(ctx context.Context, v interface{}) error {
	return ap.UnmarshalWithNamespace(ctx, defaultNamespaceApplication, v)
}

func (ap *apolloConfigCenter) UnmarshalWithNamespace(ctx context.Context, namespace string, v interface{}) error {
	var kv = map[string]string{}

	ks := ap.GetAllKeysWithNamespace(ctx, namespace)
	for _, k := range ks {
		if v, ok := ap.GetStringWithNamespace(ctx, namespace, k); ok {
			kv[k] = v
		}
	}

	return properties.UnmarshalKV(kv, v)
}

func (ap *apolloConfigCenter) UnmarshalKey(ctx context.Context, key string, v interface{}) error {
	return ap.UnmarshalKeyWithNamespace(ctx, defaultNamespaceApplication, key, v)
}

func (ap *apolloConfigCenter) UnmarshalKeyWithNamespace(ctx context.Context, namespace string, key string, v interface{}) error {
	var kv = map[string]string{}

	ks := ap.GetAllKeysWithNamespace(ctx, namespace)
	for _, k := range ks {
		if v, ok := ap.GetStringWithNamespace(ctx, namespace, k); ok {
			kv[k] = v
		}
	}

	bs, err := properties.Marshal(&kv)
	if err != nil {
		return err
	}

	return properties.UnmarshalKey(key, bs, v)
}
