package smetric

import (
	"github.com/julienschmidt/httprouter"
	"github.com/shawnfeng/sutil/slog"
	"net/http"
	"net/http/pprof"
	"runtime"
)

var (
	default_metric_location = "/metrics"
)

type Metricsprocessor struct {
	*Metrics
}

func NewMetricsprocessor() *Metricsprocessor {
	return &Metricsprocessor{DefaultMetrics}
}
func (p *Metricsprocessor) Init() error {
	slog.Infof("init metric instance")
	return nil
}

func (p *Metricsprocessor) Driver() (string, interface{}) {

	handlerFor := p.Metrics.Exportor()
	router := httprouter.New()
	router.Handler("GET", default_metric_location, handlerFor)
	// set only when there's no existing setting
	if runtime.SetMutexProfileFraction(-1) == 0 {
		// 1 out of 5 mutex events are reported, on average
		runtime.SetMutexProfileFraction(5)
	}
	router.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
	router.HandlerFunc("POST", "/debug/pprof/", pprof.Index)
	router.HandlerFunc("GET", "/debug/pprof/allocs", pprof.Index)
	router.HandlerFunc("POST", "/debug/pprof/allocs", pprof.Index)
	router.HandlerFunc("GET", "/debug/pprof/block", pprof.Index)
	router.HandlerFunc("POST", "/debug/pprof/block", pprof.Index)
	router.HandlerFunc("GET", "/debug/pprof/goroutine", pprof.Index)
	router.HandlerFunc("POST", "/debug/pprof/goroutine", pprof.Index)
	router.HandlerFunc("GET", "/debug/pprof/heap", pprof.Index)
	router.HandlerFunc("POST", "/debug/pprof/heap", pprof.Index)
	router.HandlerFunc("GET", "/debug/pprof/mutex", pprof.Index)
	router.HandlerFunc("POST", "/debug/pprof/mutex", pprof.Index)
	router.HandlerFunc("GET", "/debug/pprof/threadcreate", pprof.Index)
	router.HandlerFunc("POST", "/debug/pprof/threadcreate", pprof.Index)
	router.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	router.HandlerFunc("POST", "/debug/pprof/cmdline", pprof.Cmdline)
	router.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
	router.HandlerFunc("POST", "/debug/pprof/profile", pprof.Profile)
	router.HandlerFunc("GET", "/debug/pprof/trace", pprof.Trace)
	router.HandlerFunc("POST", "/debug/pprof/trace", pprof.Trace)
	router.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
	router.HandlerFunc("POST", "/debug/pprof/symbol", pprof.Symbol)
	//router.Handler("POST", p.location,handlerFor)
	router.HandlerFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>A Prometheus Exporter</title></head>
			<body>
			<h1>A Prometheus Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})
	return p.getListenerAddress(), router
}

func (p *Metricsprocessor) getListenerAddress() string {
	return ""
	//port := 22333
	//for {
	//	s := "127.0.0.1:" + strconv.Itoa(port)
	//	_, err := net.Dial("tcp", s)
	//	fmt.Println("dial======", s, err)
	//	if err != nil {
	//		return s
	//	}
	//	port++
	//}
}
