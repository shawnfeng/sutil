package smetric

import (
	"strconv"
	"time"
)

const (
	Name_space_palfish          = "palfish"
	Name_server_req_total       = "server_request_total"
	Name_server_duration_second = "server_duration_second"
	Label_instance              = "instance"
	Label_servname              = "servname"
	Label_servid                = "servid"
	Label_api                   = "api"
	Label_type                  = "type"
	Label_source                = "source"
	Label_status                = "status"
	Status_succ                 = 1
	Status_fail                 = 0
)

var metricReqNameKeys = []string{Name_space_palfish, Name_server_req_total}
var metricDurationNameKeys = []string{Name_space_palfish, Name_server_duration_second}

//服务相关数据收集，主要包括：成功率、响应时间
//instance:当前服务实例
//servkey:服务标识如base/sla；processor：进程类型 proc_thrift/proc_http
//duration耗时；source请求源标识，目前都是0；servid：服务id标识；funcName：监控的api标识
//err 是否处理有报错
func CollectServ(instance, servkey string, servid int, processor string, duration time.Duration, source int, funcName string, err interface{}) {
	durlabels := buildSerLabels(instance, servkey, servid, processor, source, funcName)
	DefaultMetrics.AddHistoramSampleCreateIfAbsent(metricDurationNameKeys, duration.Seconds(), durlabels, nil)
	var counterLabels []Label
	if err == nil {
		counterLabels = buildSerReqLabels(instance, servkey, servid, processor, source, funcName, Status_succ)
	} else {
		counterLabels = buildSerReqLabels(instance, servkey, servid, processor, source, funcName, Status_fail)
	}
	DefaultMetrics.IncrCounterCreateIfAbsent(metricReqNameKeys, 1.0, counterLabels)
}
func buildSerLabels(instance, servkey string, servid int, processor string, source int, funcName string) []Label {
	targetServerName := SafePromethuesValue(servkey)
	sid := strconv.Itoa(servid)
	return []Label{
		{Name: Label_instance, Value: instance},
		{Name: Label_servname, Value: targetServerName},
		{Name: Label_servid, Value: sid},
		{Name: Label_api, Value: funcName},
		{Name: Label_source, Value: strconv.Itoa(source)},
		{Name: Label_type, Value: processor},
	}
}
func buildSerReqLabels(instance, servkey string, servid int, processor string, source int, funcName string, status int) []Label {
	labels := buildSerLabels(instance, servkey, servid, processor, source, funcName)
	labels = append(labels, Label{
		Name: Label_status, Value: strconv.Itoa(status),
	})
	return labels
}
