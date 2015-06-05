// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package sutil

import (
	"hash/fnv"
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"encoding/json"
	"unicode/utf8"
	"crypto/md5"
	//"crypto/sha1"

	"sync"
	"code.google.com/p/go-uuid/uuid"
)


func HashV(addrs []string, key string) string {
	if len(addrs) == 0 {
		return ""
	}
    h := fnv.New32a()
    h.Write([]byte(key))
    hv := h.Sum32()

	return addrs[hv % uint32(len(addrs))]

}


func IsJSON(s []byte) bool {
    //var js map[string]interface{}
    var js interface{}
    return json.Unmarshal(s, &js) == nil
}

func GetUtf8Chars(s string, num int) string {
	b := []byte(s)
	rv := ""
	for i := 0; len(b)>0 && i < num; i++ {
		_, size := utf8.DecodeRune(b)
		rv += string(b[:size])
		b = b[size:]
	}

	return rv
}

var uuidMu sync.Mutex
func GetUUID() string {
	uuidMu.Lock()
	defer uuidMu.Unlock()

	uuidgen := uuid.NewUUID()
	return uuidgen.String()
}


func GetUniqueMd5() string {
	u := GetUUID()
	h := md5.Sum([]byte(u))
	return fmt.Sprintf("%x", h)
}


// 文件输出，目录不存在自动创建
func WriteFile(path string, data []byte, perm os.FileMode) error {

	idx := strings.LastIndex(path, "/")
	if idx != -1 {
		logdir := path[:idx]
		err := os.MkdirAll(logdir, 0777)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(path, data, perm)
}


type VersionCmp struct {
	ver string
}


func NewVersionCmp(ver string) *VersionCmp {
	v := &VersionCmp{
	}

	v.ver = v.fmtver(ver)
	return v
}


func (m *VersionCmp) fmtver(ver string) string {
	pvs := strings.Split(ver, ".")

	rv := ""
	for _, pv := range(pvs) {
		rv += fmt.Sprintf("%020s", pv)
	}


	return rv

}

func (m *VersionCmp) Min() string {
	return m.fmtver("0")
}

func (m *VersionCmp) Max() string {
	return m.fmtver("99999999999999999999")
}

func (m *VersionCmp) Lt(ver string) bool {
	return m.ver < m.fmtver(ver)
}

func (m *VersionCmp) Lte(ver string) bool {
	return m.ver <= m.fmtver(ver)
}

func (m *VersionCmp) Gt(ver string) bool {
	return m.ver > m.fmtver(ver)
}

func (m *VersionCmp) Gte(ver string) bool {
	return m.ver >= m.fmtver(ver)
}

func (m *VersionCmp) Eq(ver string) bool {
	return m.ver == m.fmtver(ver)
}

func (m *VersionCmp) Ne(ver string) bool {
	return m.ver != m.fmtver(ver)
}

