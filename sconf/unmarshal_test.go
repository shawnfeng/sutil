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
M.k0=a,b,c
M.k1=c,d,e
M.k2=e
M.=ddd
M=asd

[SM.a]
EE=ee
FF=ff

[SM.b]
EE=ff
FF=ee

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
			M map[string][]string
		}

		Ts1 *string

		Sm map[string]struct {
			Ee string
			Ff string
		}

		
	}


	var c2 TConf2
	fmt.Println("OK=======================")
	err = tf.Unmarshal(&c2)
	if err != nil {
		t.Errorf("err:%s", err)
	}
	fmt.Println("result:", c2, c2.Ts, c2.Sm, err)


	cfg2 := []byte(
		`
[ffff]
ZZZ=1
YYY=33
AAA=b
MMM.k0=3
MMM.k1=4
MMM.k2=5
MMM=
[uame]
[noexist]
[ts]
AAA=a
CCC=false
XXX=ddd
LLL=1,2,3,4,5
[SM.a]
EE=ee
FF=ff

[SM.b]
EE=ff
FF=ee

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
		MMM map[string]int `sep:"."`
	}
	err = tf2.Unmarshal(&c5)
	if err != nil {
		t.Errorf("error here:%s", err)
	}
	fmt.Println("result:", c5, err)



}
