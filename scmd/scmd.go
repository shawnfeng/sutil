// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package scmd

import (
	"io"
	"os"
	"fmt"
	//"time"
	"os/exec"

	"sync/atomic"
)


type Scmd struct {
	name string
	args []string

	pid  int32

}

func NewScmd(name string, arg ...string) *Scmd {

	return &Scmd {
		name: name,
		args: arg,
	}
}

func (m *Scmd) makeReaderChan(r io.Reader) (chan []byte) {
    read := make(chan []byte)

	go func() {
		for {
			// buffer设置必须放到这，不要放到for外面，否则会造成，后面的read覆盖前面的read
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
				close(read)
				return
			}
		}
	}()

	return read

}


func (m *Scmd) Start() (stdout chan []byte, stderr chan []byte, er error) {

	fmt.Println("START:", m.name, m.args)

	cmd := exec.Command(m.name, m.args...)

	stdoutRc, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderrRc, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	atomic.StoreInt32(&m.pid, int32(cmd.Process.Pid))

	return m.makeReaderChan(stdoutRc), m.makeReaderChan(stderrRc), nil

	// 不需要wait的过程，chan 使用者，可以读取chan
	// 就可以正常的判断出程序的结束
}


func (m *Scmd) Stop() error {
	pid := m.GetPid()
	p, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("find pid:%d err:%v", pid, err)
	} else {
		err = p.Kill()
		if err != nil {
			return fmt.Errorf("kill pid:%d err:%v", pid, err)
		} else {
			return nil
		}
	}
}

func (m *Scmd) GetPid() int {
	return int(atomic.LoadInt32(&m.pid))
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

	return stdout, stderr, nil

}

/*
func (m *Scmd) StartTimeoutWaitOutput(timeout time.Duration) (stdout []byte, stderr []byte, er error) {

}
*/

// 获取procid，超时停止，标准错误、输出chan返回启动、timeout启动、阻塞启动结束后才返回
// 启动停止，都有error返回，需要判断，stop时候，pid为0时候停止失败


// 被kill后，wait输出是: signal: killed
