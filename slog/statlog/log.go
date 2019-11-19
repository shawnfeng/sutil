package statlog

import (
	"context"
	"github.com/shawnfeng/lumberjack.v2"
	"github.com/shawnfeng/sutil/scontext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var (
	counter       int64
	fromTimeStamp int64
	lg            *zap.SugaredLogger
	serviceName   string
)

var emptyHeadKV = map[string]interface{}{
	scontext.ContextKeyHeadUid:     0,
	scontext.ContextKeyHeadSource:  0,
	scontext.ContextKeyHeadIp:      "",
	scontext.ContextKeyHeadRegion:  "",
	scontext.ContextKeyHeadDt:      0,
	scontext.ContextKeyHeadUnionId: "",
}

func Sync() error {
	return lg.Sync()
}

func Init(logDir, logPref string, service string) {
	InitV2(logDir, logPref, service, 10240000, 0, 0)
}

func InitV2(logDir, logPref string, service string, maxSize int, maxAge, maxBackups int) {
	serviceName = service

	logFile := ""

	if logDir != "" && logPref != "" {
		logFile = strings.Join([]string{logDir, logPref}, "/")
	}

	var out io.Writer = os.Stdout
	if logFile != "" {
		logger := lumberjack.NewLogger(logFile, maxSize, maxAge, maxBackups, true, false)

		go func() {
			for {
				secsToRotate := 3600 - time.Now().Unix()%3600
				select {
				case <-time.After(time.Second * time.Duration(secsToRotate)):
					_ = logger.Rotate()
				}
			}
		}()

		out = logger
	}

	w := zapcore.AddSync(out)
	enc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:    "",
		LevelKey:      "",
		TimeKey:       "",
		NameKey:       "",
		CallerKey:     "",
		StacktraceKey: "",
	})
	core := zapcore.NewCore(enc, w, zap.InfoLevel)
	lg = zap.New(core).Sugar()
}

func LogKV(ctx context.Context, name string, keysAndValues ...interface{}) {
	headKV := emptyHeadKV
	if chd, ok := ctx.Value(scontext.ContextKeyHead).(scontext.ContextHeader); ok {
		headKV = chd.ToKV()
	}

	body := make(map[string]interface{})
	for i := 0; i < len(keysAndValues); i += 2 {
		// ignore non-equal keys and values
		if i == len(keysAndValues)-1 {
			break
		}

		k, v := keysAndValues[i], keysAndValues[i+1]
		// fail-fast if key is not of string type
		if ks, ok := k.(string); !ok {
			break
		} else {
			body[ks] = v
		}
	}

	kvs := append([]interface{}{},
		"head", headKV,
		"ts", time.Now().Unix(),
		"service", serviceName,
		"name", name,
		"body", body)

	lg.Infow("", kvs...)
	atomic.AddInt64(&counter, 1)
}

func init() {
	Init("", "", "")
	atomic.StoreInt64(&fromTimeStamp, time.Now().Unix())
}

func LogStat() (map[string]int64, []string) {
	st := map[string]int64{
		"TOTAL": atomic.SwapInt64(&counter, 0),
		"STAMP": atomic.SwapInt64(&fromTimeStamp, time.Now().Unix()),
	}

	return st, []string{}
}
