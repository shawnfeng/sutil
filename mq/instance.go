// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/slog/slog"
)

var defaultInstanceManager = NewInstanceManager()

type MQRoleType int

const (
	RoleTypeReader MQRoleType = iota
	RoleTypeWriter
	RoleTypeDelayClient
)

var (
	InvalidMqRoleTypeStringErr = errors.New("invalid mq role type string")
)

func (t MQRoleType) String() string {
	switch t {
	case RoleTypeReader:
		return "reader"
	case RoleTypeWriter:
		return "writer"
	case RoleTypeDelayClient:
		return "delay"
	}
	// unreachable
	return ""
}

func MQRoleTypeFromInt(it int) (t MQRoleType, err error) {
	switch it {
	case 0:
		t = RoleTypeReader
	case 1:
		t = RoleTypeWriter
	case 2:
		t = RoleTypeDelayClient
	default:
		err = InvalidMqRoleTypeStringErr
	}
	return
}

type instanceConf struct {
	group     string
	role      MQRoleType
	topic     string
	groupId   string
	partition int
}

func (c *instanceConf) String() string {
	return fmt.Sprintf("%s-%d-%s-%s-%d",
		c.group, c.role, c.topic, c.groupId, c.partition)
}

func instanceConfFromString(s string) (conf *instanceConf, err error) {
	items := strings.Split(s, "-")
	if len(items) != 5 {
		return nil, fmt.Errorf("invalid instance conf string:%s", s)
	}

	conf = &instanceConf{
		group:   items[0],
		topic:   items[2],
		groupId: items[3],
	}

	it, err := strconv.Atoi(items[1])
	if err != nil {
		return nil, err
	}
	conf.role, err = MQRoleTypeFromInt(it)
	if err != nil {
		return nil, err
	}

	conf.partition, err = strconv.Atoi(items[4])
	if err != nil {
		return nil, err
	}
	return
}

type InstanceManager struct {
	instances sync.Map
	watchOnce sync.Once
	mutex     sync.Mutex
}

func NewInstanceManager() *InstanceManager {
	return &InstanceManager{}
}

func (m *InstanceManager) buildKey(conf *instanceConf) string {
	return conf.String()
}

func (m *InstanceManager) confFromKey(key string) (*instanceConf, error) {
	return instanceConfFromString(key)
}

func (m *InstanceManager) add(conf *instanceConf, in interface{}) {
	m.instances.Store(conf.String(), in)
}

func (m *InstanceManager) newInstance(ctx context.Context, conf *instanceConf) (interface{}, error) {

	switch conf.role {
	case RoleTypeReader:
		if len(conf.groupId) > 0 {
			return NewGroupReader(ctx, conf.topic, conf.groupId)
		} else {
			return NewPartitionReader(ctx, conf.topic, conf.partition)
		}

	case RoleTypeWriter:
		return NewWriter(ctx, conf.topic)

	case RoleTypeDelayClient:
		return NewDefaultDelayClient(ctx, conf.topic)

	default:
		return nil, fmt.Errorf("role %d error", conf.role)
	}
}

func (m *InstanceManager) get(ctx context.Context, conf *instanceConf) interface{} {
	fun := "InstanceManager.get -->"

	var err error
	var in interface{}
	key := m.buildKey(conf)
	in, ok := m.instances.Load(key)
	if ok == false {

		m.mutex.Lock()

		in, ok = m.instances.Load(key)
		if ok {
			m.mutex.Unlock()
			return in
		}

		slog.Infof(ctx, "%s newInstance, role:%v, topic: %s", fun, conf.role, conf.topic)
		in, err = m.newInstance(ctx, conf)
		if err != nil {
			slog.Errorf(ctx, "%s NewInstance err, topic: %s, err: %s", fun, conf.topic, err.Error())
			m.mutex.Unlock()
			return nil
		}

		in, _ = m.instances.LoadOrStore(key, in)

		m.mutex.Unlock()
	}
	return in
}

func (m *InstanceManager) applyChange(ctx context.Context, k string, change *center.Change) {
	fun := "InstanceManager.applyChange-->"
	slog.Infof(ctx, "%s apply change:%v to key:%s", fun, change, k)
	m.instances.Range(func(key, val interface{}) (ret bool) {
		ret = true

		sk, ok := key.(string)
		if !ok {
			slog.Errorf(ctx, "%s key:%v should be string", fun, key)
			return
		}

		conf, err := m.confFromKey(sk)
		if err != nil {
			slog.Errorf(ctx, "%s failed to convert key:%s to conf", fun, sk)
			return
		}

		keyParts, err := DefaultConfiger.ParseKey(ctx, k)
		if err != nil {
			slog.Errorf(ctx, "%s parse key:%s failed err:%v", fun, k, err)
			return
		}

		// NOTE: 只要 group 和 topic 相同，即认为相关的配置发生了变化
		//       为了逻辑简单，不论什么变化，都重新载入一次 instance, 不对不同的 ChangeType 单独处理
		if (keyParts.Group == conf.group || keyParts.Group == defaultRouteGroup) && keyParts.Topic == conf.topic {
			slog.Infof(ctx, "%s update instance:%v", fun, val)
			// NOTE: 关闭旧实例，重新载入新实例，若旧实例关闭失败打印日志
			if err = m.closeInstance(ctx, val, conf); err != nil {
				slog.Errorf(ctx, "%s close instance err:%v", fun, err)
			}

			in, err := m.newInstance(ctx, conf)
			if err != nil {
				m.instances.Delete(key)
				return
			}
			m.instances.Store(key, in)
		}
		return
	})
}

func (m *InstanceManager) applyChangeEvent(ctx context.Context, ce *center.ChangeEvent) {
	slog.Infoln(ctx, "got new change event:%v", ce)

	for key, change := range ce.Changes {
		// NOTE: 只需关心 MODIFY 与 DELETE 类型改变
		if change.ChangeType != center.MODIFY && change.ChangeType != center.DELETE {
			continue
		}

		m.applyChange(ctx, key, change)
	}
}

func (m *InstanceManager) watch(ctx context.Context) {
	fun := "InstanceManager.watch-->"
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
				slog.Infof(ctx, "%s context err:%v", fun, ctx.Err())
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

func (m *InstanceManager) getReader(ctx context.Context, conf *instanceConf) Reader {
	fun := "InstanceManager.getReader -->"

	in := m.get(ctx, conf)
	if in == nil {
		return nil
	}

	reader, ok := in.(Reader)
	if ok == false {
		slog.Errorf(ctx, "%s in.(Reader) err, topic: %s", fun, conf.topic)
		return nil
	}

	return reader
}

func (m *InstanceManager) getDelayClient(ctx context.Context, conf *instanceConf) *DelayClient {
	fun := "InstanceManager.getDelayClient"

	in := m.get(ctx, conf)
	if in == nil {
		return nil
	}

	client, ok := in.(*DelayClient)
	if ok == false {
		slog.Errorf(ctx, "%s in.(Reader) err, topic: %s", fun, conf.topic)
		return nil
	}
	return client
}

func (m *InstanceManager) getWriter(ctx context.Context, conf *instanceConf) Writer {
	fun := "InstanceManager.getReader -->"

	in := m.get(ctx, conf)
	if in == nil {
		return nil
	}

	writer, ok := in.(Writer)
	if ok == false {
		slog.Errorf(ctx, "%s in.(Writer) err, topic: %s", fun, conf.topic)
		return nil
	}

	return writer
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

		conf, err := m.confFromKey(skey)
		if err != nil {
			slog.Errorf(ctx, "%s key:%v, err:%s", fun, key, err)
			return false
		}

		err = m.closeInstance(ctx, value, conf)
		if err != nil {
			slog.Errorf(ctx, "%s close instance err:%v", fun, err)
		}

		m.instances.Delete(key)
		return true
	})
}

func (m *InstanceManager) closeInstance(ctx context.Context, instance interface{}, conf *instanceConf) error {
	fun := "InstanceManager.closeInstance-->"
	if conf.role == RoleTypeReader {
		reader, ok := instance.(Reader)
		if !ok {
			return fmt.Errorf("%s instance:%#v should be reader", fun, instance)
		}
		return reader.Close()
	}

	if conf.role == RoleTypeWriter {
		writer, ok := instance.(Writer)
		if !ok {
			return fmt.Errorf("%s instance:%#v should be writer", fun, instance)
		}
		return writer.Close()
	}

	return nil
}
