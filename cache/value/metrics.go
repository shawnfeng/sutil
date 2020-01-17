package value

import (
	"gitlab.pri.ibanyu.com/middleware/seaweed/xstat/xmetric/xprometheus"
)

const(
	namespace = "palfish"
	subsystem = "cache_requests"
)

var(
	buckets = []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500}

	_metricRequestDuration = xprometheus.NewHistogram(&xprometheus.HistogramVecOpts{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       "duration_ms",
		Help:       "cache requests duration(ms), include load time.",
		LabelNames: []string{"namespace", "command"},
		Buckets:    buckets,
	})
	_metricReqErr= xprometheus.NewCounter(&xprometheus.CounterVecOpts{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       "err_total",
		Help:       "cache.value error total",
		LabelNames: []string{"namespace", "command"},
	})
	_metricHits= xprometheus.NewCounter(&xprometheus.CounterVecOpts{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       "hits_total",
		Help:       "cache.value hits total",
		LabelNames: []string{"namespace", "command"},
	})
	_metricMiss= xprometheus.NewCounter(&xprometheus.CounterVecOpts{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       "miss_total",
		Help:       "cache.value miss total",
		LabelNames: []string{"namespace", "command"},
	})

)

func statReqDuration(namespace,command string, durationMS int64){
	_metricRequestDuration.With("namespace", namespace, "command", command).Observe(float64(durationMS))
}

func statReqErr(namespace, command string, err error){
	if err != nil {
		_metricReqErr.With("namespace", namespace, "command", command).Inc()
	}
	return
}
