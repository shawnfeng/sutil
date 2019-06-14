// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/slog"
	"sync"
)

var DefaultInstanceManager = NewInstanceManager()

type InstanceManager struct {
	instances sync.Map
}

func NewInstanceManager() *InstanceManager {
	return &InstanceManager{}
}

func (m *InstanceManager) buildKey(flag, namespace string) string {
	return fmt.Sprintf("%s-%s", flag, namespace)
}

func (m *InstanceManager) add(flag, namespace string, client *Client) {
	m.instances.Store(m.buildKey(flag, namespace), client)
}

func (m *InstanceManager) newInstance(ctx context.Context, namespace string) (*Client, error) {
	return NewClient(ctx, namespace)
}

func (m *InstanceManager) GetInstance(ctx context.Context, namespace string) *Client {
	fun := "InstanceManager.GetInstance -->"

	var flag string
	var err error
	var in interface{}
	key := m.buildKey(flag, namespace)
	in, ok := m.instances.Load(key)
	if ok == false {

		slog.Infof(ctx, "%s newInstance, flag: %s, namespace: %s", fun, flag, namespace)
		in, err = m.newInstance(ctx, namespace)
		if err != nil {
			slog.Errorf(ctx, "%s NewInstance err, namespace: %s, err: %s", fun, namespace, err.Error())
			return nil
		}

		in, _ = m.instances.LoadOrStore(key, in)
	}

	client, ok := in.(*Client)
	if ok == false {
		slog.Errorf(ctx, "%s in.(*Client), key:%v", fun, key)
		return nil
	}

	return client
}

func (m *InstanceManager) Close() {
	fun := "InstanceManager.Close -->"

	ctx := context.TODO()

	m.instances.Range(func(key, value interface{}) bool {
		slog.Infof(ctx, "%s key:%v", fun, key)

		skey, ok := key.(string)
		if ok == false {
			slog.Errorf(ctx, "%s key:%v", fun, key)
			return false
		}

		client, ok := value.(*Client)
		if ok == false {
			slog.Errorf(ctx, "%s value.(*Client), key:%v", fun, skey)
			return false
		}

		err := client.Close()
		if err != nil {
			slog.Errorf(ctx, "%s client.Close, key:%v", fun, skey)
			return false
		}

		m.instances.Delete(key)
		return true
	})
}
