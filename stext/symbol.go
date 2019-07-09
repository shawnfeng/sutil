// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package stext

import (
	"fmt"
	"io/ioutil"
	"bytes"

	"unicode/utf8"
)


type SymbolList struct {
	list map[rune] []byte
}

func NewSymbolList(file string) (*SymbolList, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	m := &SymbolList {
		list: make(map[rune] []byte),
	}

	items := bytes.Split(data, []byte("\n"))
	for _, it := range(items) {
		r, size := utf8.DecodeRune(it)
		if size == len(it) {
			_, ok := m.list[r]
			if ok {
				//slog.Warnf(context.TODO(), "same %c %d %s %v", r, size, it, it)
			} else {
				m.list[r] = it
			}
		} else {
			//slog.Warnf(context.TODO(), "illigal %s %d %s %v", r, size, it, it)
		}
	}

	for i := 0; i < 256; i++ {
		if i >= 'A' && i <= 'Z' {
			continue
		}

		if i >= 'a' && i <= 'z' {
			continue
		}

		if i >= '0' && i <= '9' {
			continue
		}

		m.list[rune(i)] = []byte{byte(i)}

	}


	return m, nil
}

func (m *SymbolList)Is(c rune) bool {
	_, ok := m.list[c]
	return ok
}


func (m *SymbolList) Print() {

	s := ""
	for _, v := range(m.list) {
		if len(s) == 0 {
			s = fmt.Sprintf("%s", v)
		} else {
			s = fmt.Sprintf("%s,%s", s, v)
		}
	}

	fmt.Println(s)
}


func BytesToRunesNoSymb(s *SymbolList, b []byte) []rune {
	rv := make([]rune, 0)
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)

		b = b[size:]

		if s.Is(r) {
			continue
		}


		rv = append(rv, r)
		
	}

	return rv
}



func BytesToRunes(b []byte) []rune {
	rv := make([]rune, 0)
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		rv = append(rv, r)
		b = b[size:]
		rv = append(rv, r)

	}

	return rv
}
