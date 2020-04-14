package dbrouter

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestQueryTimeout(t *testing.T) {
	ctx := context.TODO()
	router, err := NewRouter(nil)
	if err != nil {
		t.Errorf("init router failed, err: %v", err)
	}
	for i:=0;i<16;i++ {
		go func() {
			for {
				err = router.SqlExec(ctx, "COURSEWAREX", func(db *DB, tables []interface{}) (err error) {
					_, err = db.Query("select sleep(10)")
					return err
				}, "coursewarex_1")
				if err != nil {
					fmt.Printf("router exec failed, %v", err)
				}
				time.Sleep(time.Millisecond * 50)
			}
		}()
	}
	select{}
}
