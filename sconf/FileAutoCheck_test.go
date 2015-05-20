// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package sconf

import (
	"testing"
	"time"
	"fmt"
)


func Test_01(t *testing.T) {
	return

	fc := NewFileAutoCheck("./a.check")

	fmt.Println(fc)


	needup, data, err := fc.Check()
	fmt.Println(needup, string(data), err)
	fmt.Println(fc)
	time.Sleep(time.Second * time.Duration(10))
	needup, data, err = fc.Check()
	fmt.Println(needup, string(data), err)

	fmt.Println(fc)

}
