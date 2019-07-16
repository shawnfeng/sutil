package value

import (
	"context"
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
	_ = trace.InitDefaultTracer("cache.test")

	ctx := context.Background()
	cache := NewCache("base/report", "test", 60*time.Second, load)

	var test Test
	err := cache.Get(ctx, 3, &test)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	slog.Infof(ctx, "test: %v", test)

	cache.Del(ctx, 3)
	err = cache.Get(ctx, 3, &test)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	slog.Infof(ctx, "test: %v", test)
}
