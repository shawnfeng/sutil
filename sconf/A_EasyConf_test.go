package sconf

import (
	"testing"
	"strings"

	"fmt"
)


func TestToInt(t *testing.T) {
	ef := NewTierConf()

	c := make(map[string]map[string]string)

	pv := map[string]string {
		"p0": "v0",
		"p1": "v1",
		"p2": "v2",
		"p3": "1234",

	}

	pv2 := map[string]string {

		"p4": "a,  bc,\t   de   ",
		"p5": "ada;afew;  qwe",

		"p6": "1;2;q;3",
		"p7": "-1;2000;134123;3",
	}

	c["s0"] = pv
	ef.LoadFromConf(c)
	c["s1"] = pv
	ef.LoadFromConf(c)
	c["s0"] = pv2
	c["s1"] = pv2
	c["s2"] = pv2
	ef.LoadFromConf(c)

	v, err := ef.ToInt("nsection", "aa")
	fmt.Println(v, err)
	if strings.Index(err.Error(), "section empty") == -1 {
		t.Errorf("section")
	}


	v, err = ef.ToInt("s0", "npro")
	fmt.Println(v, err)
	if strings.Index(err.Error(), "property empty") == -1 {
		t.Errorf("property")
	}



	v, err = ef.ToInt("s0", "p0")
	fmt.Println(v, err)
	if err == nil {
		t.Errorf("convert")
	}


	v, err = ef.ToInt("s1", "p3")
	fmt.Println(v, err)
	if err != nil {
		t.Errorf("err")
	}


	v = ef.ToIntWithDefault("s0", "npro", 33)
	fmt.Println(v)


	s, err := ef.ToString("nsection", "aa")
	fmt.Println(s, err)
	if strings.Index(err.Error(), "section empty") == -1 {
		t.Errorf("section")
	}


	s, err = ef.ToString("s2", "p2")
	fmt.Println(s, err)


	s = ef.ToStringWithDefault("nsection", "aa", "def")
	fmt.Println(s)



	ss, err := ef.ToSliceString("nsecion", "aa", ",")
	fmt.Println(ss, err)
	if strings.Index(err.Error(), "section empty") == -1 {
		t.Errorf("section")
	}



	ss, err = ef.ToSliceString("s0", "p4", ",")
	fmt.Println(ss, err)

	ss, err = ef.ToSliceString("s0", "p4", ";")
	fmt.Println(ss, err)


	ss, err = ef.ToSliceString("s0", "p5", ",")
	fmt.Println(ss, err)


	ss, err = ef.ToSliceString("s0", "p5", ";")
	fmt.Println(ss, err)

	is, err := ef.ToSliceInt("nsecion", "aa", ",")
	fmt.Println(is, err)
	if strings.Index(err.Error(), "section empty") == -1 {
		t.Errorf("section")
	}


	is, err = ef.ToSliceInt("s1", "p6", ";")
	fmt.Println(is, err)


	is, err = ef.ToSliceInt("s1", "p7", ";")
	fmt.Println(is, err)


}
