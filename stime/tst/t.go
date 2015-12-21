package main

import (
	"time"
	"log"
)

func main() {
	tm := time.Now()

	log.Println(tm.Zone())
	log.Println(tm)
	log.Println(tm.Location())

	log.Printf("%s %T", time.UTC, time.UTC)

	log.Println(time.LoadLocation("UTC"))
	log.Println(time.LoadLocation("EST"))

	z, offset := tm.Zone()
	log.Println(time.LoadLocation("GMT+2:00"))
	log.Println(z, offset)

	loc := time.FixedZone("UTC", 10)
	log.Println(tm.In(loc))


	loc = time.FixedZone("UTC", 20)
	log.Println(tm.In(loc))


	loc = time.FixedZone("UTC", -10)
	log.Println(tm.In(loc))

	loc = time.FixedZone("UTC", 0)
	tm2 := time.Unix(1449662649, 0)
	log.Println(tm2)
	log.Println(tm2.In(loc))

}
