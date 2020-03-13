package dbrouter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
)

func TestMongoTimeout(t *testing.T) {
	ctx := context.TODO()
	router, err := NewRouter(nil)
	if err != nil {
		t.Errorf("init router failed, err: %v", err)
	}
	query := func(c *mgo.Collection) error {
		n, err := c.Count()
		fmt.Println(n)
		return err
	}
	for i := 0; i < 16; i++ {
		go func() {
			for {
				err = router.MongoExecEventual(ctx, "STAT", "dbrouter_oprecord", query)
				if err != nil {
					fmt.Printf("router exec failed, %v", err)
				}
				time.Sleep(time.Millisecond * 50)
			}
		}()

	}
	select{}
}
