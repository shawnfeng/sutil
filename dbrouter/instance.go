// Copyright 2014 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"sync"
)

//var DefaultInstanceManager = NewInstanceManager(Factory)

type Instancer interface {
	GetType() string
	Close() error
}

type FactoryFunc func(ctx context.Context, key string) (in Instancer, err error)

type InstanceManager struct {
	instances sync.Map
	shadowInstances sync.Map
	factory   FactoryFunc
}

func NewInstanceManager(factory FactoryFunc, dbChangeChan chan dbInstanceChange) *InstanceManager {
	instanceManager := &InstanceManager{
		factory: factory,
	}

	go instanceManager.dbInsChangeHandler(context.Background(), dbChangeChan)

	return instanceManager
}

func (m *InstanceManager) Get(ctx context.Context, key string) Instancer {
	fun := "InstanceManager.get -->"

	var err error
	var in interface{}
	var instanceMap sync.Map

	// TODO 确定是否压测的标识
	group := scontext.GetGroup(ctx)
	switch group {
	case "":
		instanceMap = m.instances
	case "default":
		instanceMap = m.instances
	case "xxx":
		instanceMap = m.shadowInstances
	default:
		// TODO 这种情况容不容易出现？
		slog.Errorf(ctx, "%s invalid context group: %s", fun, group)
		return nil
	}
	/*if isTest, ok := ctx.Value("xxx").(bool); ok {
		if isTest {
			instanceMap = m.shadowInstances
		} else {
			instanceMap = m.instances
		}
	} else {
		instanceMap = m.instances
	}*/

	in, ok := instanceMap.Load(key)
	if ok == false {

		slog.Infof(ctx, "%s newInstance, key: %s", fun, key)
		in, err = m.factory(ctx, key)
		if err != nil {
			slog.Errorf(ctx, "%s NewInstance err, key: %s, err: %s", fun, key, err.Error())
			return nil
		}

		in, _ = instanceMap.LoadOrStore(key, in)
	}

	tmp, ok := in.(Instancer)
	if ok == false {
		slog.Errorf(ctx, "%s in.(Instancer) false, key: %s", fun, key)
		return nil
	}

	return tmp
}

func (m *InstanceManager) Close() {
	fun := "InstanceManager.Close -->"

	m.instances.Range(func(key, value interface{}) bool {
		slog.Infof(context.TODO(), "%s key:%v", fun, key)

		in, ok := value.(Instancer)
		if ok == false {
			slog.Errorf(context.TODO(), "%s value.(Instancer) false, key: %v", fun, key)
			return false
		}

		in.Close()
		m.instances.Delete(key)

		return true
	})

	m.shadowInstances.Range(func(key, value interface{}) bool {
		slog.Infof(context.TODO(), "%s key:%v", fun, key)

		in, ok := value.(Instancer)
		if ok == false {
			slog.Errorf(context.TODO(), "%s value.(Instancer) false, key: %v", fun, key)
			return false
		}

		in.Close()
		m.shadowInstances.Delete(key)

		return true
	})
}

func (m *InstanceManager) dbInsChangeHandler(ctx context.Context, dbChangeChan chan dbInstanceChange) {
	fun := "InstanceManager.dbInsChangeHandler -->"
	for dbInsChange := range dbChangeChan {
		slog.Infof(ctx, "%s receive db instance changes: %+v", fun, dbInsChange)
		for _, insName := range dbInsChange.dbInsChanges {
			m.closeDbInstance(ctx, insName)
		}

		for _, insName := range dbInsChange.shadowDbInsChanges {
			m.closeShadowDbInstance(ctx, insName)
		}
	}
}

func (m *InstanceManager) closeDbInstance(ctx context.Context, insName string) {
	fun := "InstanceManager.closeDbInstance -->"
	if ins, ok := m.instances.Load(insName); ok {
		m.instances.Delete(insName)
		if in, ok := ins.(Instancer); ok {
			if err := in.Close(); err == nil {
				slog.Infof(ctx, "%s succeed to close db instance %s", fun, insName)
			} else {
				slog.Warnf(ctx, "%s close db instance %s error: %s", fun, insName, err.Error())
			}
		} else {
			slog.Warnf(ctx, "%s close db instance %s error, not Instancer type", fun, insName)
		}
	}
}

func (m *InstanceManager) closeShadowDbInstance(ctx context.Context, insName string) {
	fun := "InstanceManager.closeShadowDbInstance -->"
	if ins, ok := m.shadowInstances.Load(insName); ok {
		m.shadowInstances.Delete(insName)
		if in, ok := ins.(Instancer); ok {
			if err := in.Close(); err == nil {
				slog.Infof(ctx, "%s succeed to close shadow db instance %s", fun, insName)
			} else {
				slog.Warnf(ctx, "%s close shadow db instance %s error: %s", fun, insName, err.Error())
			}
		} else {
			slog.Warnf(ctx, "%s close shadow db instance %s error, not Instancer type", fun, insName)
		}
	}
}
