package sutil

import (
	"fmt"
	"testing"

)


func TestWriteFile(t *testing.T) {

	err := WriteFile("aa", []byte("abcde\n"), 0600)

	if err != nil {
		t.Errorf("%s", err)
	}


	err = WriteFile("log/aa", []byte("abcde\n"), 0600)

	if err != nil {
		t.Errorf("%s", err)
	}


	err = WriteFile("log/log/aa", []byte("abcde\n"), 0600)

	if err != nil {
		t.Errorf("%s", err)
	}


	err = WriteFile("log/", []byte("abcde\n"), 0600)
	if err == nil {
		t.Errorf("%s", err)
	}

}



func TestVersion(t *testing.T) {

	v := NewVersionCmp("1.2.3")

	//fmt.Println(v.fmtver("1.2.3-beta"))
	fmt.Println(v.fmtver(""))

	if v.Lt("1.2.3") {
		t.Errorf("hhh")
	}

	if !v.Lte("1.2.3") {
		t.Errorf("hhh")
	}

	if !v.Gte("1.2.3") {
		t.Errorf("hhh")
	}


	if v.Ne("1.2.3") {
		t.Errorf("hhh")
	}

	if !v.Eq("1.2.3") {
		t.Errorf("hhh")
	}



	if v.Lt("1.1.3") {
		t.Errorf("hhh")
	}

	if !v.Gt("1.1.3") {
		t.Errorf("hhh")
	}

	if !v.Lt("2.0.1") {
		t.Errorf("hhh")
	}

	if v.Gt("2.0.1") {
		t.Errorf("hhh")
	}



	if !v.Lt("1.2.3.1") {
		t.Errorf("hhh")
	}


	if v.Gt("1.2.3.1") {
		t.Errorf("hhh")
	}


	if v.Lt("1.2.2") {
		t.Errorf("hhh")
	}


	if !v.Gt("1.2.2") {
		t.Errorf("hhh")
	}


	if v.Lt("1.2.2.9.9") {
		t.Errorf("hhh")
	}


	if !v.Gt("1.2.2.9.9") {
		t.Errorf("hhh")
	}


	if !v.Lt("1.10.3") {
		t.Errorf("hhh")
	}


	if v.Gt("1.10.3") {
		t.Errorf("hhh")
	}


	if !v.Lt("10.10.3") {
		t.Errorf("hhh")
	}


	if v.Gt("10.10.3") {
		t.Errorf("hhh")
	}


	fmt.Println(v.fmtver(v.Min()))
	fmt.Println(v.fmtver(v.Max()))
	fmt.Println(v.fmtver("1.2.3"))
	if v.Gt(v.Max()) {
		t.Errorf("hhh")
	}


	if v.Lt(v.Min()) {
		t.Errorf("hhh")
	}


}
