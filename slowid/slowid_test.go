// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package slowid


import (
	"time"
	"log"
	"testing"

)



func TestDo(t *testing.T) {

	log.Println(MaxWorkerId, MaxSequence)
	log.Println(time.Now().Zone())

	sf, _ := NewSlowid(0)

	for i := 0; i < 10; i++ {
		nid, _ := sf.Next()
		log.Println(nid, nid >> 11, Since+nid>>11, time.Now().UnixNano()/1000000, (time.Now().UnixNano()/1000000 - Since) << 11)
	}


	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond)
		nid, _ := sf.Next()




		log.Println(nid, nid >> 11)
	}

}
