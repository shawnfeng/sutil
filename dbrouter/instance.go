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

type FactoryFunc func(ctx context.Context, key, group string) (in Instancer, err error)

type InstanceManager struct {
	instances map[string]*sync.Map
	factory   FactoryFunc
	groups    []string
}

func NewInstanceManager(factory FactoryFunc, dbChangeChan chan dbConfigChange, groups []string) *InstanceManager {
	instanceManager := &InstanceManager{
		factory: factory,
		instances: make(map[string]*sync.Map),
		groups: groups,
	}

	go instanceManager.dbInsChangeHandler(context.Background(), dbChangeChan)

	return instanceManager
}

func (m *InstanceManager) Get(ctx context.Context, key string) Instancer {
	fun := "InstanceManager.get -->"

	var err error
	var in interface{}
	var instanceMap *sync.Map
	group := scontext.GetGroup(ctx)

	isConfig := false
	for _, configGroup := range m.groups {
		if group == configGroup {
			isConfig = true
		}
	}

	if !isConfig {
		group = DefaultGroup
	}

	if ins, ok := m.instances[group]; ok {
		instanceMap = ins
	} else {
		m.instances[group] = new(sync.Map)
		instanceMap = m.instances[group]
	}

	in, ok := instanceMap.Load(key)
	if ok == false {

		slog.Infof(ctx, "%s newInstance, key: %s", fun, key)
		in, err = m.factory(ctx, key, group)
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

	for group, ins := range m.instances {
		ins.Range(func(key, value interface{}) bool {
			slog.Infof(context.TODO(), "%s key:%v", fun, key)

			in, ok := value.(Instancer)
			if ok == false {
				slog.Errorf(context.TODO(), "%s value.(Instancer) false, key: %v", fun, key)
				return false
			}

			in.Close()
			ins.Delete(key)

			return true
		})

		delete(m.instances, group)
	}
}

func (m *InstanceManager) dbInsChangeHandler(ctx context.Context, dbChangeChan chan dbConfigChange) {
	fun := "InstanceManager.dbInsChangeHandler -->"
	var originGroups []string
	for dbInsChange := range dbChangeChan {
		slog.Infof(ctx, "%s receive db instance changes: %+v", fun, dbInsChange)
		originGroups = m.groups
		m.groups = dbInsChange.dbGroups
		for group, insNames := range dbInsChange.dbInstanceChange {
			for _, insName := range insNames {
				m.closeDbInstance(ctx, insName, group)
			}
		}

		for _, group := range originGroups {
			if !isInStringList(group, m.groups) {
				delete(m.instances, group)
				slog.Infof(ctx, "%s succeed delete group %s", fun, group)
			}
		}
	}
}

func isInStringList(item string, strList []string) bool {
	for _, str := range strList {
		if str == item {
			return true
		}
	}

	return false
}

func (m *InstanceManager) closeDbInstance(ctx context.Context, insName, group string) {
	fun := "InstanceManager.closeDbInstance -->"
	if _, ok := m.instances[group]; !ok {
		return
	}
	if ins, ok := m.instances[group].Load(insName); ok {
		m.instances[group].Delete(insName)
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
