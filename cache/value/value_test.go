package value

import (
	"context"
	cache2 "github.com/shawnfeng/sutil/cache"
	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/trace"

	//"fmt"
	"github.com/shawnfeng/sutil/slog/slog"
	"testing"
	"time"
)

type Test struct {
	Id int64
}

func load(key interface{}) (value interface{}, err error) {

	//return nil, fmt.Errorf("not found")
	return &Test{
		Id: 1,
	}, nil
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	_ = trace.InitDefaultTracer("cache.test")
	_ = center.Init(ctx, "base/report", []string{"application", "infra.cache"})

	cache := NewCache("base/report", "test", 60*time.Second, load)
	_ = SetConfiger(ctx, cache2.ConfigerTypeApollo)

	var test Test
	err := cache.Get(ctx, 7, &test)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	slog.Infof(ctx, "test: %v", test)

	cache.Del(ctx, 7)
	err = cache.Get(ctx, 7, &test)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	slog.Infof(ctx, "test: %v", test)

	time.Sleep(2 * time.Second)
}
