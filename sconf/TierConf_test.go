package sconf

import (
	"testing"

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
	if err == nil {
		t.Errorf("log cfg is err")
	}

	fmt.Println(tf)

	stf := fmt.Sprintf("%s", tf)
	if stf !=
		"&{map[]}" {
		t.Errorf("error:%s", stf)
	}




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

	fmt.Println(tf)


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

