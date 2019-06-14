package trace

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)

// singleton
var bt *backTracer

func InitDefaultTracer(serviceName string) error {
	return InitTracer(TRACER_TYPE_JAEGER, serviceName)
}

func InitTracer(tracerType string, serviceName string) error {
	// only init once
	if bt != nil {
		return nil
	}
	switch tracerType {
	case TRACER_TYPE_JAEGER:
		return initJaeger(serviceName)
	default:
		return fmt.Errorf("unsupported tracer type:%s", tracerType)
	}
}

func CloseTracer() error {
	if bt != nil {
		return bt.closer.Close()
	}
	return nil
}

type backTracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func initJaeger(serviceName string) error {
	fun := "initJaeger-->"
	defaultConfig := defaultConfigurator.GetConfig(serviceName)
	if cfg, ok := defaultConfig.TracerConfig.(config.Configuration); !ok {
		return fmt.Errorf("wrong tracer config:%v", cfg)
	} else {
		tracer, closer, err := cfg.NewTracer()
		if err != nil {
			return err
		}
		bt = &backTracer{tracer, closer}

		opentracing.SetGlobalTracer(bt.tracer)
		fmt.Printf("%s succeed for %s\n", fun, serviceName)
		return nil
	}
}
