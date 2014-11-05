package stime

import (
	"time"

)

type BackOffCtrl struct {
	// 退避的最大值
	ceil time.Duration
	// 退避的起始值
	step time.Duration
	backtime time.Duration

	reset chan bool
}

func NewBackOffCtrl(step time.Duration, ceil time.Duration) *BackOffCtrl {
	return &BackOffCtrl {
		ceil: ceil,
		step: step,
		backtime: 0,
		reset: make(chan bool),
	}

}

// 执行退避，会发生阻塞
func (m *BackOffCtrl) BackOff() {

	select {
	case <-m.reset:
	case <-time.After(m.backtime):
		if m.backtime <= 0 {
			m.backtime = m.step
		} else {
			m.backtime = m.backtime * 2
		}

		if m.backtime >= m.ceil {
			m.backtime = m.ceil
		}

	}

}

// 终止退避过程，reset退避状态
func (m *BackOffCtrl) Reset() {
	m.backtime = 0
	select {
	case m.reset <-true:
	default:
	}

}


