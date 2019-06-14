// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package snetutil

import (
	"context"
	"testing"
	"net/http"

    "github.com/julienschmidt/httprouter"
	"github.com/shawnfeng/sutil/slog"
)


func TestDo(t *testing.T) {

	slog.Infoln(context.TODO(), "INIT")
	router := httprouter.New()
	router.GET("/test0", HttpRequestWrapper(FactoryReq0))
	router.GET("/test1", HttpRequestWrapper(FactoryReq1))
	router.GET("/test2", HttpRequestWrapper(FactoryReq2))


	http.ListenAndServe(":12345", router)

}



type Req0 struct {
}

func FactoryReq0() HandleRequest {
	r := new(Req0)
	return r
}

func (m *Req0) Handle(r *HttpRequest) HttpResponse {
	fun := "Req0.Handle -->"

	slog.Infof(context.TODO(), "%s url:%s method:%s body:%s", fun, r.URL(), r.Method(), r.Body().Binary())

	return  NewHttpRespString(http.StatusOK, "OK\n")

}





type Req1 struct {
}

func FactoryReq1() HandleRequest {
	r := new(Req1)
	return r
}

func (m *Req1) Handle(r *HttpRequest) HttpResponse {
	fun := "Req1.Handle -->"

	slog.Infof(context.TODO(), "%s url:%s method:%s body:%s", fun, r.URL(), r.Method(), r.Body().Binary())

	h := make(http.Header)
	h.Set("set", "setv")
	h.Add("add", "addv")
	h.Add("add", "addv1")

	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, &http.Cookie{
		Name: "Set0",
		Value: "SetValue0",
	})

	cookies = append(cookies, &http.Cookie{
		Name: "Set1",
		Value: "SetValue1",
	})

	cookies = append(cookies, &http.Cookie{
		Name: "Set2",
		Value: "SetValue2",
	})

	return &HttpRespJson{
		Status: 201,
		Body: map[string]interface{}{
			"name": "value",
		},
		Header: h,
		Cookies: cookies,
	}

}




type Req2 struct {
}

func FactoryReq2() HandleRequest {
	r := new(Req2)
	return r
}

func (m *Req2) Handle(r *HttpRequest) HttpResponse {
	fun := "Req2.Handle -->"

	slog.Infof(context.TODO(), "%s url:%s method:%s body:%s", fun, r.URL(), r.Method(), r.Body().Binary())

	return &HttpRespString{
		Status: 202,
		Body: "req2",

	}

}


