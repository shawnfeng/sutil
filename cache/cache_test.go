package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Test struct {
	ID  int `json:"id"`
}

func (p *Test) Marshal() ([]byte, error){
	return json.Marshal(&p)
}

func (p *Test) Unmarshal(data []byte) error{
	return json.Unmarshal(data, p)
}

// cache miss load
func (p *Test) Load(key string) error{
	p.ID = 1
	return nil
}

func TestNewCommonCache(t *testing.T) {
	id := 3
	c, err := NewCommonCache("base/changeboard", "test", 10, 60)
	if err != nil {
		t.Errorf("NewCommonCache err: %s", err.Error())
	}
	var test Test
	test.ID = id
	err = c.Set("test", &test)
	if err != nil {
		t.Errorf("Set err: %s", err.Error())
	}
	err = c.Get("test", &test)
	if err != nil {
		t.Errorf("Set err: %s", err.Error())
	}
	assert.Equal(t, id, test.ID)
}

func TestCache_Set(t *testing.T) {
	ctx := context.Background()
	c, err  := NewCacheByNamespace(ctx, "test/test","test",60)
	if err != nil {
		t.Errorf("NewCacheByNamespace err: %s", err.Error())
	}
	WatchUpdate(ctx)
	var test Test
	test.ID = 3
	err = c.Set("test", &test)
	if err != nil {
		t.Errorf("Set err: %s", err.Error())
	}
	t.Log(test)
}

func TestCache_Get(t *testing.T) {
	ctx := context.Background()
	c, err  := NewCacheByNamespace(ctx, "base/test","test",60)
	if err != nil {
		t.Errorf("NewCacheByNamespace err: %s", err.Error())
	}
	//WatchUpdate(ctx)
	var test Test
	for {
		err = c.Get("test", &test)
		if err != nil {
			t.Errorf("Get err: %s", err.Error())
			return
		}
		fmt.Println(test)
		time.Sleep(time.Second)
	}
	//t.Log(test)
	//c.Del("test")
	//// load
	//err = c.Get("test", &test)
	//if err != nil {
	//	t.Errorf("Get err: %s", err.Error())
	//}
	//assert.Equal(t, 1, test.ID)
}
