// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package ssync


import (
	"fmt"
	"time"
	"testing"

	"sync"
)


func testsMu(t *testing.T) {
	var mm sync.Mutex
	mm.Unlock()
	//mm.Lock()
	//mm.Lock()

}

func testMu1(t *testing.T) {

	var mu Mutex

	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("do lock", i)
			mu.Lock()
			fmt.Println("do lock ok", i)
		}(i)
	}

	//mu.Unlock()

	isl := mu.Trylock()
	fmt.Println("trylock", isl)


	time.Sleep(time.Second)
	mu.Unlock()
	fmt.Println("unlock")
	time.Sleep(time.Second)


}


func TestMu2(t *testing.T) {

	var mu Mutex

	isl := mu.Trylock()
	fmt.Println("trylock", isl)

	isl = mu.Trylock()
	fmt.Println("trylock", isl)

	isl = mu.Trylock()
	fmt.Println("trylock", isl)

	isl = mu.Trylock()
	fmt.Println("trylock", isl)


	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println("do lock", i)
			mu.Lock()
			fmt.Println("do lock ok", i)
		}(i)
	}

	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println("do unlock", i)
			mu.Unlock()
			fmt.Println("do unlock ok", i)
		}(i)
	}

	time.Sleep(time.Second)

}



func testMu3(t *testing.T) {

	var once sync.Once
	onceBody := func() {
		fmt.Println("Only once")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
			fmt.Println("loop")
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}



}
