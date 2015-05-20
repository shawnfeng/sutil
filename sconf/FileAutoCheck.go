// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package sconf

import (
	"os"
	"fmt"
	"io/ioutil"
	"crypto/sha1"
)


type FileAutoCheck struct {
	file string
	modtime int64
	filehash string

}

func NewFileAutoCheck(file string) *FileAutoCheck {

	return &FileAutoCheck {
		file: file,

	}
}

func (m *FileAutoCheck) Check() (bool, []byte, error) {
	info, err := os.Stat(m.file)

	if err != nil {
		return false, nil, err
	}

	stamp := info.ModTime().Unix()
	if stamp == m.modtime {
		// 不需要更新
		return false, nil, nil
	}
	m.modtime = stamp
	data, err := ioutil.ReadFile(m.file)
	if err != nil {
		return false, nil, err
	}

	h := sha1.Sum(data)
	hex := fmt.Sprintf("%x", h)
	if hex == m.filehash {
		// hash值一样不需要更改
		return false, nil, nil
	}
	m.filehash = hex

	return true, data, nil
}
