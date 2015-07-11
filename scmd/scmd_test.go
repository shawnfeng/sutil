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
	testStart0(t)
}


func testStart0(t *testing.T) {
	
	//c := NewScmd("grep", "A", "/home/fenggx/shawn_go/src/github.com/shawnfeng/sutil/scmd/*")  // * 执行是有错误的，因为  *  通配符是shell支持的
	c := NewScmd("grep", "A", "/home/fenggx/shawn_go/src/github.com/shawnfeng/sutil/scmd/scmd.go")

	//c := NewScmd("echo", "AAAA", "BBBB")

	stdout, stderr, err := c.StartWaitOutput()
	if err != nil {
		t.Errorf("%s", err)
	}


	log.Printf("o:%s e:%s", stdout, stderr)


}



func testStart1(t *testing.T) {

	c := NewScmd("sleep", "2")

	stdout, stderr, err := c.StartWaitOutput()
	if err != nil {
		t.Errorf("%s", err)
	}


	log.Printf("o:%s e:%s", stdout, stderr)


}



func testStartStop(t *testing.T) {

	c := NewScmd("sleep", "10")

	go func() {
		stdout, stderr, err := c.StartWaitOutput()
		if err != nil {
			t.Errorf("%s", err)
		}
		log.Printf("o:%s e:%s", stdout, stderr)
	}()

	time.Sleep(time.Second)

	err := c.Stop()
	if err != nil {
		t.Errorf("%s", err)
	}

}


