package stime

import (
	"time"
	"sync/atomic"
)

type BackOffCtrl struct {
	// 退避的最大值
	ceil int64
	// 退避的起始值
	step int64
	backtime int64

	reset chan bool
}

func NewBackOffCtrl(step time.Duration, ceil time.Duration) *BackOffCtrl {
	return &BackOffCtrl {
		ceil: ceil.Nanoseconds(),
		step: step.Nanoseconds(),
		backtime: 0,
		reset: make(chan bool),
	}

}

func (m *BackOffCtrl) SetCtrl(step time.Duration, ceil time.Duration) {
	atomic.StoreInt64(&m.step, step.Nanoseconds())
	atomic.StoreInt64(&m.ceil, ceil.Nanoseconds())
	m.Reset()
}

// 执行退避，会发生阻塞
func (m *BackOffCtrl) BackOff() {

	select {
	case <-m.reset:
	case <-time.After(time.Duration(atomic.LoadInt64(&m.backtime))):
		if atomic.LoadInt64(&m.backtime) <= 0 {
			atomic.StoreInt64(&m.backtime, atomic.LoadInt64(&m.step))
		} else {
			//m.backtime = m.backtime * 2
			atomic.StoreInt64(&m.backtime, atomic.LoadInt64(&m.backtime)*2)
		}

		if atomic.LoadInt64(&m.backtime) >= atomic.LoadInt64(&m.ceil) {
			atomic.StoreInt64(&m.backtime, atomic.LoadInt64(&m.ceil))
		}

	}

}

// 终止退避过程，reset退避状态
func (m *BackOffCtrl) Reset() {
	atomic.StoreInt64(&m.backtime, 0)
	select {
	case m.reset <-true:
	default:
	}

}


