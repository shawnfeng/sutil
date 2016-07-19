// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


/*
基于snowflake算法改造，慢id生成器，适合id产生不是非常快的场景
id表示可以再52bit完成，用double表示不会丢失精度，javascript等弱类型语音可以直接使用
基于毫秒时间戳，每毫秒最多产生2个id，过快会自动阻塞，直到毫秒递增
*/

package slowid


import (

	"sync"
	"time"
	"fmt"
)



const (
	WorkerIdBits = 10
	SequenceBits = 1

	MaxWorkerId  = -1 ^ (-1 << WorkerIdBits)
	MaxSequence  = -1 ^ (-1 << SequenceBits)
)

var (
	Since int64 = time.Date(2014, 11, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1000000
)

type Slowid struct {
	lastTimestamp uint64
	workerId      uint32
	sequence      uint32
	lock          sync.Mutex
}

func (sf *Slowid) uint64() uint64 {
	return (sf.lastTimestamp << (WorkerIdBits + SequenceBits)) |
		(uint64(sf.workerId) << SequenceBits) |
		(uint64(sf.sequence))
}

func (sf *Slowid) Next() (int64, error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	ts := timestamp()
	//fmt.Println(ts, sf.lastTimestamp, sf.sequence)
	if ts == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & MaxSequence
		if sf.sequence == 0 {
			ts = tilNextMillis(ts)
		}
	} else {
		sf.sequence = 0
	}

	if ts < sf.lastTimestamp {
		return 0, fmt.Errorf("Invalid timestamp: %v - precedes %v", ts, sf)
	}
	sf.lastTimestamp = ts
	return int64(sf.uint64()),  nil
}


func NewSlowid(workerId int) (*Slowid, error) {
	if workerId < 0 || workerId > MaxWorkerId {
		return nil, fmt.Errorf("Worker id %v is invalid", workerId)
	}
	return &Slowid{workerId: uint32(workerId)}, nil
}

func timestamp() uint64 {
	return uint64(time.Now().UnixNano()/1000000 - Since)
}

func tilNextMillis(ts uint64) uint64 {

	un := time.Now().UnixNano()
	dw := (un/1000000+1)*1000000-un
	<-time.After(time.Duration(dw))

	i := timestamp()
	//fmt.Println("TOO FAST GET", un, dw, i, ts)
	for i <= ts {
		//fmt.Println("get", i)
		i = timestamp()
	}
	return i
}

