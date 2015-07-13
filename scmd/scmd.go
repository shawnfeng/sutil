// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package scmd

import (
	"io"
	"os"
	"fmt"
	"time"
	"os/exec"

	"sync"
)


type Scmd struct {
	name string
	args []string

	muPid sync.Mutex
	pid  int

}

func NewScmd(name string, arg ...string) *Scmd {

	return &Scmd {
		name: name,
		args: arg,
	}
}

func (m *Scmd) opPid(opfun func()) {
	m.muPid.Lock()
	defer m.muPid.Unlock()

	opfun()

}

func (m *Scmd) makeReaderChan(r io.Reader) (chan []byte, chan bool) {
    read := make(chan []byte)
	over := make(chan bool)

	go func() {
		for {
			// buffer设置必须放到这，不要放到for外面，否则会造成，后面的read覆盖前面的read
			b := make([]byte, 1024)
			n, err := r.Read(b)
			fmt.Println("Debug", n, err)

			// https://golang.org/pkg/io/#Reader
			// Callers should always process the n > 0 bytes returned before considering the error err. Doing so correctly handles I/O errors that happen after reading some bytes and also both of the allowed EOF behaviors.
			if n > 0 {
				//fmt.Printf("READ: n:%d s:%s\n", n, b[0:n])
				read <- b[0:n]

			}

			if err != nil {
				fmt.Println("READER err:%s", err)
				close(read)
				over <- true
				return
			}
		}
	}()

	return read, over

}


func (m *Scmd) waitExit(cmd *exec.Cmd, stdoutOver, stderrOver chan bool) {
	// StdoutPipe returns a pipe that will be connected to the command's standard output when the command starts.

	// Wait will close the pipe after seeing the command exit, so most callers need not close the pipe themselves; however,
	// an implication is that it is incorrect to call Wait before all reads from the pipe have completed. For the same reason,
	// it is incorrect to call Run when using StdoutPipe. See the example for idiomatic usage.

	// 按照同步逻辑执行，wait先于pipe读取调用，wait就会阻塞在那里了
	// 异步测试没发现什么问题，这里实现，没有严格按照read完，才wait，因为是异步逻辑，可能会有问题么？
	// 最后还是这么做了，想把整个pid操作安全化完全，但是这里要wait保证调用时候完全不阻塞完成，防止对mutex占用时间过长
	fmt.Println("WAIT: call")

	<-stdoutOver
	fmt.Println("WAIT: stdout over")

	<-stderrOver
	fmt.Println("WAIT: stderr over")

	m.opPid(func() {
		if err := cmd.Wait(); err != nil {
			fmt.Printf("WAIT: err:%s\n", err)
		}

		m.pid = 0
	})
}

func (m *Scmd) Start() (stdout chan []byte, stderr chan []byte, er error) {

	//fmt.Println("START:", m.name, m.args)

	m.opPid(func() {
		if m.pid != 0 {
			er = fmt.Errorf("cmd been start pid:%d", m.pid)
			return
		}

		cmd := exec.Command(m.name, m.args...)

		stdoutRc, err := cmd.StdoutPipe()
		if err != nil {
			er = fmt.Errorf("init stdout pipe err:%s", err)
			return
		}

		stderrRc, err := cmd.StderrPipe()
		if err != nil {
			er = fmt.Errorf("init stderr pipe err:%s", err)
			return
		}

		err = cmd.Start()
		if err != nil {
			er = fmt.Errorf("start process err:%s", err)
			return
		}

		m.pid = cmd.Process.Pid

		// test call here, will deadlock
		// 并且这时候去查这个pid找不到，why？？
		fmt.Printf("START pid:%d\n", m.pid)
		//m.waitExit(cmd)

		var stdoutOver, stderrOver chan bool
		stdout, stdoutOver = m.makeReaderChan(stdoutRc)
		stderr, stderrOver = m.makeReaderChan(stderrRc)

		go m.waitExit(cmd, stdoutOver, stderrOver)

	})

	return

}


func (m *Scmd) Stop() (er error) {
    m.opPid(func() {
		pid := m.pid
		p, err := os.FindProcess(pid)
		if err != nil {
			er = fmt.Errorf("find pid:%d err:%v", pid, err)
		} else {
			err = p.Kill()
			if err != nil {
				er = fmt.Errorf("kill pid:%d err:%v", pid, err)
			} else {
				m.pid = 0
			}
		}
	})

	return
}

func (m *Scmd) GetPid() (pid int) {
	m.opPid(func() {
		pid = m.pid
	})
	return
}

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
				fmt.Printf("STDOUT n:%d %s\n", len(b), b)
				stdout = append(stdout, b...)
			}

		case b, ok := <-stderrChan:
			if !ok {
				stderrChan = nil
			} else {
				fmt.Printf("STDERR n:%d %s\n", len(b), b)
				stderr = append(stderr, b...)
			}
		}

		if stdoutChan == nil &&
			stderrChan == nil {
			break
		}
	}

	return stdout, stderr, nil

}


func (m *Scmd) StartTimeoutWaitOutput(timeout time.Duration) (stdout []byte, stderr []byte, er error) {

	return
}


// 获取procid，超时停止，标准错误、输出chan返回启动、timeout启动、阻塞启动结束后才返回
// 启动停止，都有error返回，需要判断，stop时候，pid为0时候停止失败


// 被kill后，wait输出是: 
// -9 signal: killed
// 不带 signal: terminated
