// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package scmd

import (
	"time"
	"testing"

	"log"
)

func TestIt(t *testing.T) {
	//testProg0(t)
	//testProg1(t)
	//testProg2(t)

	testStart0(t)
	//time.Sleep(20 * time.Second)
	//testStart1(t)
	//testStartStop(t)
}

func testProg1(t *testing.T) {
	//p, _, _, err := Newprogress("grep", "A", "/home/fenggx/shawn_go/src/github.com/shawnfeng/sutil/scmd/*")
	p, _, _, err := Newprogress("sh", "t.sh")

	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(time.Millisecond)

	log.Println("call stop")
	err = p.Stop(time.Second*5)
	log.Println("stop timeout check", err)

}

func testProg0(t *testing.T) {
	
	//p, _, _, err := Newprogress("grep", "A", "/home/fenggx/shawn_go/src/github.com/shawnfeng/sutil/scmd/*")
	p, _, _, err := Newprogress("sh", "t.sh")

	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(time.Second*10)

	err = p.Stop(1)
	log.Println("stop 0", err)
	time.Sleep(time.Second)
	err = p.Stop(1)
	log.Println("stop 1", err)

	err = p.Stop(1)
	time.Sleep(time.Second)
	log.Println("stop 2", err)


	time.Sleep(time.Second*10)

}


func testProg2(t *testing.T) {
	
	//p, _, _, err := Newprogress("grep", "A", "/home/fenggx/shawn_go/src/github.com/shawnfeng/sutil/scmd/*")
	p, stdoutChan, stderrChan, err := Newprogress("sh", "t.sh")

	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(time.Second*1)

	go func() {
		err := p.Stop(time.Second * 15)
		log.Println("stop 0", err)
	}()

	//time.Sleep(time.Second)

	go func() {
		err := p.Stop(time.Second * 15)
		log.Println("stop 1", err)
	}()

	//time.Sleep(time.Second)


	go func() {
		err := p.Stop(time.Second * 15)
		log.Println("stop 2", err)
	}()


	for {
		select {
		case b, ok := <-stdoutChan:
			if !ok {
				stdoutChan = nil
			} else {
				log.Printf("STDOUT n:%d %s\n", len(b), b)
			}

		case b, ok := <-stderrChan:
			if !ok {
				stderrChan = nil
			} else {
				log.Printf("STDERR n:%d %s\n", len(b), b)
			}
		}

		if stdoutChan == nil &&
			stderrChan == nil {
			break
		}
	}

	log.Println("WAIT OVER")

	time.Sleep(time.Second*10)

}




func testStart0(t *testing.T) {
	
	c := NewScmd("grep", "A", "/home/fenggx/shawn_go/src/github.com/shawnfeng/sutil/scmd/*")  // * 执行是有错误的，因为  *  通配符是shell支持的
	//c := NewScmd("grep", "A", "/home/fenggx/shawn_go/src/github.com/shawnfeng/sutil/scmd/scmd.go")
	//c := NewScmd("ls", "/tmp")

	//c := NewScmd("echo", "AAAA", "BBBB")

	stdout, stderr, err := c.StartWaitOutput()
	if err != nil {
		t.Errorf("%s", err)
	}


	log.Printf("o:%s e:%s c:%s", stdout, stderr, c.prog)

	time.Sleep(time.Second)
	log.Printf("o:%s e:%s c:%s", stdout, stderr, c.prog)
}



func testStart1(t *testing.T) {

	c := NewScmd("sleep", "20")

	stdout, stderr, err := c.StartWaitOutput()
	if err != nil {
		t.Errorf("%s", err)
	}


	log.Printf("o:%s e:%s", stdout, stderr)
	time.Sleep(time.Second)

}



func testStartStop(t *testing.T) {

	//c := NewScmd("sh", "t.sh")
	c := NewScmd("sleep", "10")

	go func() {
		log.Printf("go do start")
		stdout, stderr, err := c.StartWaitOutput()
		if err != nil {
			t.Errorf("%s", err)
		}
		log.Printf("o:%s e:%s", stdout, stderr)
	}()

	log.Printf("call go do start goroutine")
	time.Sleep(time.Second*2)


	log.Printf("do stop")
	err := c.Stop(time.Second)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	log.Printf("do stop ok")

	time.Sleep(11*time.Second)
}


