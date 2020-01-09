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

var initApolloLock sync.Mutex

var apolloCenter center.ConfigCenter

var apolloSpanFilterConfig *spanFilterConfig

func InitTraceSpanFilter() error {
	fun := "TraceSpanFilter.init --> "
	ctx := context.Background()

	if err := initApolloCenter(ctx); err != nil {
		return err
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

func initApolloCenter(ctx context.Context) error {
	if apolloCenter != nil {
		return nil
	}

	initApolloLock.Lock()
	defer initApolloLock.Unlock()

	if apolloCenter != nil {
		return nil
	}

	newCenter, err := center.NewConfigCenter(center.ApolloConfigCenter)
	if err != nil {
		return fmt.Errorf("new config center error, %s", err.Error())
	}

	namespaceList := []string{center.DefaultApolloTraceNamespace}
	err = newCenter.Init(ctx, center.DefaultApolloMiddlewareService, namespaceList)
	if err != nil {
		return fmt.Errorf("init apollo with service %s namespace %s error, %s",
			center.DefaultApolloMiddlewareService, strings.Join(namespaceList, " "), err.Error())
	}

	newCenter.StartWatchUpdate(ctx)

	apolloCenter = newCenter
	return nil
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
