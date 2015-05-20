// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package sconf

import (
	"testing"
	"strings"
	"fmt"
)

// 测试格式错误
func TestTierConfErr(t *testing.T) {

	cfg0 := []byte(
		`
a
`)

	tf := NewTierConf()

	err := tf.Load(cfg0)
	fmt.Printf("TestTierConfErr %v\n", err)
	if err == nil {
		t.Errorf("log cfg is err")
	}

	fmt.Println(tf)


}


// 测试基本
func TestTierConf(t *testing.T) {

	cfg0 := []byte(
		`
[log]
LogDir =
LogLevel = TRACE
ConnPort = 00
[empty]
`)

	tf := NewTierConf()

	err := tf.Load(cfg0)
	if err != nil {
		t.Errorf("log cfg err")
	}

	sp, err := tf.StringCheck()
	if err != nil {
		t.Errorf("log cfg err")
	}


	fmt.Println(tf)

	fmt.Println("TestTierConf:", sp)




}

// 测试多层
func TestTierConfMulti(t *testing.T) {

	cfg0 := []byte(
		`
[log]
LogDir =
LogLevel = TRACE
ConnPort = 00
Fuck = you
[empty]
`)

	tf := NewTierConf()

	err := tf.Load(cfg0)
	if err != nil {
		t.Errorf("log cfg err")
	}


	fmt.Println(tf)


	cfg1 := []byte(
		`
[link]
ConnPort = 9988
HeartIntv = 300


[log]
LogDir = ./lua
LogLevel = TRACE
`)


	err = tf.Load(cfg1)
	if err != nil {
		t.Errorf("log cfg err")
	}


	fmt.Println(tf)

	fmt.Println(tf.GetConf())



}


func TestTierConfVar(t *testing.T) {
	fun := "TestTierConfVar"

	cfg0 := []byte(
		`
[log]
LogDir =33
LogLevel = TRACE
ConnPort = 00
[empty]

[lvar]
hh = lvar/${log.ConnPort}

h0 = h0/${lvar.h1}
h1 = h1/${lvar.h2}
h2 = h2/${lvar.h3}
h3 = h3/${lvar.h4}
h4 = h4/${lvar.h0}

[tvar]
cc = abc
c0 = a/${log.LogDir}
c1 = ${log.LogLevel}
c2 = ${  log.LogLevel }/cc
c3 = ${  log.LogLevel}/bb
c4 = ${  log.  LogLevel}
c5 = ${  log  .  LogLevel  }/d
c6 = ${  log  .  LogLevel  }/a/c/${log.LogLevel}/${log.ConnPort}/dd
c7 = ${  log  .  LogLevel  }/a/c/d/${lvar.hh}/${log.ConnPort}
c8 = ${  log  .  LogLevel  }/a/c/d/${lvar.hh}/${log.ConnPort}/a
c9 = ${log.notexist}/a/c/d/${lvar.hh}/${log.ConnPort}/a
c10 = ${log}/a/c/d/${lvar.hh}/${log.ConnPort}/a
c11 = ${ log.notexist.dd}/a/c/d/${lvar.hh}/${log.ConnPort}/a

c12 = ${log.LogLevel}/${tvar.c12}


c13 = aa/c/${lvar.h0}
`)

	tf := NewTierConf()

	err := tf.Load(cfg0)
	if err != nil {
		t.Errorf("log cfg err")
	}

	fmt.Println(fun, tf)

	p, v := "c0", "a/33"
	c, err := tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("hh")
	}


	p, v = "c1", "TRACE"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}


	p, v = "cc", "abc"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}



	p, v = "c2", "TRACE/cc"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}


	p, v = "c3", "TRACE/bb"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}



	p, v = "c4", "TRACE"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}


	p, v = "c5", "TRACE/d"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}


	p, v = "c6", "TRACE/a/c/TRACE/00/dd"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}


	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}



	p, v = "c7", "TRACE/a/c/d/lvar/00/00"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}


	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}


	p, v = "c8", "TRACE/a/c/d/lvar/00/00/a"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}


	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}


	p, v = "c9", "${log.notexist}/a/c/d/lvar/00/00/a"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}


	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}


	p, v = "c10", "${log}/a/c/d/lvar/00/00/a"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}


	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}


	p, v = "c11", "${ log.notexist.dd}/a/c/d/lvar/00/00/a"
	c, err = tf.ToString("tvar", p)
	if err != nil {
		t.Errorf("%s", err)
	}


	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}



	c, err = tf.ToString("tvar", "c12")
	if err == nil {
		t.Errorf("here")
		return
	}
	if strings.Index(err.Error(), "cyclic reference") == -1 {
		t.Errorf("here")
	}

	fmt.Printf("%s,[%s]\n", fun, err)


	p, v = "c12", "ddddd"
	c = tf.ToStringWithDefault("tvar", p, "ddddd")


	fmt.Printf("%s,[%s]\n", fun, c)
	if c != v {
		t.Errorf("here")
	}



	c, err = tf.ToString("tvar", "c13")
	if err == nil {
		t.Errorf("here")
		return
	}

	if strings.Index(err.Error(), "cyclic reference") == -1 {
		t.Errorf("here")
	}

	fmt.Printf("%s,[%s]\n", fun, err)

}

