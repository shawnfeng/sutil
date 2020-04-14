package dbrouter

import (
	"gitlab.pri.ibanyu.com/middleware/seaweed/xstat/xmetric/xprometheus"
)

const (
	namespace = "palfish"
	subsystem = "db"
)

var (
	_metricReqErr = xprometheus.NewCounter(&xprometheus.CounterVecOpts{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       "request_err_total",
		Help:       "db request err total",
		LabelNames: []string{xprometheus.LabelSource},
	})
)

func statReqErr(cluster, table string, err error) {
	if err != nil {
		source := cluster + "." + table
		_metricReqErr.With(xprometheus.LabelSource, source).Inc()
	}
	return
}
