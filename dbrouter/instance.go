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
	instanceMu sync.RWMutex
	groupMu    sync.RWMutex
	instances  map[string]Instancer
	factory    FactoryFunc
	groups     []string
}

func NewInstanceManager(factory FactoryFunc, dbChangeChan chan dbConfigChange, groups []string) *InstanceManager {
	instanceManager := &InstanceManager{
		instances: make(map[string]Instancer),
		factory:   factory,
		groups:    groups,
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
	var in Instancer
	group := scontext.GetGroup(ctx)

	if group != DefaultGroup {
		if !m.isInGroup(group) {
			group = DefaultGroup
		}
	}

	key := m.buildKey(instance, group)
	in, ok := m.getInstance(ctx, key)
	if ok == false {
		slog.Infof(ctx, "%s newInstance, instance: %s", fun, instance)
		in, err = m.buildInstance(ctx, instance, group)
		if err != nil {
			slog.Errorf(ctx, "%s NewInstance err, instance: %s, err: %s", fun, instance, err.Error())
			return nil
		}

	}

	return in
}

func (m *InstanceManager) getInstance(ctx context.Context, key string) (Instancer, bool) {
	m.instanceMu.RLock()
	defer m.instanceMu.RUnlock()

	in, ok := m.instances[key]
	return in, ok
}

func (m *InstanceManager) isInGroup(group string) bool {
	m.groupMu.RLock()
	defer m.groupMu.RUnlock()

	for _, configGroup := range m.groups {
		if group == configGroup {
			return true
		}
	}
	return false
}

func (m *InstanceManager) buildInstance(ctx context.Context, instance, group string) (Instancer, error) {
	if group != DefaultGroup {
		if !m.isInGroup(group) {
			group = DefaultGroup
		}
	}
	key := m.buildKey(instance, group)

	m.instanceMu.Lock()
	defer m.instanceMu.Unlock()

	if in, ok := m.instances[key]; ok {
		return in, nil
	}

	in, err := m.factory(ctx, instance, group)
	if err != nil {
		return nil, err
	}

	m.instances[key] = in
	return in, nil
}

func (m *InstanceManager) Close() {
	fun := "InstanceManager.Close -->"
	m.instanceMu.Lock()
	defer m.instanceMu.Unlock()

	for key, in := range m.instances {
		slog.Infof(context.TODO(), "%s key:%v", fun, key)
		go in.Close()
		delete(m.instances, key)
	}
}

func (m *InstanceManager) dbInsChangeHandler(ctx context.Context, dbChangeChan chan dbConfigChange) {
	fun := "InstanceManager.dbInsChangeHandler -->"
	for dbInsChange := range dbChangeChan {
		slog.Infof(ctx, "%s receive db instance changes: %+v", fun, dbInsChange)
		m.handleGroupChange(dbInsChange.dbGroups)
		m.handleDbInsChange(ctx, dbInsChange.dbInstanceChange)
	}
}

func (m *InstanceManager) handleGroupChange(groups []string) {
	m.groupMu.Lock()
	defer m.groupMu.Unlock()

	m.groups = groups
}

func (m *InstanceManager) handleDbInsChange(ctx context.Context, dbInstanceChange map[string][]string) {
	for group, insNames := range dbInstanceChange {
		for _, insName := range insNames {
			m.closeDbInstance(ctx, insName, group)
		}
	}
}

func (m *InstanceManager) closeDbInstance(ctx context.Context, insName, group string) {
	fun := "InstanceManager.closeDbInstance -->"
	key := m.buildKey(insName, group)

	m.instanceMu.Lock()
	defer m.instanceMu.Unlock()

	if in, ok := m.instances[key]; ok {
		delete(m.instances, key)
		go func() {
			if err := in.Close(); err == nil {
				slog.Infof(ctx, "%s succeed to close db instance: %s group: %s", fun, insName, group)
			} else {
				slog.Warnf(ctx, "%s close db instance: %s group: %s error: %s", fun, insName, group, err.Error())
			}
		}()
	}
}
