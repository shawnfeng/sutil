package statlog

import (
	"context"
	"sync/atomic"
	"time"

	"gitlab.pri.ibanyu.com/middleware/seaweed/xlog"
)

var (
	counter       int64
	fromTimeStamp int64
)

// Init init statlog in xlog
func Init(logDir, fileName string, service string) {
	xlog.InitStatLog(logDir, fileName)
	xlog.SetStatLogService(service)
}

// Sync sync of stat log
func Sync() {
	xlog.StatLogSync()
}

// LogKV 打印info日志，进行统计
func LogKV(ctx context.Context, name string, keysAndValues ...interface{}) {
	xlog.StatInfow(ctx, name, keysAndValues...)
	atomic.AddInt64(&counter, 1)
}

func init() {
	atomic.StoreInt64(&fromTimeStamp, time.Now().Unix())
}

// LogStat 统计日志打印信息
func LogStat() (map[string]int64, []string) {
	st := map[string]int64{
		"TOTAL": atomic.SwapInt64(&counter, 0),
		"STAMP": atomic.SwapInt64(&fromTimeStamp, time.Now().Unix()),
	}

	return st, []string{}
}
