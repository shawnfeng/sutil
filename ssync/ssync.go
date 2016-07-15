// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package ssync

import (
	"sync"
	//"fmt"
)



type Mutex struct {
	once sync.Once
	mu chan bool
}


func (m *Mutex) initLock() {
	m.once.Do(func(){
		//fmt.Println("init lock")
		m.mu = make(chan bool, 1)
	})
}

func (m *Mutex) Lock() {
	m.initLock()
	m.mu <- true
}

func (m *Mutex) Unlock() {
	select {
	case <- m.mu:
	default:
		panic("ssync: unlock of unlocked mutex")
	}
}


func (m *Mutex) Trylock() bool {
	m.initLock()
	select {
	case m.mu <- true:
		return true
	default:
		return false
	}
}
