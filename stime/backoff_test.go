package stime

import (
	"testing"
	"time"
	"log"
	"fmt"
)

func TestDayBeginStamp(t *testing.T) {

	zone, offset := time.Now().Zone()
	fmt.Println(zone, offset)


	//now := time.Now().Unix()

	// CST 2015-04-11
	tmbegin := int64(1428681600)
	log.Println(time.Unix(tmbegin, 0).Format("2006-01-02/15:04:05"))
	for i := 0; i < 3600*24*5; i++ {
		tm := int64(tmbegin+int64(i))
		stmd := time.Unix(tm, 0).Format("2006-01-02")
		stm := fmt.Sprintf("%s/00:00:00", stmd)



		begin := DayBeginStamp(tm)

		sbegin := time.Unix(begin, 0).Format("2006-01-02/15:04:05")

		if sbegin != stm {
			t.Errorf("error %d %s", tm, sbegin)
		}

		//log.Println(time.Unix(tm, 0).Format("2006-01-02/15:04:05"), stm, begin, sbegin)
	}


}

func TestBackoffReset(t *testing.T) {
	log.Println("BackOffRest Begin")
	bo := NewBackOffCtrl(time.Second * 1, time.Second*10)

	go func() {
		bg := time.Now().Unix()
		bo.BackOff() // 0s
		log.Printf("BackOffRest routine %d", time.Now().Unix() - bg)
		bo.BackOff() // 1s
		log.Printf("BackOffRest routine %d", time.Now().Unix() - bg)
		bo.BackOff() // 2s
		log.Printf("BackOffRest routine %d", time.Now().Unix() - bg)
		bo.BackOff() // 1s
		log.Printf("BackOffRest routine %d", time.Now().Unix() - bg)
		if time.Now().Unix() - bg != 4 {
			t.Errorf("BackOffRest reset err")
		} else {
			log.Println("BackOffRest OK Reset")
		}
	}()

	time.Sleep(time.Second * time.Duration(4))


	log.Printf("breset %v", bo)
	bo.Reset()

	log.Printf("areset %v", bo)
	bg := time.Now().Unix()
	bo.BackOff() // 02
	log.Printf("BackOffRest b1 %d", time.Now().Unix() - bg)
	bo.BackOff() // 1s
	log.Printf("BackOffRest b2 %d", time.Now().Unix() - bg)

	ttt := time.Now().Unix() - bg
	if ttt != 1 {
		t.Errorf("BackOffRest reset Continue err:%d", ttt)
	} else {
		log.Println("BackOffRest OK Reset Continue")
	}


	log.Printf("set before reset %v", bo)
	bo.SetCtrl(time.Second * 2, time.Second*5)
	log.Printf("set end areset %v", bo)

}


func TestBackoff(t *testing.T) {
	bo := NewBackOffCtrl(time.Second * 1, time.Second*10)

	for i := uint32(0); i < 8; i++ {
		log.Printf("BackOff %d Begin", i)
		bg := time.Now().Unix()

		if i == 0 {
			bo.BackOff()
		}
		bo.BackOff()
		intv := time.Now().Unix() - bg

		if intv > 10 {
			t.Errorf("BackOff ceil err %d", intv)
		}

		if intv != 1 << i {
			if 1 << i > 10 && intv == 10 {
				log.Printf("BackOff ceil %d", i)
			} else {
				t.Errorf("BackOff time %d err", i)
			}
		}

		log.Printf("BackOff %d End", i)

	}


}

