// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package scmd

import (
	"os"
	"fmt"
	"time"
	"sync"
)


type Scmd struct {
	name string
	args []string

	muProg sync.Mutex
	prog *progress
}

func NewScmd(name string, arg ...string) *Scmd {

	return &Scmd {
		name: name,
		args: arg,
	}
}


func (m *Scmd) Start() (stdout chan []byte, stderr chan []byte, er error) {
	m.muProg.Lock()
	defer m.muProg.Unlock()

	if m.prog != nil {
		er = fmt.Errorf("cmd:%s args:%s been start pid:%d", m.name, m.args, m.prog.GetPid())
		return
	}

	m.prog, stdout, stderr, er = Newprogress(m.name, m.args...)
	return
}

// stop失败时候，返回错误原因
// 其他成功，即使已经stop，再stop也不报错误
func (m *Scmd) Stop(timeout time.Duration) error {
	m.muProg.Lock()
	defer m.muProg.Unlock()

	if m.prog == nil {
		//return fmt.Errorf("cmd:%s args:%s not start", m.name, m.args)
		return nil
	}

	if m.prog.IsStop() {
		m.prog = nil
		return nil
	}

	err := m.prog.Stop(timeout)

	if m.prog.IsStop() {
		m.prog = nil
		return nil

	} else {
		return err
	}
}

// 给进程发信号
func (m *Scmd) Signal(sig os.Signal) error {
	m.muProg.Lock()
	defer m.muProg.Unlock()

	if m.prog == nil {
		return fmt.Errorf("process can not handle")
	}

	if m.prog.IsStop() {
		return fmt.Errorf("process been stop")
	}

	return m.prog.Signal(sig)

}


func (m *Scmd) IsStop() bool {
	m.muProg.Lock()
	defer m.muProg.Unlock()

	if m.prog == nil {
		return true
	}

	return m.prog.IsStop()
}

/*
func (m *Scmd) Restart(timeout time.Duration) error {
	err := m.Stop()
	if err != nil {
		return err
	}

}
*/

// 阻塞调用， 直到程序结束才会返回
func (m *Scmd) StartWaitOutput() (stdout []byte, stderr []byte, er error) {
	stdoutChan, stderrChan, err := m.Start()
	if err != nil {
		return nil, nil, err
	}

	stdout = make([]byte, 0)
	stderr = make([]byte, 0)
	for {
		select {
		case b, ok := <-stdoutChan:
			if !ok {
				stdoutChan = nil
			} else {
				//fmt.Printf("STDOUT n:%d %s\n", len(b), b)
				stdout = append(stdout, b...)
			}

		case b, ok := <-stderrChan:
			if !ok {
				stderrChan = nil
			} else {
				//fmt.Printf("STDERR n:%d %s\n", len(b), b)
				stderr = append(stderr, b...)
			}
		}

		if stdoutChan == nil &&
			stderrChan == nil {
			break
		}
	}
	// 设置结束
	go m.Stop(time.Second)

	return stdout, stderr, nil

}




func (m *Scmd) StartTimeoutWaitOutput(timeout time.Duration) (stdout []byte, stderr []byte, er error) {

	over := make(chan bool)

	go func() {
		stdout, stderr, er = m.StartWaitOutput()
		select {
		case over <- true:
		default:
		}
	}()


	select {
	case <-over:
	case <-time.After(timeout):
		// stop 可能会失败，但是调用这个函数的
		// 这种问题应该不需要考虑
		m.Stop(time.Second)

		er =  fmt.Errorf("run timeout")
	}

	return
}
