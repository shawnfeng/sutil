package dbrouter

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shawnfeng/sutil/sconf/center"
	"github.com/shawnfeng/sutil/slog/slog"
)

var (
	configCenter center.ConfigCenter
)

const (
	globalGranularityKey = "global.granularity" // 精度
	globalThresholdKey   = "global.threshold"   // 精度内阈值
	globalBreakerGapKey  = "global.breakergap"  // 触发熔断后的熔断间隔,单位: 秒

	checkTick         = time.Millisecond * 25
	defaultThreshold  = 10
	defaultBreakerGap = 10
)

// TOOD 简单计数法实现熔断操作，后续改为滑动窗口或三方组件的方式
type BreakerManager struct {
	lock     sync.Mutex
	Breakers map[string]*Breaker
}

type Breaker struct {
	Rejected      int32
	RejectedStart int64
	Count         int32
}

var bm *BreakerManager

func statBreaker(cluster, table string, err error) {
	if err != nil && (strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "invalid connection")) {
		key := concat(cluster, "_", table)
		bm.lock.Lock()
		if _, ok := bm.Breakers[key]; !ok {
			breaker := new(Breaker)
			breaker.Run()
			bm.Breakers[key] = breaker
		}
		breaker := bm.Breakers[key]
		bm.lock.Unlock()
		atomic.AddInt32(&breaker.Count, 1)
	}
}

func Entry(cluster, table string) bool {
	key := concat(cluster, "_", table)
	bm.lock.Lock()
	breaker := bm.Breakers[key]
	bm.lock.Unlock()
	if breaker != nil {
		return atomic.LoadInt32(&breaker.Rejected) != 1
	}
	return true
}

func (breaker *Breaker) Run() {
	go func() {
		granularityStr, exist := configCenter.GetStringWithNamespace(context.TODO(), center.DefaultApolloMysqlNamespace, globalGranularityKey)
		if !exist {
			slog.Warnf(context.TODO(), "dbrouter: get granularity from apollo failed, exist: %v", exist)
			granularityStr = "1s"
		}
		granularity, err := time.ParseDuration(granularityStr)
		if err != nil {
			slog.Warnf(context.TODO(), "dbrouter: granularity in apollo is invalid, %s", granularityStr)
			granularity = time.Second * 1
		}
		granularityTickC := time.Tick(granularity)
		checkTickC := time.Tick(checkTick)
		for {
			select {
			case <-granularityTickC:
				atomic.StoreInt32(&breaker.Count, 0)
				// check 1s/checkTick times in 1s
			case <-checkTickC:
				threshold, exist := configCenter.GetIntWithNamespace(context.TODO(), center.DefaultApolloMysqlNamespace, globalThresholdKey)
				if !exist {
					slog.Warnf(context.TODO(), "dbrouter: get threshold from apollo failed, exist: %v", exist)
					threshold = defaultThreshold
				}
				breakerGap, exist := configCenter.GetIntWithNamespace(context.TODO(), center.DefaultApolloMysqlNamespace, globalBreakerGapKey)
				if !exist {
					slog.Warnf(context.TODO(), "dbrouter: get breakGap from apollo failed, exist: %v", exist)
					breakerGap = defaultBreakerGap
				}
				if atomic.LoadInt32(&breaker.Count) > int32(threshold) {
					atomic.StoreInt32(&breaker.Rejected, 1)
					breaker.RejectedStart = time.Now().Unix()
				} else {
					now := time.Now().Unix()
					if now-breaker.RejectedStart > int64(breakerGap) {
						atomic.StoreInt32(&breaker.Rejected, 0)
					}
				}
			}
		}
	}()
}

func initConfig() error {
	var err error
	configCenter, err = center.NewConfigCenter(center.ApolloConfigCenter)
	if err != nil {
		return err
	}
	err = configCenter.Init(context.TODO(), center.DefaultApolloMiddlewareService, []string{center.DefaultApolloMysqlNamespace})
	if err != nil {
		return err
	}
	return nil
}

func concat(strings ...string) string {
	var buffer bytes.Buffer
	for _, s := range strings {
		buffer.WriteString(s)
	}
	return buffer.String()
}

func init() {
	bm = &BreakerManager{Breakers: make(map[string]*Breaker)}
	err := initConfig()
	if err != nil {
		slog.Panicf(context.TODO(), "dbrouter: init apollo config failed, err: %v", err)
	}
}
