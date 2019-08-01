package smetric

import (
	"github.com/julienschmidt/httprouter"
	"github.com/shawnfeng/sutil/slog"
	"net/http"
	"net/http/pprof"
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
	router.Handler("GET", "/debug/pprof/", &routeHandlerAdapter{pprof.Index})
	router.Handler("POST", "/debug/pprof/", &routeHandlerAdapter{pprof.Index})
	router.Handler("GET", "/debug/pprof/allocs", &routeHandlerAdapter{pprof.Index})
	router.Handler("POST", "/debug/pprof/allocs", &routeHandlerAdapter{pprof.Index})
	router.Handler("GET", "/debug/pprof/block", &routeHandlerAdapter{pprof.Index})
	router.Handler("POST", "/debug/pprof/block", &routeHandlerAdapter{pprof.Index})
	router.Handler("GET", "/debug/pprof/goroutine", &routeHandlerAdapter{pprof.Index})
	router.Handler("POST", "/debug/pprof/goroutine", &routeHandlerAdapter{pprof.Index})
	router.Handler("GET", "/debug/pprof/heap", &routeHandlerAdapter{pprof.Index})
	router.Handler("POST", "/debug/pprof/heap", &routeHandlerAdapter{pprof.Index})
	router.Handler("GET", "/debug/pprof/mutex", &routeHandlerAdapter{pprof.Index})
	router.Handler("POST", "/debug/pprof/mutex", &routeHandlerAdapter{pprof.Index})
	router.Handler("GET", "/debug/pprof/threadcreate", &routeHandlerAdapter{pprof.Index})
	router.Handler("POST", "/debug/pprof/threadcreate", &routeHandlerAdapter{pprof.Index})
	router.Handler("GET", "/debug/pprof/cmdline", &routeHandlerAdapter{pprof.Cmdline})
	router.Handler("POST", "/debug/pprof/cmdline", &routeHandlerAdapter{pprof.Cmdline})
	router.Handler("GET", "/debug/pprof/profile", &routeHandlerAdapter{pprof.Profile})
	router.Handler("POST", "/debug/pprof/profile", &routeHandlerAdapter{pprof.Profile})
	router.Handler("GET", "/debug/pprof/trace", &routeHandlerAdapter{pprof.Trace})
	router.Handler("POST", "/debug/pprof/trace", &routeHandlerAdapter{pprof.Trace})
	router.Handler("GET", "/debug/pprof/symbol", &routeHandlerAdapter{pprof.Symbol})
	router.Handler("POST", "/debug/pprof/symbol", &routeHandlerAdapter{pprof.Symbol})
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

type routeHandlerAdapter struct {
	ProfFunc func(w http.ResponseWriter, r *http.Request)
}

func (m *routeHandlerAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.ProfFunc(w, r)
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
