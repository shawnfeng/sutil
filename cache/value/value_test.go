package value

import (
	"context"
	//"fmt"
	"github.com/shawnfeng/sutil/slog"
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
	cache := NewCache("base/report", "test", 60*time.Second, load)

	var test Test
	err := cache.Get(ctx, 3.5, &test)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	slog.Infof("test: %v", test)

	//	cache.Del(ctx, 1)
	err = cache.Get(ctx, 1, &test)
	if err != nil {
		t.Errorf("get err: %v", err)
	}
	slog.Infof("test: %v", test)
}
