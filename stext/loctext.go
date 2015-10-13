// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package stext


import (
	//"fmt"
	"strconv"

)



func doFmt(f string) int {
	if len(f) == 0 {
		return 0
	}

	idx, err := strconv.Atoi(f)
	if err != nil {
		return -1
	}

	return idx
}

// %3@ %@
func LocationText(format string, params ...string) string {

	var res string

	sz := len(params)

	// 格式化字符的位置
	fb := -1
	fe := -1

	// 截断非格式化字符的起始位置
	cb := 0
	// 记录是第几个没有标记位置(%@)的格式化字符
	pi := 1

	for i, c := range format {
		//fmt.Printf("begin i=%d fb=%d fe=%d cb=%d pi=%d c=%c res=%s\n", i, fb, fe, cb, pi, c, res)

		if c == '%' {
			if fb == -1 {
				fb = i

			} else {
				// 双%%转义到一个%
				if fb == i-1 {
					res += format[cb:i]
					cb = i+1
					fb = -1
				} else {
					fb = i
				}
			}

		} else if c == '@' {
			if fb >= 0 {
				res += format[cb:fb]
				cb = fb

				// 格式化翻译
				fe = i
				idx := doFmt(format[fb+1:fe])
				if idx < 0 {
					res += format[fb:fe+1]
				} else if idx == 0 {
					if pi <= sz {
						res += params[pi-1]
						pi++
					}
				} else {
					if idx <= sz {
						res += params[idx-1]
					}
				}

				cb = fe+1
				fb = -1
				fe = -1
			}
		}


		//fmt.Printf("end i=%d fb=%d fe=%d cb=%d pi=%d c=%c res=%s\n", i, fb, fe, cb, pi, c, res)

	}

	if cb < len(format) {
		res += format[cb:]
	}

	//fmt.Printf("over fb=%d fe=%d cb=%d pi=%d res=%s\n", fb, fe, cb, pi, res)

	return res
}


