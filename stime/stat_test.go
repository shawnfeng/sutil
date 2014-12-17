package stime

import (
	"testing"

	"time"
	"log"
)

func TestStat(t *testing.T) {
	st := NewTimeStat()

	log.Println(st.Millisecond())
	log.Println(st.Microsecond())
	d := st.Duration()
	log.Printf("%d\n", d)
	log.Println(d, int64(d), d.Nanoseconds(), d.Seconds(), d.Minutes(), d.Hours())

	log.Println(st.Nanosecond())




}



func TestStat2(t *testing.T) {
	st := NewTimeStat()

	time.Sleep(time.Millisecond * time.Duration(2))

	if st.Millisecond() != 2 {
		log.Println("time stat error")
		t.Errorf("time stat error")
	}

	log.Println(st.Millisecond())
	log.Println(st.Microsecond())
	d := st.Duration()
	log.Println(d)

	log.Println(st.Nanosecond())


	st.Reset()
	time.Sleep(time.Millisecond * time.Duration(1))
	log.Println(st.Duration())

	log.Println(d-st.Duration())
	if st.Millisecond() != 1 {
		log.Println("time stat error")
		t.Errorf("time stat error")
	}

}
