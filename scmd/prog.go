// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package scmd

import (
	"io"
	"os"
	"time"
	"fmt"
	"os/exec"

	"sync"
	"sync/atomic"
)


type progress struct {
	muCmd sync.Mutex
	cmd *exec.Cmd

	waitNotify chan bool

	isStop   int32
}


func (m *progress) opCmd(opfun func()) {
	m.muCmd.Lock()
	defer m.muCmd.Unlock()

	opfun()

}

func makeReaderChan(r io.Reader) (chan []byte, chan bool) {
    read := make(chan []byte)
	over := make(chan bool)

	go func() {
		for {
			// buffer设置必须放到这，不要放到for外面，否则会造成，后面的read覆盖前面的read，
			// chan 传递的[]byte是引用关系
			b := make([]byte, 1024)
			n, err := r.Read(b)
			//fmt.Println("Debug", n, err)

			// https://golang.org/pkg/io/#Reader
			// Callers should always process the n > 0 bytes returned before considering the error err. Doing so correctly handles I/O errors that happen after reading some bytes and also both of the allowed EOF behaviors.
			if n > 0 {
				//fmt.Printf("READ: n:%d s:%s\n", n, b[0:n])
				read <- b[0:n]

			}

			if err != nil {
				//fmt.Println("READER err:%s", err)
				close(read)
				over <- true
				return
			}
		}
	}()

	return read, over

}


func (m *progress) waitExit(stdoutOver, stderrOver chan bool) {
	// StdoutPipe returns a pipe that will be connected to the command's standard output when the command starts.

	// Wait will close the pipe after seeing the command exit, so most callers need not close the pipe themselves; however,
	// an implication is that it is incorrect to call Wait before all reads from the pipe have completed. For the same reason,
	// it is incorrect to call Run when using StdoutPipe. See the example for idiomatic usage.

	// 按照同步逻辑执行，wait先于pipe读取调用，wait就会阻塞在那里了
	// 异步测试没发现什么问题，这里实现，没有严格按照read完，才wait，因为是异步逻辑，可能会有问题么？
	// 最后还是这么做了，想把整个pid操作安全化完全，但是这里要wait保证调用时候完全不阻塞完成，防止对mutex占用时间过长
	//fmt.Println("WAIT: call")

	<-stdoutOver
	//fmt.Println("WAIT: stdout over")

	<-stderrOver
	//fmt.Println("WAIT: stderr over")

	m.opCmd(func() {
		if err := m.cmd.Wait(); err != nil {
			//fmt.Printf("WAIT: err:%s\n", err)
		}
		select {
		case m.waitNotify <- true:
		default:
		}

		// 同时可能多个同时stop过程,一个成功后，其他的就不要在，timeout了，也算成功
		// 让所有在read的都返回
		close(m.waitNotify)
		atomic.StoreInt32(&m.isStop, 1)
	})
}

func Newprogress(name string, arg ...string) (prog *progress, stdout chan []byte, stderr chan []byte, er error) {

	//fmt.Println("START:", m.name, m.args)

	cmd := exec.Command(name, arg...)

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

	// test call here, will deadlock
	// 并且这时候去查这个pid找不到，why？？
	//fmt.Printf("START pid:%d\n", cmd.Process.Pid)
	//m.waitExit(cmd)

	var stdoutOver, stderrOver chan bool
	stdout, stdoutOver = makeReaderChan(stdoutRc)
	stderr, stderrOver = makeReaderChan(stderrRc)


	prog = &progress {
		cmd: cmd,
		waitNotify: make(chan bool),
	}


	go prog.waitExit(stdoutOver, stderrOver)


	return

}


// 进程没有挂断，重复kill没有问题
// stop完成后保证进程已经挂掉
// 这个接口是goroutine安全的
func (m *progress) Stop(timeout time.Duration) (er error) {
    m.opCmd(func() {
		err := m.cmd.Process.Kill()
		if err != nil {
			er = fmt.Errorf("kill pid:%d err:%v", m.cmd.Process.Pid, err)
		}
	})

	if er == nil {
		// 等待进程wait返回
		select {
		case <-m.waitNotify:
		case <-time.After(timeout):
			er = fmt.Errorf("kill pid:%d stop timeout", m.cmd.Process.Pid)
		}

	}

	return
}

// 发送信号
func (m *progress) Signal(sig os.Signal) (er error) {
    m.opCmd(func() {
		err := m.cmd.Process.Signal(sig)
		if err != nil {
			er = fmt.Errorf("send signal pid:%d err:%v", m.cmd.Process.Pid, err)
		}
	})

	return er
}


// 这个接口是goroutine安全的
func (m *progress) IsStop() bool {
	return atomic.LoadInt32(&m.isStop) == 1
}

func (m *progress) GetPid() (pid int) {
	m.opCmd(func() {
		pid = m.cmd.Process.Pid
	})
	return
}


// 获取procid，超时停止，标准错误、输出chan返回启动、timeout启动、阻塞启动结束后才返回
// 启动停止，都有error返回，需要判断，stop时候，pid为0时候停止失败


// 被kill后，wait输出是: 
// -9 signal: killed
// 不带 signal: terminated


// wait sh 脚本，
// sh 被干掉后，wait可以获取都立即的返回，但是pipe确没有马上EOF
// 上面是通过findprogress测试的

// 如果直接用cmd 中start获取到process kill，测试没有发现这个问题(发现不是，是启动立即stop是这个效果)
// 原因是应该因为，还没有产生输出时候，就EOF了，所以就终止了


// StdoutPipe returns a pipe that will be connected to the command's standard output when the command starts.

// Wait will close the pipe after seeing the command exit, so most callers need not close the pipe themselves; however,
// an implication is that it is incorrect to call Wait before all reads from the pipe have completed. For the same reason,
// it is incorrect to call Run when using StdoutPipe. See the example for idiomatic usage.

// 按照同步逻辑执行，wait先于pipe读取调用，wait就会阻塞在那里了
// 异步测试没发现什么问题，这里实现，没有严格按照read完，才wait，因为是异步逻辑，可能会有问题么？

// 让wait就单纯的去wait，不要调整pid的逻辑了


// os
// func (p *Process) Release() error
// Release releases any resources associated with the Process p, rendering it unusable in the future. Release only needs to be called if Wait is not.
