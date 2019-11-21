package spromths

import (
	"context"
	"testing"
	"time"
)

var client *PromethsInstance
var ctx = context.Background()

func init() {
	client, _ = NewPromethsInstance("http://sla.prometheus.pri.ibanyu.com/")
}
func TestPromeths_Query(t *testing.T) {
	query := `
 histogram_quantile(0.99, sum(rate(base_sla_rtcquality_duration_second_bucket{api=~"9002"}[10m]))by (servname,le,api) ) 
`
	//	query:=`
	// sum(increase(base_sla_rtcquality_duration_second_sum{api=~"9002"}[10m]))/sum(increase(base_sla_rtcquality_duration_second_count{api=~"9002"}[10m]))
	//`
	value, e := client.Query(ctx, query, time.Now())
	t.Log(value.Type())
	t.Log(value, "====", e)
}
func TestPromeths_QueryVector(t *testing.T) {
	query := `
 sum(increase(base_sla_rtcquality_duration_second_sum{api=~"9002"}[10m]))/sum(increase(base_sla_rtcquality_duration_second_count{api=~"9002"}[10m]))
`
	value, e := client.QueryVector(ctx, query, time.Now())
	for i, v := range value {
		t.Log(i, "===", v.Value, v.Metric)
	}
	t.Log(value, "====", e)
}
func TestPromeths_QueryVector2(t *testing.T) {
	query := `
 histogram_quantile(0.99, sum(rate(base_sla_rtcquality_duration_second_bucket{api=~"9002"}[10m]))by (servname,le,api) ) 
`
	value, e := client.QueryVector(ctx, query, time.Now())
	for i, v := range value {
		t.Log(i, "===", v.Value, v.Metric)
	}
	t.Log(value, "====", e)
}
