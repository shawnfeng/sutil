package main

import (
	"encoding/json"
	"log"
)

type tinter map[string]string


func (m *tinter) load() {
	*m = make(map[string]string)
	(*m)["aa"] = "bb"
}


type tinter0 struct {
	a string
	b int
}


func (m *tinter0) load() {
	*m = tinter0{"aa", 11}
}


func ttt1() {

	var t tinter
	if &t == nil {
		// not print
		log.Println("nil 0")
	}
	log.Printf("& %s %T", &t, &t)

	st, err := json.Marshal(&t)
	// null|<nil>
	log.Printf("%s|%v", st, err)

	var tt tinter
	err = json.Unmarshal(st, &tt)
	log.Printf("un %s|%v", tt, err)
	if tt == nil {
		log.Println("nil 0.1")
	}

	log.Println(t)
	if t == nil {
		log.Println("nil 1")
	}

	t.load()
	if t == nil {
		log.Println("nil 2")
	}

	log.Println(t)
}



func ttt2() {

	var t tinter0
	if &t == nil {
		// not print
		log.Println("nil 0")
	}
	log.Printf("& %s %T", &t, &t)

	st, err := json.Marshal(&t)
	// {}|<nil>
	log.Printf("%s|%v", st, err)

	var tt tinter0
	err = json.Unmarshal(st, &tt)
	log.Printf("un %s|%v", tt, err)


	log.Println(t)

	t.load()


	log.Println(t)
}


func ttt3() {

	var t *tinter0
	if t == nil {
		// not print
		log.Println("nil 0")
	}
	log.Printf("& %s %T", &t, &t)

	st, err := json.Marshal(t)
	// null|<nil>
	log.Printf("%s|%v", st, err)

	var tt *tinter0
	err = json.Unmarshal(st, tt)
	log.Printf("un %s|%v", tt, err)


	log.Println(t)

	//t.load()

}


func main() {

	ttt3()
}
