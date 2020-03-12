// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/shawnfeng/sutil/cache/constants"
	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/slog/slog"
)

const (
	defaultGroup = "default"
	keySep       = "-"
)

type InstanceConf struct {
	Group     string
	Namespace string
	Wrapper   string
}

func (m *InstanceConf) String() string {
	return fmt.Sprintf("group:%s namespace:%s wrapper:%s", m.Group, m.Namespace, m.Wrapper)
}

func instanceConfFromString(s string) (conf *InstanceConf, err error) {
	items := strings.Split(s, keySep)
	if len(items) != 3 {
		return nil, fmt.Errorf("invalid instance conf string:%s", s)
	}

	conf = &InstanceConf{
		Group:     items[0],
		Namespace: items[1],
		Wrapper:   items[2],
	}
	return conf, nil
}

var DefaultInstanceManager = NewInstanceManager()

type InstanceManager struct {
	instances sync.Map
	watchOnce sync.Once
}

func NewInstanceManager() *InstanceManager {
	return &InstanceManager{}
}

func (m *InstanceManager) buildKey(conf *InstanceConf) string {
	return strings.Join([]string{
		conf.Group,
		conf.Namespace,
		conf.Wrapper,
	}, keySep)
}

func (m *InstanceManager) add(key string, client *Client) {
	m.instances.Store(key, client)
}

func (m *InstanceManager) newInstance(ctx context.Context, conf *InstanceConf) (*Client, error) {
	return NewClient(ctx, conf.Namespace, conf.Wrapper)
}

func (m *InstanceManager) GetInstance(ctx context.Context, conf *InstanceConf) (*Client, error) {
	fun := "InstanceManager.GetInstance -->"

	var err error
	var in interface{}
	key := m.buildKey(conf)
	in, ok := m.instances.Load(key)
	if ok == false {

		slog.Infof(ctx, "%s newInstance with conf:%v", fun, conf)
		in, err = m.newInstance(ctx, conf)
		if err != nil {
			slog.Errorf(ctx, "%s NewInstance err: %v", fun, err)
			return nil, err
		}

		in, _ = m.instances.LoadOrStore(key, in)
	}

	client, ok := in.(*Client)
	if ok == false {
		err := fmt.Errorf("in.(*Client), key:%v", key)
		slog.Errorf(ctx, "%s %s", fun, err.Error())
		return nil, err
	}

	return client, nil
}

func (m *InstanceManager) applyChange(ctx context.Context, key string, change *center.Change) {
	fun := "InstanceManager.applyChange-->"
	slog.Infof(ctx, "%s apply change:%v to key:%v", fun, change, key)
	m.instances.Range(func(k, v interface{}) (ret bool) {
		ret = true

		sk, ok := k.(string)
		if !ok {
			slog.Errorf(ctx, "%s key:%v should be string", fun, key)
			return
		}

		conf, err := instanceConfFromString(sk)
		if err != nil {
			slog.Errorf(ctx, "%s convert instances key:%s err:%v", fun, sk, err)
			return
		}

		keyParts, err := DefaultConfiger.ParseKey(ctx, key)
		if err != nil {
			slog.Errorf(ctx, "%s parse change key:%s err:%v", fun, key, err)
			return
		}

		// NOTE: 只要 namespace 和 group 相同，即认为相关的配置发生了变化
		//       为了逻辑简单，不论什么变化，都重新载入一次 instance，不对不同的 ChangeType 单独处理
		if (keyParts.Group == conf.Group || keyParts.Group == constants.DefaultRouteGroup) && keyParts.Namespace == conf.Namespace {
			slog.Infof(ctx, "%s update instance:%v", fun, v)
			// NOTE: 关闭旧实例，重新载入新实例，若旧实例关闭失败打印日志
			if err = m.closeInstance(ctx, v); err != nil {
				slog.Errorf(ctx, "%s close instance err:%v", fun, err)
			}

			in, err := m.newInstance(ctx, conf)
			if err != nil {
				m.instances.Delete(k)
				return
			}
			m.instances.Store(k, in)
		}

		return
	})
}

func (m *InstanceManager) applyChangeEvent(ctx context.Context, ce *center.ChangeEvent) {
	fun := "InstanceManager.applyChangeEvent-->"
	slog.Infof(ctx, "%s got new change event:%v", fun, ce)

	for key, change := range ce.Changes {
		if change.ChangeType != center.MODIFY && change.ChangeType != center.DELETE {
			continue
		}

		m.applyChange(ctx, key, change)
	}
}

func (m *InstanceManager) Watch(ctx context.Context) {
	fun := "InstanceManager.Watch-->"
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			slog.Errorf(ctx, "%s recover err: %v, stack: %s", fun, err, string(buf))
		}
	}()
	m.watchOnce.Do(func() {
		slog.Infof(ctx, "%s start watching updates", fun)
		ceChan := DefaultConfiger.Watch(ctx)
	Loop:
		for {
			select {
			case <-ctx.Done():
				slog.Infof(ctx, "%s context done err:%v", fun, ctx.Err())
				return
			case ce, ok := <-ceChan:
				if !ok {
					slog.Infof(ctx, "%s change event channel closed", fun)
					break Loop
				}
				m.applyChangeEvent(ctx, ce)
			}
		}
	})
}

func (m *InstanceManager) Close() {
	fun := "InstanceManager.Close -->"

	ctx := context.TODO()

	m.instances.Range(func(key, value interface{}) bool {
		slog.Infof(ctx, "%s key:%v", fun, key)

		err := m.closeInstance(ctx, value)
		if err != nil {
			slog.Errorf(ctx, "%s close instance err:%v", fun, err)
			return false
		}

		m.instances.Delete(key)
		return true
	})
}

func (m *InstanceManager) closeInstance(ctx context.Context, instance interface{}) error {
	fun := "InstanceManager.closeInstance-->"
	client, ok := instance.(*Client)
	if !ok {
		return fmt.Errorf("%s instance:%#v should be cache.redis.redis.Client", fun, instance)
	}
	return client.Close(ctx)
}
