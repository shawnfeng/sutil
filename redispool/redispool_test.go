package redispool

import (
	"testing"
	"github.com/fzzy/radix/redis"
	"github.com/shawnfeng/sutil/slog"
)

func TestLuaLoad(t *testing.T) {
	pool := NewRedisPool()

	err := pool.LoadLuaFile("Test", "./test.luad")
	slog.Infoln(err)
	if err == nil || err.Error() != "open ./test.luad: no such file or directory" {
		t.Errorf("error here")
	}

	err = pool.LoadLuaFile("Test", "./test.lua")
	if err != nil {
		t.Errorf("error here")
	}

	addr := "localhost:9600"
    args := []interface{}{
		2,
		"key1",
		"key2",
		"argv1",
		"argv2",
	}

	rp := pool.EvalSingle(addr, "Nothave", args)

	slog.Infoln(rp)

	if "get lua sha1 add:localhost:9600 key:Nothave err:lua not find" != rp.String() {
		t.Errorf("error here")
	}

	rp = pool.EvalSingle(addr, "Test", args)

	slog.Infoln(rp)
	if rp.Type == redis.ErrorReply {
		t.Errorf("error here")
	}

	if rp.String() != "key1key2argv1argv222" {
		t.Errorf("error here")
	}


}
