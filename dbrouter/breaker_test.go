package dbrouter

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestEntry(t *testing.T) {
	go func() {
		for {
			fmt.Printf("Entry: %v\n", Entry("group", "test"))
			time.Sleep(time.Millisecond * 1000)
		}
	}()

	go func() {
		for i := 0; i < 20; i++ {
			statBreaker("group", "test", errors.New("timeout"))
		}
	}()
	select {}
}
