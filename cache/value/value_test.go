package value

import (
	"context"
	"github.com/shawnfeng/sutil/cache"
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

	c := NewCache("base/report", "test", 60*time.Second, load)
	_ = SetConfiger(ctx, cache.ConfigerTypeApollo)
	WatchUpdate(ctx)

	var test Test
	err := c.Get(ctx, 7, &test)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	slog.Infof(ctx, "test: %v", test)

	c.Del(ctx, 7)
	c.Load(ctx, 7)
	err = c.Get(ctx, 7, &test)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	slog.Infof(ctx, "test: %v", test)

	time.Sleep(2 * time.Second)
}
