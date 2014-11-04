package stime

import (
	"testing"

	"time"
	"log"
)

func TestStat(t *testing.T) {
	st := NewTimeStat("Test")

	log.Println(st.StatLog())

	log.Println(st.Millisecond())
	log.Println(st.Microsecond())
	log.Println(st.Nanosecond())

}



func TestStat2(t *testing.T) {
	st := NewTimeStat("Test")

	time.Sleep(time.Millisecond * time.Duration(2))

	if st.Millisecond() != 2 {
		log.Println("time stat error")
	}

	log.Println(st.Millisecond())
	log.Println(st.Microsecond())
	log.Println(st.Nanosecond())


}
