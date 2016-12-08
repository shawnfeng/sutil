// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package snetutil


import (
	"fmt"
	"bytes"
	"strings"
	"strconv"
	"io/ioutil"
	"net/http"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"encoding/json"
    "github.com/julienschmidt/httprouter"

	"github.com/shawnfeng/sutil/slog"
)

func DoWriteResponse(w http.ResponseWriter, header http.Header, cookies []*http.Cookie, status int, body io.Reader) {
	fun := "DoWriteResponse -->"

	h := w.Header()
	for k, v := range header {
		h[k] = v
	}

	for _, c := range cookies {
		http.SetCookie(w, c)
	}

	w.WriteHeader(status)

	n, err := io.Copy(w, body)
	if err != nil {
		slog.Warnf("%s white response n:%d err:%s", fun, n, err)
	}

// WriteHeader 必须在Copy之前调用在起作用否则会出错误: http: multiple response.WriteHeader calls
// w.Header 必须在WriteHeader，Copy之前调用，否则不起作用
// SetCookie 必须在WriteHeader，Copy之前调用，否则不起作用

}

// http response interface
type HttpResponse interface {
	WriteResponse(http.ResponseWriter)
}

// 定义了几种产用的类型的response

// json形式的response
type HttpRespJson struct {
	Status int
	Body interface{}
	Header http.Header
	Cookies []*http.Cookie
}


func (m *HttpRespJson) WriteResponse(w http.ResponseWriter) {
	fun := "HttpRespJson.Marshal -->"
	resp, err := json.Marshal(m.Body)

	if err != nil {
		slog.Warnf("%s json unmarshal err:%s", fun, err)
		DoWriteResponse(w, m.Header, m.Cookies, m.Status, strings.NewReader(err.Error()))

	} else {
		DoWriteResponse(w, m.Header, m.Cookies, m.Status, bytes.NewReader(resp))

	}

}


func NewHttpRespJson200(body interface{}) HttpResponse {
	return &HttpRespJson{http.StatusOK, body, nil, nil}
}


func NewHttpRespJson(status int, body interface{}) HttpResponse {
	return &HttpRespJson{status, body, nil, nil}
}


// byte 形式的response
type HttpRespBytes struct {
	Status int
	Body []byte

	Header http.Header
	Cookies []*http.Cookie
}

func (m *HttpRespBytes) WriteResponse(w http.ResponseWriter) {
	DoWriteResponse(w, m.Header, m.Cookies, m.Status, bytes.NewReader(m.Body))

}

func NewHttpRespBytes(status int, body []byte) HttpResponse {
	return &HttpRespBytes{status, body, nil, nil}
}


// string 形式的response
type HttpRespString struct {
	Status int
	Body string

	Header http.Header
	Cookies []*http.Cookie
}

func (m *HttpRespString) WriteResponse(w http.ResponseWriter) {
	DoWriteResponse(w, m.Header, m.Cookies, m.Status, strings.NewReader(m.Body))

}

func NewHttpRespString(status int, body string) HttpResponse {
	return &HttpRespString{status, body, nil, nil}
}


// redirect
type HttpRespRedirect struct {
	Status int
	Url string
	R *HttpRequest
}

func NewHttpRespRedirect(r *HttpRequest, status int, redirectUrl string) HttpResponse {
	return &HttpRespRedirect{status, redirectUrl, r}
}


func (m *HttpRespRedirect) WriteResponse(w http.ResponseWriter) {
	http.Redirect(w, m.R.Request(), m.Url, m.Status)
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


func (m *reqArgs) Bool(key string) bool {
	fun := "reqArgs.Bool -->"
	v := m.r.Get(key)
	if len(v) == 0 {
		return false
	}

	i, err := strconv.ParseBool(v)
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
	r *http.Request
	body []byte
}

func (m *reqBody) Binary() []byte {
	fun := "reqBody.Binary"
	if m.body == nil {
		body, err := ioutil.ReadAll(m.r.Body);
		if err != nil {
			slog.Errorf("%s read body %s", fun, err.Error())
		}
		m.body = body
	}

	return m.body
}

// https://golang.org/pkg/net/http/#Request
// For server requests the Request Body is always non-nil
// but will return EOF immediately when no body is present.
// The Server will close the request body. The ServeHTTP
// Handler does not need to.
func (m *reqBody) Reader() io.ReadCloser {
	return m.r.Body
}



func (m *reqBody) Json(js interface{}) error {

    dc := json.NewDecoder(bytes.NewBuffer(m.Binary()))
    dc.UseNumber()
    err := dc.Decode(js)
	if err != nil {
		return fmt.Errorf("json unmarshal %s", err.Error())
	} else {
		return nil
	}

}


func (m *reqBody) FormValue(key string) string {
	fun := "reqBody.FormValue -->"
	// 获取到content-type，并根据其类型来决策是从r.MultipartForm，获取数据
	// 还是r.PostForm中获取数据，r.Form实际上市把query中的postform中的，mutlpartform都搞到一起了
	// r.PostFrom 对应的content-type为 application/x-www-form-urlencoded
	// r.MultipartForm 对应的 multipart/form-data

	// 仅仅是为让内部触发对form的parse过程
	m.r.FormValue(key)


	// 参照http package中parsePostForm 实现
	ct := m.r.Header.Get("Content-Type")
	// RFC 2616, section 7.2.1 - empty type
	//   SHOULD be treated as application/octet-stream
	if ct == "" {
		ct = "application/octet-stream"
	}
	var err error
	ct, _, err = mime.ParseMediaType(ct)
	if err != nil {
		slog.Errorf("%s parsemediatype err:%s", fun, err)
	}


	if ct == "application/x-www-form-urlencoded" {
		if vs := m.r.PostForm[key]; len(vs) > 0 {
			return vs[0]
		}

	} else if ct == "multipart/form-data" {
		if m.r.MultipartForm == nil {
			slog.Errorf("%s multipart/form-data parse nil", fun)
			return ""
		}


		if vs := m.r.MultipartForm.Value[key]; len(vs) > 0 {
			return vs[0]
		}

	}

	return ""
}

func (m *reqBody) FormValueJson(key string, js interface{}) error {

    dc := json.NewDecoder(strings.NewReader(m.FormValue(key)))
    dc.UseNumber()
    err := dc.Decode(js)
	if err != nil {
		return fmt.Errorf("json unmarshal %s", err.Error())
	} else {
		return nil
	}
}



func (m *reqBody) FormFile(key string) ([]byte, *multipart.FileHeader, error) {

	file, head, err := m.r.FormFile(key)
	if err != nil {
		return nil, nil, fmt.Errorf("get form file err:%s", err)
	}

    data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, fmt.Errorf("get form file data err:%s", err)
	}

	return data, head, nil


}


// ============================
// 没有body类的请求
type HttpRequest struct {
	r *http.Request

	query *reqArgs
	params *reqArgs

	body *reqBody
}

func (m *HttpRequest) Query() *reqArgs {
	return m.query
}

func (m *HttpRequest) Params() *reqArgs {
	return m.params
}

func (m *HttpRequest) Body() *reqBody {
	return m.body
}


func (m *HttpRequest) URL() *url.URL {
	return m.r.URL
}

func (m *HttpRequest) Method() string {
	return m.r.Method
}

func (m *HttpRequest) RemoteAddr() string {
	return m.r.RemoteAddr
}


func (m *HttpRequest) Header() http.Header {
	return m.r.Header
}


func (m *HttpRequest) Request() *http.Request {
	return m.r
}



func NewHttpRequest(r *http.Request, ps httprouter.Params) (*HttpRequest, error) {
	return &HttpRequest {
		r: r,
		query: NewreqArgs(&reqQuery{r: r,}),
		params: NewreqArgs(&reqParams{ps}),
		body: &reqBody{r: r,},
	}, nil
}



func NewHttpRequestJsonBody(r *http.Request, ps httprouter.Params, js interface{}) (*HttpRequest, error) {
	hrb, err := NewHttpRequest(r, ps)
	if err != nil {
		return hrb, err
	}

	err = hrb.Body().Json(js)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal %s", err.Error())
	}


	return hrb, nil

}


type HandleRequest interface {
	Handle(*HttpRequest) HttpResponse
}

type FactoryHandleRequest func() HandleRequest


func HttpRequestWrapper(fac FactoryHandleRequest) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	fun := "HttpRequestWrapper -->"

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req, err := NewHttpRequest(r, ps)
		if err != nil {
			slog.Warnf("%s new request err:%s", fun, err)
			http.Error(w, "new request err:"+err.Error(), 400)
			return
		}

		resp := fac().Handle(req)

		resp.WriteResponse(w)
	}

}

func HttpRequestJsonBodyWrapper(fac FactoryHandleRequest) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	fun := "HttpRequestJsonBodyWrapper -->"

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		newme := fac()
		req, err := NewHttpRequestJsonBody(r, ps, newme)
		if err != nil {
			slog.Warnf("%s body json err:%s", fun, err)
			http.Error(w, "json unmarshal err:"+err.Error(), 400)
			return
		}

		resp := newme.Handle(req)

		resp.WriteResponse(w)
	}

}


// 测试get 获取body ok
// 测试mutlibody 直接获取body,ok
// 测试 application/x-www-form-urlencoded
// 测试 multipart/form-data
