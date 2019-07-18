// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package stext

import (
	"context"
	"testing"
	"unicode/utf8"
	"github.com/shawnfeng/sutil/slog/slog"
)


func TestSymb(t *testing.T) {
	s, err := NewSymbolList("symbol.list")
	if err != nil {
		slog.Errorln(context.TODO(), "load symbol")
	}


	b := []byte("\t!@#$%^&*不同国家的女人难受时如何用同一句话安慰？神回复：欧美 - You need cry，dear（你需要哭出来，宝贝）中国 - 有你\n的快递儿！")
	bb := []byte("\t!@#$%^&*不同国家的女人难受时如何用同一句话安慰？神回复：欧美 - You need cry，dear（你需要哭出来，宝贝）中国 - 有你\n的快递儿！")


	cmp := "\t!@#$%^&*？： -   ，（，） - \n！"
	pick := ""
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if s.Is(r) {
			slog.Infof(context.TODO(), "%c", r)
			pick += string(b[:size])
		}
		b = b[size:]
	}

	slog.Infoln(context.TODO(), "cmp:", cmp)
	slog.Infoln(context.TODO(), "pick:", pick)

	if cmp != pick {
		t.Errorf("not ok")
	}

	rv := BytesToRunesNoSymb(s, bb)

	slog.Infof(context.TODO(), "%s   %s", bb, rv)

	for i := 0; i < len(rv); i++ {
		slog.Infof(context.TODO(), "%c", rv[i])
	}



	rv = BytesToRunes(bb)

	slog.Infof(context.TODO(), "%s   %s", bb, rv)

	for i := 0; i < len(rv); i++ {
		slog.Infof(context.TODO(), "%c", rv[i])
	}

	//s.Print()

}
