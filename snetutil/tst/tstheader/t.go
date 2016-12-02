package main

import (
	"net/http"
	"log"
)

func main() {
	h := make(http.Header)
	h.Set("set", "setv")
	h.Add("add", "addv")
	h.Add("add", "addv")
	log.Println(h)
}


