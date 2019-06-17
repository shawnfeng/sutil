// Copyright 2014 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
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
	factory   FactoryFunc
}

func NewInstanceManager(factory FactoryFunc) *InstanceManager {
	return &InstanceManager{
		factory: factory,
	}
}

func (m *InstanceManager) Get(ctx context.Context, key string) Instancer {
	fun := "InstanceManager.get -->"

	var err error
	var in interface{}
	in, ok := m.instances.Load(key)
	if ok == false {

		slog.Infof(ctx, "%s newInstance, key: %s", fun, key)
		in, err = m.factory(ctx, key)
		if err != nil {
			slog.Errorf(ctx, "%s NewInstance err, key: %s, err: %s", fun, key, err.Error())
			return nil
		}

		in, _ = m.instances.LoadOrStore(key, in)
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
}
