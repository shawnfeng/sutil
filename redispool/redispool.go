// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redispool

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"

	//	"os"
	"sync"
	"time"
	//	"reflect"

	"github.com/fzzy/radix/redis"

	"github.com/shawnfeng/sutil/slog"
)

const (
	TIMEOUT_INTV int64 = 200
	POOL_SIZE    int   = 512
)

type RedisEntry struct {
	client *redis.Client
	addr   string
	stamp  int64
}

func (self *RedisEntry) String() string {
	return fmt.Sprintf("%p@%s@%d", self.client, self.addr, self.stamp)

}

func (self *RedisEntry) Cmd(args []interface{}) *redis.Reply {
	value := args[0].(string)

	return self.client.Cmd(value, args[1:]...)

}

func (self *RedisEntry) close() {
	fun := "RedisEntry.close"
	slog.Infof("%s re:%s", fun, self)

	err := self.client.Close()
	if err != nil {
		slog.Infof("%s err re:%s err:%s", fun, self, err)
	}

}

type luaScript struct {
	sha1 string
	data []byte
}

type RedisPool struct {
	mu      sync.RWMutex
	clipool map[string]chan *RedisEntry

	muLua sync.RWMutex
	luas  map[string]*luaScript

	poolLen int
}

func (self *RedisPool) add(addr string) (*RedisEntry, error) {
	fun := "RedisPool.add"
	slog.Infof("%s addr:%s", fun, addr)

	c, err := redis.DialTimeout("tcp", addr, time.Duration(300)*time.Second)
	if err != nil {
		return nil, err
	}

	en := &RedisEntry{
		client: c,
		addr:   addr,
		stamp:  time.Now().Unix(),
	}

	return en, nil

}

func (self *RedisPool) getPool(addr string) chan *RedisEntry {
	//fun := "RedisPool.getPool -->"

	self.mu.RLock()

	tmp, ok := self.clipool[addr]
	self.mu.RUnlock()

	if !ok {
		self.mu.Lock()
		tmp, ok = self.clipool[addr]
		if !ok {
			tmp = make(chan *RedisEntry, self.poolLen)
			self.clipool[addr] = tmp
		}
		self.mu.Unlock()
	}

	return tmp
}

func (self *RedisPool) get(addr string) (*RedisEntry, error) {
	//fun := "RedisPool.get"

	po := self.getPool(addr)

	var err error
	var entry *RedisEntry
	select {
	case entry = <-po:
		//slog.Infof("%s get:%s len:%d", fun, addr, len(po))
	default:
		entry, err = self.add(addr)
	}

	return entry, err
}

func (self *RedisPool) payback(addr string, re *RedisEntry) {
	fun := "RedisPool.payback"

	po := self.getPool(addr)

	select {
	case po <- re:
		//slog.Infof("%s payback:%s len:%d", fun, addr, len(po))
	default:
		slog.Errorf("%s full not payback:%s len:%d", fun, addr, len(po))
		re.close()
	}
}

// 只对一个redis执行命令
func (self *RedisPool) CmdSingleRetry(addr string, cmd []interface{}, retrytimes int) *redis.Reply {
	fun := "RedisPool.CmdSingleRetry"
	c, err := self.get(addr)
	if err != nil {
		es := fmt.Sprintf("get conn retrytimes:%d addr:%s err:%s", retrytimes, addr, err)
		slog.Infoln(fun, es)
		return &redis.Reply{Type: redis.ErrorReply, Err: errors.New(es)}
	}

	rp := c.Cmd(cmd)
	if rp.Type == redis.ErrorReply {
		if retrytimes > 0 {
			slog.Warnf("%s redis Cmd try:%d error %s", fun, retrytimes, rp)
		} else {
			slog.Warnf("%s redis Cmd try:%d error %s", fun, retrytimes, rp)
		}
		if rp.String() == "EOF" {
			if retrytimes > 0 {
				return rp
			}
			// redis 连接timeout，重试一次
			return self.CmdSingleRetry(addr, cmd, retrytimes+1)
		}

		c.close()
	} else {
		self.payback(addr, c)
	}

	return rp

}

func (self *RedisPool) CmdSingle(addr string, cmd []interface{}) *redis.Reply {
	return self.CmdSingleRetry(addr, cmd, 0)

}

func (self *RedisPool) sha1Lua(key string) (string, error) {
	self.muLua.RLock()
	defer self.muLua.RUnlock()
	if v, ok := self.luas[key]; ok {
		return v.sha1, nil
	} else {
		return "", errors.New("lua not find")
	}

}

func (self *RedisPool) dataLua(key string) ([]byte, error) {
	self.muLua.RLock()
	defer self.muLua.RUnlock()

	if v, ok := self.luas[key]; ok {
		return v.data, nil
	} else {
		return []byte{}, errors.New("lua not find")
	}

}

func (self *RedisPool) LoadLuaFile(key, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	h := sha1.Sum(data)
	hex := fmt.Sprintf("%x", h)

	slog.Infof("RedisPool.loadLuaFile key:%s sha1:%s file:%s", key, hex, file)

	self.muLua.Lock()
	defer self.muLua.Unlock()

	self.luas[key] = &luaScript{
		sha1: hex,
		data: data,
	}

	return nil

}

// lua 脚本执行的快捷命令
func (self *RedisPool) EvalSingle(addr string, key string, cmd_args []interface{}) *redis.Reply {
	fun := "RedisPool.EvalSingle"
	sha1, err := self.sha1Lua(key)
	if err != nil {
		es := fmt.Sprintf("get lua sha1 add:%s key:%s err:%s", addr, key, err)
		return &redis.Reply{Type: redis.ErrorReply, Err: errors.New(es)}
	}

	cmd := append([]interface{}{"evalsha", sha1}, cmd_args...)
	rp := self.CmdSingle(addr, cmd)
	if rp.Type == redis.ErrorReply && rp.String() == "NOSCRIPT No matching script. Please use EVAL." {
		slog.Infoln(fun, "load lua", addr)
		cmd[0] = "eval"
		cmd[1], _ = self.dataLua(key)
		rp = self.CmdSingle(addr, cmd)
	}

	return rp
}

func (self *RedisPool) Cmd(multi_args map[string][]interface{}) map[string]*redis.Reply {
	rv := make(map[string]*redis.Reply)
	for k, v := range multi_args {
		rv[k] = self.CmdSingle(k, v)
	}

	return rv

}

func HashRedis(addrs []string, key string) string {
	h := fnv.New32a()
	h.Write([]byte(key))
	hv := h.Sum32()

	return addrs[hv%uint32(len(addrs))]

}

func NewRedisPool(poolLen int) *RedisPool {
	if poolLen <= 0 {
		poolLen = POOL_SIZE
	}
	return &RedisPool{
		clipool: make(map[string]chan *RedisEntry),
		luas:    make(map[string]*luaScript),
		poolLen: poolLen,
	}
}

//////////
//TODO
// OK 1. timeout remove
// 2. multi addr channel get
// 3. single addr multi cmd
// 4. pool conn ceil controll
