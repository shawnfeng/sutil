package dbrouter

import (
	"errors"
	"testing"
	"time"
)

func TestEntry(t *testing.T) {
	go func() {
		for {
			Entry("test")
			time.Sleep(time.Millisecond * 1000)
		}
	}()

	go func() {
		for i := 0; i < 20; i++ {
			statBreaker("test", errors.New("timeout"))
		}
	}()
	select {}
}
