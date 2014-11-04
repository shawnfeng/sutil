package stime

import (
	"time"

)

type BackOffCtrl struct {
	// 退避的最大值
	ceil int32
	// 退避累计
	cn uint32

	reset chan bool
}

func NewBackOffCtrl(ceil int32) *BackOffCtrl {
	return &BackOffCtrl {
		ceil: ceil,
		cn: 0,
		reset: make(chan bool),
	}

}

// 执行退避，会发生阻塞
func (m *BackOffCtrl) BackOff() {

	backtime := int32(1 << m.cn)
	if 1 << m.cn > m.ceil {
		backtime = m.ceil
	} else {
		m.cn++
	}

	select {
	case <-m.reset:
	case <-time.After(time.Second * time.Duration(backtime)):
	}


}

// 终止退避过程，reset退避状态
func (m *BackOffCtrl) Reset() {
	m.cn = 0
	select {
	case m.reset <-true:
	default:
	}

}


