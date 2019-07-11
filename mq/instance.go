// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"fmt"
	"github.com/shawnfeng/sutil/scontext"
	"github.com/shawnfeng/sutil/slog/slog"
	"strings"
	"sync"
)

var DefaultInstanceManager = NewInstanceManager()

const (
	RoleTypeReader = iota
	RoleTypeWriter
)

type InstanceManager struct {
	instances sync.Map
}

func NewInstanceManager() *InstanceManager {
	return &InstanceManager{}
}

func (m *InstanceManager) buildKey(flag string, role int, topic, groupId string, partition int) string {
	return fmt.Sprintf("%s-%d-%s-%s-%d", flag, role, topic, groupId, partition)
}

func (m *InstanceManager) add(flag string, role int, topic, groupId string, partition int, in interface{}) {
	m.instances.Store(m.buildKey(flag, role, topic, groupId, partition), in)
}

func (m *InstanceManager) getRole(key string) (int, error) {
	items := strings.Split(key, "-")
	if len(items) != 5 {
		return 0, fmt.Errorf("key error, key:%s", key)
	}

	if items[1] == "0" {
		return RoleTypeReader, nil
	}

	if items[1] == "1" {
		return RoleTypeWriter, nil
	}

	return 0, fmt.Errorf("key error, key:%s", key)
}

func (m *InstanceManager) newInstance(ctx context.Context, role int, topic, groupId string, partition int) (interface{}, error) {

	switch role {
	case RoleTypeReader:
		if len(groupId) > 0 {
			return NewGroupReader(ctx, topic, groupId)
		} else {
			return NewPartitionReader(ctx, topic, partition)
		}

	case RoleTypeWriter:
		return NewWriter(ctx, topic)

	default:
		return nil, fmt.Errorf("role %d error", role)
	}
}

func (m *InstanceManager) get(ctx context.Context, role int, topic, groupId string, partition int) interface{} {
	fun := "InstanceManager.get -->"

	group := scontext.GetGroup(ctx)

	var err error
	var in interface{}
	key := m.buildKey(group, role, topic, groupId, partition)
	in, ok := m.instances.Load(key)
	if ok == false {

		slog.Infof(ctx, "%s newInstance, role:%d, topic: %s", fun, role, topic)
		in, err = m.newInstance(ctx, role, topic, groupId, partition)
		if err != nil {
			slog.Errorf(ctx, "%s NewInstance err, topic: %s, err: %s", fun, topic, err.Error())
			return nil
		}

		in, _ = m.instances.LoadOrStore(key, in)
	}
	return in
}

func (m *InstanceManager) getReader(ctx context.Context, role int, topic, groupId string, partition int) Reader {
	fun := "InstanceManager.getReader -->"

	in := m.get(ctx, role, topic, groupId, partition)
	if in == nil {
		return nil
	}

	reader, ok := in.(Reader)
	if ok == false {
		slog.Errorf(ctx, "%s in.(Reader) err, topic: %s", fun, topic)
		return nil
	}

	return reader
}

func (m *InstanceManager) getWriter(ctx context.Context, role int, topic, groupId string, partition int) Writer {
	fun := "InstanceManager.getReader -->"

	in := m.get(ctx, role, topic, groupId, partition)
	if in == nil {
		return nil
	}

	writer, ok := in.(Writer)
	if ok == false {
		slog.Errorf(ctx, "%s in.(Writer) err, topic: %s", fun, topic)
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

		role, err := m.getRole(skey)
		if err != nil {
			slog.Errorf(ctx, "%s key:%v, err:%s", fun, key, err)
			return false
		}

		if role == RoleTypeReader {
			reader, ok := value.(Reader)
			if ok == false {
				slog.Errorf(ctx, "%s value.(Reader), key:%v", fun, key)
				return false
			}

			reader.Close()
		}

		if role == RoleTypeWriter {
			writer, ok := value.(Writer)
			if ok == false {
				slog.Errorf(ctx, "%s value.(Writer), key:%v", fun, key)
				return false
			}

			writer.Close()
		}

		m.instances.Delete(key)
		return true
	})
}
