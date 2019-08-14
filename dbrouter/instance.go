// Copyright 2014 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"fmt"
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
	instanceMu sync.Mutex
	instances  sync.Map
	factory    FactoryFunc
	groups     []string
}

func NewInstanceManager(factory FactoryFunc, dbChangeChan chan dbConfigChange, groups []string) *InstanceManager {
	instanceManager := &InstanceManager{
		factory: factory,
		groups: groups,
	}

	go instanceManager.dbInsChangeHandler(context.Background(), dbChangeChan)

	return instanceManager
}

func (m *InstanceManager) buildKey(instance, group string) string {
	return fmt.Sprintf("%s-%s", group, instance)
}

func (m *InstanceManager) Get(ctx context.Context, instance string) Instancer {
	fun := "InstanceManager.Get -->"

	var err error
	var in interface{}
	group := scontext.GetGroup(ctx)

	if !m.isInGroup(group) {
		group = DefaultGroup
	}

	key := m.buildKey(instance, group)
	in, ok := m.instances.Load(key)
	if ok == false {
		slog.Infof(ctx, "%s newInstance, instance: %s", fun, instance)
		in, err = m.buildInstance(ctx, instance, group)
		if err != nil {
			slog.Errorf(ctx, "%s NewInstance err, instance: %s, err: %s", fun, instance, err.Error())
			return nil
		}

	}

	tmp, ok := in.(Instancer)
	if ok == false {
		slog.Errorf(ctx, "%s in.(Instancer) false, key: %s", fun, key)
		return nil
	}

	return tmp
}

func (m *InstanceManager) isInGroup(group string) bool {
	for _, configGroup := range m.groups {
		if group == configGroup {
			return true
		}
	}
	return false
}

func (m *InstanceManager) buildInstance(ctx context.Context, instance, group string) (interface{}, error) {
	m.instanceMu.Lock()
	defer m.instanceMu.Unlock()
	if group != DefaultGroup {
		if !m.isInGroup(group) {
			group = DefaultGroup
		}
	}
	tmp, err := m.factory(ctx, instance, group)
	if err != nil {
		return nil, err
	}

	key := m.buildKey(instance, group)
	in, _ := m.instances.LoadOrStore(key, tmp)

	return in, nil
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
}

func (m *InstanceManager) dbInsChangeHandler(ctx context.Context, dbChangeChan chan dbConfigChange) {
	fun := "InstanceManager.dbInsChangeHandler -->"
	for dbInsChange := range dbChangeChan {
		slog.Infof(ctx, "%s receive db instance changes: %+v", fun, dbInsChange)
		m.groups = dbInsChange.dbGroups
		m.handleDbInsChange(ctx, dbInsChange.dbInstanceChange)
	}
}

func (m *InstanceManager) handleDbInsChange(ctx context.Context, dbInstanceChange map[string][]string) {
	m.instanceMu.Lock()
	defer m.instanceMu.Unlock()
	for group, insNames := range dbInstanceChange {
		for _, insName := range insNames {
			m.closeDbInstance(ctx, insName, group)
		}
	}
}

func (m *InstanceManager) closeDbInstance(ctx context.Context, insName, group string) {
	fun := "InstanceManager.closeDbInstance -->"
	key := m.buildKey(insName, group)
	if ins, ok := m.instances.Load(key); ok {
		m.instances.Delete(key)
		if in, ok := ins.(Instancer); ok {
			if err := in.Close(); err == nil {
				slog.Infof(ctx, "%s succeed to close db instance: %s group: %s", fun, insName, group)
			} else {
				slog.Warnf(ctx, "%s close db instance: %s group: %s error: %s", fun, insName, group, err.Error())
			}
		} else {
			slog.Warnf(ctx, "%s close db instance: %s group: %s error, not Instancer type", fun, insName, group)
		}
	}
}
