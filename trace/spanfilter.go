package trace

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/slog/slog"
	"net/http"
	"strings"
	"sync"
)

const FilterUrls = "span_filter_urls"

const ListConfigSep = ","

var once sync.Once

var apolloCenter center.ConfigCenter

var apolloSpanFilterConfig *spanFilterConfig

func InitTraceSpanFilter() error {
	fun := "TraceSpanFilter.init --> "
	ctx := context.Background()

	initApolloCenter(ctx)
	if apolloCenter == nil {
		return fmt.Errorf("trace apollo center is nil")
	}

	urls, ok := apolloCenter.GetStringWithNamespace(ctx, center.DefaultApolloTraceNamespace, FilterUrls)
	if !ok {
		return fmt.Errorf("not get %s from apollo namespace %s", FilterUrls, center.DefaultApolloTraceNamespace)
	}
	slog.Infof(ctx, "%s get %s from apollo: %s", fun, FilterUrls, urls)

	urlList := strings.Split(urls, ListConfigSep)

	apolloSpanFilterConfig = &spanFilterConfig{
		urls: urlList,
	}

	apolloCenter.RegisterObserver(ctx, apolloSpanFilterConfig)
	return nil
}

func initApolloCenter(ctx context.Context) {
	fun := "ApolloCenter.init --> "

	if apolloCenter != nil {
		return
	}

	once.Do(func() {
		var err error
		apolloCenter, err = center.NewConfigCenter(center.ApolloConfigCenter)
		if err != nil {
			slog.Errorf(ctx, "%s new config center error, center type: %d, err: %s", fun, center.ApolloConfigCenter, err.Error())
			return
		}

		err = apolloCenter.Init(ctx, center.DefaultApolloMiddlewareService, []string{center.DefaultApolloTraceNamespace})
		if err != nil {
			slog.Errorf(ctx, "%s init apollo config center error, service name: %s, namespaces: %s, err: %s",
				fun, center.DefaultApolloMiddlewareService, center.DefaultApolloTraceNamespace, err.Error())
			return
		}

		apolloCenter.StartWatchUpdate(ctx)
	})
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

	if event.Namespace != center.DefaultApolloTraceNamespace {
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
