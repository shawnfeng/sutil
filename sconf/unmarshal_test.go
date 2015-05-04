package sconf

import (
	"testing"
	"fmt"
)


type TConf struct {
    Uame   string "user name"
    Passwd string "user passsword"
	Fuck int `sconf:"ffff"`
	Girl int64

	// 不是指针的、是指针的，指针为空的或者不为空的
	Ts *struct {
		AAA string
		BBB int
	}
}


func TestUnmar(t *testing.T) {
	cfg0 := []byte(
		`
[ffff]
[uame]
[noexist]
[ts]
AAA=a
CCC=true
BBB=233
LLL=1,2,3
`)

	tf := NewTierConf()

	err := tf.Load(cfg0)
	if err != nil {
		t.Errorf("log cfg err")
	}


	var c TConf
	fmt.Println(c)
	err = tf.Unmarshal(c)
	fmt.Println(c, err)

	tt := 0
	err = tf.Unmarshal(&tt)
	fmt.Println(c, err)



	err = tf.Unmarshal(&c)
	fmt.Println("result:", c, err)


	type TConf2 struct {
		Uname   string
		Passwd string
		Fuck int
		Girl int64

		// 不是指针的、是指针的，指针为空的或者不为空的
		Ts *struct {
			AAA string
			BBB uint8

			CCC bool

			LLL []int `sep:"," sconf:"lll"`
		}

		Ts1 *string
	}


	var c2 TConf2
	fmt.Println("OK=======================")
	err = tf.Unmarshal(&c2)
	fmt.Println("result:", c2, c2.Ts, err)


	cfg2 := []byte(
		`
[ffff]
ZZZ=1
YYY=33
AAA=b
[uame]
[noexist]
[ts]
AAA=a
CCC=false
XXX=ddd
LLL=1,2,3,4,5
`)


	tf2 := NewTierConf()

	err = tf2.Load(cfg2)
	if err != nil {
		t.Errorf("log cfg err")
	}


	err = tf2.Unmarshal(&c2)
	fmt.Println("result2:", c2, c2.Ts, err)



	var c3 map[int]string
	err = tf2.Unmarshal(&c3)
	if err == nil {
		t.Errorf("error here")
	}
	fmt.Println("result:", c3, err)


	//======
	var c4 map[string]string
	err = tf2.Unmarshal(&c4)
	if err == nil {
		t.Errorf("error here")
	}
	fmt.Println("result:", c4, err)


	//======
	var c5 map[string] struct {
		YYY int
		ZZZ string
		AAA string
	}
	err = tf2.Unmarshal(&c5)
	if err != nil {
		t.Errorf("error here:%s", err)
	}
	fmt.Println("result:", c5, err)



}
