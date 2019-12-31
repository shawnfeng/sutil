package trace

import (
	"context"
	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/slog/slog"
	"net/http"
	"strings"
	"sync"
)

const (
	ServiceName      = "trace"
	DefaultNameSpace = "application"
)

const FilterUrls = "span_filter_urls"

const ListConfigSep = ","

var apolloCenter center.ConfigCenter

var apolloSpanFilterConfig *spanFilterConfig

func init() {
	fun := "trace.init --> "
	ctx := context.Background()

	var err error
	apolloCenter, err = center.NewConfigCenter(center.ApolloConfigCenter)
	if err != nil {
		slog.Errorf(ctx, "%s new config center error, center type: %d, err: %s", fun, center.ApolloConfigCenter, err.Error())
		return
	}

	err = apolloCenter.Init(ctx, ServiceName, nil)
	if err != nil {
		slog.Errorf(ctx, "%s init apollo config center error, service name: %s, namespaces: %s", fun, ServiceName, "")
		return
	}

	urls, ok := apolloCenter.GetString(ctx, FilterUrls)
	if !ok {
		slog.Errorf(ctx, "%s get %s from apollo failed", fun, FilterUrls)
		return
	}
	slog.Infof(ctx, "%s get %s from apollo: %s", fun, FilterUrls, urls)

	urlList := strings.Split(urls, ListConfigSep)

	apolloSpanFilterConfig = &spanFilterConfig{
		urls: urlList,
	}

	apolloCenter.StartWatchUpdate(ctx)
	apolloCenter.RegisterObserver(ctx, apolloSpanFilterConfig)
}

type spanFilterConfig struct {
	mu sync.RWMutex

	urls []string
}

func (m *spanFilterConfig) updateUrls(urls []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.urls = urls
}

func (m *spanFilterConfig) filterUrl(url string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.urls {
		if u == url {
			return false
		}
	}

	return true
}

func (m *spanFilterConfig) HandleChangeEvent(event *center.ChangeEvent) {
	fun := "spanFilterConfig.HandleChangeEvent --> "
	ctx := context.Background()

	if event.Namespace != DefaultNameSpace {
		return
	}

	for key, change := range event.Changes {
		if key == FilterUrls {
			slog.Infof(ctx, "%s get key %s from apollo, old value: %s, new value: %s", fun, key, change.OldValue, change.NewValue)
			urlList := strings.Split(change.NewValue, ListConfigSep)
			m.updateUrls(urlList)
		}
	}
}

func UrlSpanFilter(r *http.Request) bool {
	if apolloSpanFilterConfig != nil {
		return apolloSpanFilterConfig.filterUrl(r.URL.Path)
	}

	return true
}
