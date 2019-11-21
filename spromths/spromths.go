package spromths

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"net"
	"net/http"
	"time"
)

type PromethsInstance struct {
	Api v1.API
}

func NewPromethsInstanceWichApi(api v1.API) *PromethsInstance {
	return &PromethsInstance{
		Api: api,
	}
}

func NewPromethsInstance(addr string) (*PromethsInstance, error) {
	client, err := api.NewClient(api.Config{
		Address: addr,
		RoundTripper: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConnsPerHost: 20,
		},
	})
	if err != nil {
		return nil, err
	}
	promethuesApi := v1.NewAPI(client)
	return NewPromethsInstanceWichApi(promethuesApi), nil
}
func (m *PromethsInstance) Query(ctx context.Context, query string, time time.Time) (model.Value, error) {
	return m.Api.Query(ctx, query, time)
}
func (m *PromethsInstance) QueryRange(ctx context.Context, query string, r v1.Range) (model.Value, error) {
	return m.Api.QueryRange(ctx, query, r)
}

func (m *PromethsInstance) QueryMatrix(ctx context.Context, query string, time time.Time) (model.Matrix, error) {
	value, err := m.Query(ctx, query, time)
	if err != nil {
		return nil, err
	}
	v, ok := value.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("not Matrix")
	}
	return v, nil
}

func (m *PromethsInstance) QueryVector(ctx context.Context, query string, time time.Time) (model.Vector, error) {
	value, err := m.Query(ctx, query, time)
	if err != nil {
		return nil, err
	}
	v, ok := value.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("not Vector")
	}
	return v, nil
}

func (m *PromethsInstance) QueryScalar(ctx context.Context, query string, time time.Time) (*model.Scalar, error) {
	value, err := m.Query(ctx, query, time)
	if err != nil {
		return nil, err
	}
	v, ok := value.(*model.Scalar)
	if !ok {
		return nil, fmt.Errorf("not Scalar")
	}
	return v, nil
}

func (m *PromethsInstance) QueryString(ctx context.Context, query string, time time.Time) (*model.String, error) {
	value, err := m.Query(ctx, query, time)
	if err != nil {
		return nil, err
	}
	v, ok := value.(*model.String)
	if !ok {
		return nil, fmt.Errorf("not String")
	}
	return v, nil
}

func (m *PromethsInstance) QueryRangeMatrix(ctx context.Context, query string, r v1.Range) (model.Matrix, error) {
	value, err := m.QueryRange(ctx, query, r)
	if err != nil {
		return nil, err
	}
	v, ok := value.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("not Matrix")
	}
	return v, nil
}

func (m *PromethsInstance) QueryRangeVector(ctx context.Context, query string, r v1.Range) (model.Vector, error) {
	value, err := m.QueryRange(ctx, query, r)
	if err != nil {
		return nil, err
	}
	v, ok := value.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("not Vector")
	}
	return v, nil
}

func (m *PromethsInstance) QueryRangeScalar(ctx context.Context, query string, r v1.Range) (*model.Scalar, error) {
	value, err := m.QueryRange(ctx, query, r)
	if err != nil {
		return nil, err
	}
	v, ok := value.(*model.Scalar)
	if !ok {
		return nil, fmt.Errorf("not Scalar")
	}
	return v, nil
}

func (m *PromethsInstance) QueryRangeString(ctx context.Context, query string, r v1.Range) (*model.String, error) {
	value, err := m.QueryRange(ctx, query, r)
	if err != nil {
		return nil, err
	}
	v, ok := value.(*model.String)
	if !ok {
		return nil, fmt.Errorf("not String")
	}
	return v, nil
}
