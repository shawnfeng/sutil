package redispool

import (
	"fmt"
	"errors"

//	"os"
	"time"
	"sync"
//	"reflect"

	"github.com/fzzy/radix/redis"

	"github.com/shawnfeng/sutil/slog"
)

const (
	TIMEOUT_INTV int64 = 200
)

type RedisEntry struct {
	client *redis.Client
	addr string
	stamp int64
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


type RedisPool struct {
	mu sync.Mutex
	clipool map[string][]*RedisEntry
}



func (self *RedisPool) add(addr string) (*RedisEntry, error) {
	fun := "RedisPool.add"
	slog.Infof("%s addr:%s", fun, addr)

	c, err := redis.DialTimeout("tcp", addr, time.Duration(300)*time.Second)
	if err != nil {
		return nil, err
	}

	en := &RedisEntry {
		client: c,
		addr: addr,
		stamp: time.Now().Unix(),
	}

	return en, nil


}

func (self *RedisPool) rmTimeout(rs *[]*RedisEntry) bool {
	fun := "RedisPool.rmTimeout"
	// 每次只检查一个最老的超时
	if len(*rs) > 0 && (time.Now().Unix()-(*rs)[0].stamp) > TIMEOUT_INTV {
		slog.Infof("%s rm timeout:%s", fun, (*rs)[0])
		(*rs)[0].close()
		*rs = (*rs)[1:]
		return true
	} else {
		return false
	}

}

func (self *RedisPool) getCache(addr string) *RedisEntry {
	fun := "RedisPool.getCache"
	//slog.Traceln(fun, "call", addr, self)

	self.mu.Lock()
	self.mu.Unlock()
	rs, ok := self.clipool[addr]
	if ok {
		if self.rmTimeout(&rs) {
			self.clipool[addr] = rs
		}
		if len(rs) == 0 {
			return nil
		} else {
			tmp := rs[len(rs)-1]
			self.clipool[addr] = rs[:len(rs)-1]

			nowstp := time.Now().Unix()
			if nowstp - tmp.stamp > TIMEOUT_INTV {
				// 对于超时的连接不再使用
				slog.Infof("%s rm timeout:%s", fun, tmp)
				tmp.close()
				return nil
			} else {
				// 更新使用时间戳
				tmp.stamp = nowstp
				return tmp

			}

		}

	} else {
		return nil
	}

}

func (self *RedisPool) payback(addr string, re *RedisEntry) {
	fun := "RedisPool.payback"
	//slog.Traceln(fun, "call", addr, self)

	self.mu.Lock()
	self.mu.Unlock()


	if rs, ok := self.clipool[addr]; ok {

		self.clipool[addr] = append(rs, re)

	} else {
		self.clipool[addr] = []*RedisEntry{re, }

	}

	slog.Tracef("%s addr:%s re:%s len:%d", fun, addr, re, len(self.clipool[addr]))

	//slog.Traceln(fun, "end", addr, self)


}

func (self *RedisPool) get(addr string) (*RedisEntry, error) {
	if r := self.getCache(addr); r != nil {
		return r, nil
	} else {
		return self.add(addr)
	}
}





// 只对一个redis执行命令
func (self *RedisPool) CmdSingleRetry(addr string, cmd []interface{}, retrytimes int) *redis.Reply {
	fun := "RedisPool.CmdSingleRetry"
	c, err := self.get(addr)
	if err != nil {
		es := fmt.Sprintf("get conn retrytimes:%d addr:%s err:%s", retrytimes, addr, err)
		slog.Infoln(fun, es)
		return &redis.Reply{Type: redis.ErrorReply, Err:errors.New(es)}
	}

	rp := c.Cmd(cmd)
	if rp.Type == redis.ErrorReply {
		slog.Errorf("%s redis Cmd try:%d error %s", fun, retrytimes, rp)
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

func (self *RedisPool) Cmd(multi_args map[string][]interface{}) map[string]*redis.Reply {
	rv := make(map[string]*redis.Reply)
	for k, v := range multi_args {
		rv[k] = self.CmdSingle(k, v)
	}

	return rv

}

func NewRedisPool() *RedisPool {
	return &RedisPool{
		clipool: make(map[string][]*RedisEntry),
	}
}


//////////
//TODO
// OK 1. timeout remove
// 2. multi addr channel get
// 3. single addr multi cmd
// 4. pool conn ceil controll



