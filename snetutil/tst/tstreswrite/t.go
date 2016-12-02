package main

import (
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sendstrailers", func(w http.ResponseWriter, req *http.Request) {
		// Before any call to WriteHeader or Write, declare
		// the trailers you will set during the HTTP
		// response. These three headers are actually sent in
		// the trailer.
		w.Header().Set("Trailer", "AtEnd1, AtEnd2")
		w.Header().Add("Trailer", "AtEnd3")

		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header

		http.SetCookie(w, &http.Cookie{
			Name: "Set0",
			Value: "SetValue0",
		})


		http.SetCookie(w, &http.Cookie{
			Name: "Set1",
			Value: "SetValue1",
		})

		w.WriteHeader(http.StatusOK)

		w.Header().Set("AtEnd1", "value 1")
		http.SetCookie(w, &http.Cookie{
			Name: "Set2",
			Value: "SetValue2",
		})

		io.WriteString(w, "This HTTP response has both headers before this text and trailers at the end.\n")
		w.Header().Set("AtEnd2", "value 2")
		w.Header().Set("AtEnd3", "value 3") // These will appear as trailers.


		http.SetCookie(w, &http.Cookie{
			Name: "Set3",
			Value: "SetValue3",
		})
	})


	http.ListenAndServe(":12345", mux)
}

// WriteHeader 必须在io.WriteString之前调用在起作用否则会出错误: http: multiple response.WriteHeader calls
// w.Header 必须在WriteHeader，io.WriteString之前调用，否则不起作用
// SetCookie 必须在WriteHeader，io.WriteString之前调用，否则不起作用
