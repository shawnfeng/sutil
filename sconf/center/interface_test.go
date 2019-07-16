package center

import (
	"context"
	"testing"
)

const (
	testService   = "base/authapi"
	testKey       = "key1"
	testNamespace = "db"
)

func TestGetTypedVar(t *testing.T) {
	ctx := context.Background()

	_ = Init(ctx, testService, []string{"application"})
	defer Stop(ctx)

	strVal, ok := GetString(ctx, testKey)
	if !ok {
		t.Errorf("strVal expect:%s to exist", testKey)
	}

	if strVal != "1" {
		t.Errorf("strVal expect:1 got:%v", strVal)
	}

	boolVal, ok := GetBool(ctx, testKey)

	if !ok {
		t.Errorf("boolVal expect:%s to exist", testKey)
	}

	if boolVal != true {
		t.Error("boolVal expect:true got:false")
	}

	intVal, ok := GetInt(ctx, testKey)

	if !ok {
		t.Errorf("intVal expect:%s to exist", testKey)
	}

	if intVal != 1 {
		t.Errorf("intVal expect:1 got:%d", intVal)
	}
}

func TestGetTypedVarWithNamespace(t *testing.T) {
	ctx := context.Background()

	_ = Init(ctx, testService, []string{"application", testNamespace})
	defer Stop(ctx)

	strVal, ok := GetStringWithNamespace(ctx, testNamespace, testKey)
	if ok {
		t.Errorf("strVal expect:%s to not exist", testKey)
	}

	if strVal != "" {
		t.Errorf("strVal expect:'' got:%v", strVal)
	}

	boolVal, ok := GetBoolWithNamespace(ctx, testNamespace, testKey)
	if boolVal != false {
		t.Error("boolVal expect:false got:true")
	}
	if ok {
		t.Errorf("boolVal expect:%s to not exist", testKey)
	}

	intVal, ok := GetIntWithNamespace(ctx, testNamespace, testKey)
	if ok {
		t.Errorf("intVal expect:%s to not exist", testKey)
	}
	if intVal != 0 {
		t.Errorf("intVal expect:0 got:%d", intVal)
	}
}
