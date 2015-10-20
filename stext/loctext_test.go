// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package stext

import (
	"testing"
	"github.com/shawnfeng/sutil/slog"
)


func ck(t *testing.T, correct, format string, params ...string) {

	res := LocationText(format, params...)
	slog.Infoln(format, params, res)

	if res != correct {
		t.Errorf("correct:%s format:%s params:%s res:%s", correct, format, params, res)
	}
}

func TestLocText(t *testing.T) {
	// 空测试
	ck(t, "", "", "123")
	ck(t, "123", "%@", "123")
	ck(t, "123", "%@%@", "123")
	ck(t, "123%@", "%@%%@", "123")
	ck(t, "123%a@@", "%@%a@@", "123")

	ck(t, "%abc", "%@%abc")
	ck(t, "123%abc", "%@%abc", "123")

	// 放在开头
	ck(t, "123abc", "%@abc", "123")



	// 放在末尾
	ck(t, "世界你好abc123", "世界你好abc%@", "123")


	// %%
	ck(t, "abc%@", "abc%%@", "123")


	// %% 混合
	ck(t, "abc%@123", "abc%%@%@", "123")

	// 多参数
	ck(t, "456abc%@123", "%@abc%%@%@", "456", "123")


	// 多余参数
	ck(t, "456abc%@123", "%@abc%%@%@%@%@%@", "456", "123")


	// 指定位置参数
	ck(t, "123abc%@456123", "%2$@abc%%@%@%@%@%@", "456", "123")


	// 指定位置参数
	ck(t, "123abc%@456123456", "%2$@abc%%@%@%@%@%1$@", "456", "123")


	// 指定位置参数
	ck(t, "123 abc %@ 456 123 456", "%2$@ abc %%@ %@ %@ %1$@", "456", "123")

	// 纯指定位置参数
	ck(t, "123 abc 456", "%2$@ abc %1$@", "456", "123")


	// 越界指定
	ck(t, "123 abc", "%2$@ abc%10$@", "456", "123")

	// 非法指定
	ck(t, "123 abc%-2@ abc%-2$@", "%2$@ abc%-2@ abc%-2$@", "456", "123")
	ck(t, "123 abc%ab@ abc%ab$@", "%2$@ abc%ab@ abc%ab$@", "456", "123")

	// %0@
	ck(t, "123 abc%ab@456123 789", "%2$@ abc%ab@%0$@%@ %@", "456", "123", "789")

	// %$@ $1@
	ck(t, "%1@ 123 abc%ab@456123 789 %$@", "%1@ %2$@ abc%ab@%0$@%@ %@ %$@", "456", "123", "789")

	return

}





func TestLocText2(t *testing.T) {
	// 空测试
	ck(t, "", "", "123")
	ck(t, "123", "%s", "123")
	ck(t, "123", "%s%s", "123")
	ck(t, "123%s", "%s%%s", "123")
	ck(t, "123%ass", "%s%ass", "123")

	ck(t, "%abc", "%s%abc")
	ck(t, "123%abc", "%s%abc", "123")

	// 放在开头
	ck(t, "123abc", "%sabc", "123")



	// 放在末尾
	ck(t, "世界你好abc123", "世界你好abc%s", "123")


	// %%
	ck(t, "abc%s", "abc%%s", "123")


	// %% 混合
	ck(t, "abc%s123", "abc%%s%s", "123")

	// 多参数
	ck(t, "456abc%s123", "%sabc%%s%s", "456", "123")


	// 多余参数
	ck(t, "456abc%s123", "%sabc%%s%s%s%s%s", "456", "123")


	// 指定位置参数
	ck(t, "123abc%s456123", "%2$sabc%%s%s%s%s%s", "456", "123")


	// 指定位置参数
	ck(t, "123abc%s456123456", "%2$sabc%%s%s%s%s%1$s", "456", "123")


	// 指定位置参数
	ck(t, "123 abc %s 456 123 456", "%2$s abc %%s %s %s %1$s", "456", "123")

	// 纯指定位置参数
	ck(t, "123 abc 456", "%2$s abc %1$s", "456", "123")


	// 越界指定
	ck(t, "123 abc", "%2$s abc%10$s", "456", "123")

	// 非法指定
	ck(t, "123 abc%-2s abc%-2$s", "%2$s abc%-2s abc%-2$s", "456", "123")
	ck(t, "123 abc%abs abc%ab$s", "%2$s abc%abs abc%ab$s", "456", "123")

	// %0s
	ck(t, "123 abc%abs456123 789", "%2$s abc%abs%0$s%s %s", "456", "123", "789")

	// %$s $1s
	ck(t, "%1s 123 abc%abs456123 789 %$s", "%1s %2$s abc%abs%0$s%s %s %$s", "456", "123", "789")

	return

}






