// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"errors"
	"fmt"
	"github.com/shawnfeng/sutil/slog/slog"
	"strconv"
	"strings"
	"sync"
)

var DefaultInstanceManager = NewInstanceManager()

type MQRoleType int

const (
	RoleTypeReader MQRoleType = iota
	RoleTypeWriter
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
	default:
		err = InvalidMqRoleTypeStringErr
	}
	return
}

type instanceConf struct {
	group string
	role MQRoleType
	topic string
	groupId string
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
		group: items[0],
		topic: items[2],
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

		slog.Infof(ctx, "%s newInstance, role:%v, topic: %s", fun, conf.role, conf.topic)
		in, err = m.newInstance(ctx, conf)
		if err != nil {
			slog.Errorf(ctx, "%s NewInstance err, topic: %s, err: %s", fun, conf.topic, err.Error())
			return nil
		}

		in, _ = m.instances.LoadOrStore(key, in)
	}
	return in
}

func (m *InstanceManager) watchInstance(ctx context.Context) {

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

		if conf.role == RoleTypeReader {
			reader, ok := value.(Reader)
			if ok == false {
				slog.Errorf(ctx, "%s value.(Reader), key:%v", fun, key)
				return false
			}

			reader.Close()
		}

		if conf.role == RoleTypeWriter {
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
