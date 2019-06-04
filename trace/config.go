package trace

import (
	"fmt"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"os"
)

const (
	TRACER_TYPE_JAEGER = "jaeger"
)

const (
	CONFIG_TYPE_SIMPLE = iota
	CONFIG_TYPE_ETCD
)

const (
	ENV_JAEGER_AGENT_HOST = "JAEGER_AGENT_HOST"
	ENV_JAEGER_AGENT_PORT = "JAEGER_AGENT_PORT"
)

const (
	JAEGER_DEBUG_HEADER         = "trace-debug-id"
	JAEGER_BAGGAGE_HEADER       = "trace-baggage"
	TRACE_CONTEXT_HEADER_NAME   = "banyu-trace-id"
	TRACE_BAGGAGE_HEADER_PREFIX = "banyuctx-"
)

var defaultConfigurator = NewSimpleConfigurator()

type Config struct {
	TracerType   string
	ServiceName  string
	TracerConfig interface{}
}

type Configurator interface {
	GetConfig(serviceName string) *Config
}

func NewConfigurator(configType int) (Configurator, error) {

	switch configType {
	case CONFIG_TYPE_SIMPLE:
		return NewSimpleConfigurator(), nil
	case CONFIG_TYPE_ETCD:
		return NewEtcdConfigurator(), nil
	default:
		return nil, fmt.Errorf("configType %d error", configType)
	}
}

type SimpleConfig struct {
}

func NewSimpleConfigurator() Configurator {
	return &SimpleConfig{}
}

func (m *SimpleConfig) GetConfig(serviceName string) *Config {
	// 测试环境
	var agentHost = "127.0.0.1"
	var agentPort = "6831"

	if h := os.Getenv(ENV_JAEGER_AGENT_HOST); h != "" {
		agentHost = h
	}

	if p := os.Getenv(ENV_JAEGER_AGENT_PORT); p != "" {
		agentPort = p
	}

	return &Config{
		TracerType:  TRACER_TYPE_JAEGER,
		ServiceName: serviceName,
		TracerConfig: config.Configuration{
			ServiceName: serviceName,
			Disabled:    false,
			RPCMetrics:  false,
			Sampler: &config.SamplerConfig{
				Type:  "const",
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LocalAgentHostPort: fmt.Sprintf("%s:%s", agentHost, agentPort),
			},
			Headers: &jaeger.HeadersConfig{
				JaegerDebugHeader:        JAEGER_DEBUG_HEADER,
				JaegerBaggageHeader:      JAEGER_BAGGAGE_HEADER,
				TraceContextHeaderName:   TRACE_CONTEXT_HEADER_NAME,
				TraceBaggageHeaderPrefix: TRACE_BAGGAGE_HEADER_PREFIX,
			},
		},
	}
}

type EtcdConfig struct {
	etcdAddrs []string
}

func NewEtcdConfigurator() Configurator {
	// TODO:

	return &EtcdConfig{
		etcdAddrs: []string{},
	}
}

func (m *EtcdConfig) GetConfig(serviceName string) *Config {
	return &Config{
		TracerType:   TRACER_TYPE_JAEGER,
		ServiceName:  serviceName,
		TracerConfig: nil,
	}
}
