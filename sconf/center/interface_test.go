package center

import (
	"context"
	"testing"
)

const (
	testService   = "authapi"
	testKey       = "key1"
	testNamespace = "db"
)

func TestGetTypedVar(t *testing.T) {
	ctx := context.Background()

	_ = Init(ctx, testService, []string{"application"})
	defer Stop(ctx)

	strVal := GetString(ctx, testKey)
	if strVal != "1" {
		t.Errorf("strVal expect:1 got:%v", strVal)
	}

	boolVal := GetBool(ctx, testKey)
	if boolVal != true {
		t.Error("boolVal expect:true got:false")
	}

	intVal := GetInt(ctx, testKey)
	if intVal != 1 {
		t.Errorf("intVal expect:1 got:%d", intVal)
	}
}

func TestGetTypedVarWithNamespace(t *testing.T) {
	ctx := context.Background()

	_ = Init(ctx, testService, []string{"application", testNamespace})
	defer Stop(ctx)

	strVal := GetStringWithNamespace(ctx, testNamespace, testKey)
	if strVal != "0" {
		t.Errorf("strVal expect:0 got:%v", strVal)
	}

	boolVal := GetBoolWithNamespace(ctx, testNamespace, testKey)
	if boolVal != false {
		t.Error("boolVal expect:false got:true")
	}

	intVal := GetIntWithNamespace(ctx, testNamespace, testKey)
	if intVal != 0 {
		t.Errorf("intVal expect:0 got:%d", intVal)
	}
}

