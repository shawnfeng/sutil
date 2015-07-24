// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package snetutil


import (
	"fmt"
	"bytes"
	"strconv"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/json"
    "github.com/julienschmidt/httprouter"

	"github.com/shawnfeng/sutil/slog"
)

// http response interface
type HttpResponse interface {
	Marshal() (int, []byte)
}

// 定义了几种产用的类型的response

// json形式的response
type HttpRespJson struct {
	status int
	resp interface{}
}


func (m *HttpRespJson) Marshal() (int, []byte) {
	fun := "HttpRespJson.Marshal -->"
	resp, err := json.Marshal(m.resp)

	if err != nil {
		slog.Warnf("%s json unmarshal err:%s", fun, err)
	}

	return m.status, resp
}


func NewHttpRespJson200(r interface{}) HttpResponse {
	return &HttpRespJson{200, r}
}


func NewHttpRespJson(status int, r interface{}) HttpResponse {
	return &HttpRespJson{status, r}
}


// byte 形式的response
type HttpRespBytes struct {
	status int
	resp []byte
}

func (m *HttpRespBytes) Marshal() (int, []byte) {
	return m.status, m.resp
}

func NewHttpRespBytes(status int, resp []byte) HttpResponse {
	return &HttpRespBytes{status, resp}
}


// string 形式的response
type HttpRespString struct {
	status int
	resp string
}

func (m *HttpRespString) Marshal() (int, []byte) {
	return m.status, []byte(m.resp)
}

func NewHttpRespString(status int, resp string) HttpResponse {
	return &HttpRespString{status, resp}
}

// ===============================================
type keyGet interface {
	Get(key string) string
}


type reqArgs struct {
	r keyGet
}

func NewreqArgs(r keyGet) *reqArgs {
	return &reqArgs{r}
}


func (m *reqArgs) String(key string) string {
	return m.r.Get(key)
}

func (m *reqArgs) Int(key string) int {
	fun := "reqArgs.Int -->"
	v := m.r.Get(key)
	if len(v) == 0 {
		return 0
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		slog.Warnf("%s parse int v:%s err:%s", fun, v, err)
	}
	return i

}

func (m *reqArgs) Int32(key string) int32 {
	fun := "reqArgs.Int32 -->"
	v := m.r.Get(key)
	if len(v) == 0 {
		return 0
	}


	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		slog.Warnf("%s parse int32 v:%s err:%s", fun, v, err)
	}

	return int32(i)
}

func (m *reqArgs) Int64(key string) int64 {
	fun := "reqArgs.Int64 -->"
	v := m.r.Get(key)
	if len(v) == 0 {
		return 0
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		slog.Warnf("%s parse int64 v:%s err:%s", fun, v, err)
	}
	return i

}
// =========================
type reqQuery struct {
	r *http.Request
	q url.Values
}

func (m *reqQuery) Get(key string) string {
	fun := "reqQuery.Get -->"
	if m.q == nil {
		if m.r.URL != nil {
			var err error
			m.q, err = url.ParseQuery(m.r.URL.RawQuery)
			if err != nil {
				slog.Warnf("%s parse query q:%s err:%s", fun, m.r.URL.RawQuery, err)
			}
		}

		if m.q == nil {
			m.q = make(url.Values)
		}

		slog.Debugf("%s parse query q:%s err:%s", fun, m.r.URL.RawQuery, m.q)
	}


	return m.q.Get(key)
}


type reqParams struct {
	p httprouter.Params
}

func (m *reqParams) Get(key string) string {
	return m.p.ByName(key)
}

// ==========
type reqBody struct {
	body []byte
}

func (m *reqBody) Body() []byte {
	return m.body
}

func NewreqBody(body []byte) *reqBody {
	return &reqBody {
		body: body,
	}

}


// ============================
// 没有body类的请求
type HttpRequestNoBody struct {
	URL *url.URL
	Method string
	Query *reqArgs
	Params *reqArgs

}


func NewHttpRequestNoBody(r *http.Request, ps httprouter.Params) (*HttpRequestNoBody, error) {
	return &HttpRequestNoBody {
		URL: r.URL,
		Method: r.Method,
		Query: NewreqArgs(&reqQuery{r: r,}),
		Params: NewreqArgs(&reqParams{ps}),
	}, nil
}


type HttpRequestCommonBody struct {
	HttpRequestNoBody
	Body *reqBody
}


func NewHttpRequestCommonBody(r *http.Request, ps httprouter.Params) (*HttpRequestCommonBody, error) {

	body, err := ioutil.ReadAll(r.Body);
	if err != nil {
		return nil, fmt.Errorf("read body %s", err.Error())
	}

	return &HttpRequestCommonBody {
		HttpRequestNoBody: HttpRequestNoBody{r.URL, r.Method, NewreqArgs(&reqQuery{r: r}), NewreqArgs(&reqParams{ps})},
		Body: NewreqBody(body),
	}, nil
}



func NewHttpRequestJsonBody(r *http.Request, ps httprouter.Params, js interface{}) (*HttpRequestCommonBody, error) {
	hrb, err := NewHttpRequestCommonBody(r, ps)
	if err != nil {
		return hrb, err
	}

    dc := json.NewDecoder(bytes.NewBuffer(hrb.Body.Body()))
    dc.UseNumber()
    err = dc.Decode(js)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal %s", err.Error())
	}


	return hrb, nil

}


type HandleNoBody interface {
	Handle(*HttpRequestNoBody) HttpResponse
	// 构造自己一个副本，如果结构本身保存了
	// 本请求的数据，则factory必须new一个新的
	// 否则没有共享数据问题，可以返回自己当前的指针就好
	Factory() HandleNoBody
}

type HandleCommonBody interface {
	Handle(*HttpRequestCommonBody) HttpResponse

	// 对于存在body，并使用json自动unmarshal
	// 一定注意使用一般时候你都需要new一个新的
	// 除非你想让请求之间通过某种技巧来关联，否则。。。
	Factory() HandleCommonBody
}

func HttpNoBodyWrapper(h HandleNoBody) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	fun := "HttpNoBodyWrapper -->"

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req, err := NewHttpRequestNoBody(r, ps)
		if err != nil {
			slog.Warnf("%s body json err:%s", fun, err)
			http.Error(w, "request err", 400)
			return
		}

		resp := h.Factory().Handle(req)
		status, rs := resp.Marshal()

		if status == 200 {
			fmt.Fprintf(w, "%s", rs)
		} else {
			http.Error(w, string(rs), status)
		}
	}

}



func HttpCommonBodyWrapper(h HandleCommonBody) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	fun := "HttpCommmonBodyWrapper -->"

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req, err := NewHttpRequestCommonBody(r, ps)
		if err != nil {
			slog.Warnf("%s body json err:%s", fun, err)
			http.Error(w, "request err", 400)
			return
		}

		resp := h.Factory().Handle(req)
		status, rs := resp.Marshal()

		if status == 200 {
			fmt.Fprintf(w, "%s", rs)
		} else {
			http.Error(w, string(rs), status)
		}
	}

}



func HttpJsonBodyWrapper(h HandleCommonBody) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	fun := "HttpJsonBodyWrapper -->"

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		newme := h.Factory()
		req, err := NewHttpRequestJsonBody(r, ps, newme)
		if err != nil {
			slog.Warnf("%s body json err:%s", fun, err)
			http.Error(w, "request err", 400)
			return
		}

		resp := newme.Handle(req)
		status, rs := resp.Marshal()

		if status == 200 {
			fmt.Fprintf(w, "%s", rs)
		} else {
			http.Error(w, string(rs), status)
		}
	}

}
