package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/shawnfeng/sutil/snetutil"
	"gitlab.pri.ibanyu.com/middleware/delayqueue/model"
	"gitlab.pri.ibanyu.com/middleware/delayqueue/processor"
	"net/http"
	"strings"
	"time"
)

const (
	defaultToken        = "01E0SSK0DJ9XX4PDFJCD3DN7WX"
	defaultRequestSleep = 300 * time.Millisecond
)

type AckHandler interface {
	Ack(ctx context.Context) error
}

type DelayHandler struct {
	cli   *DelayClient
	jobID string
}

func NewDelayHandler(cli *DelayClient, jobID string) *DelayHandler {
	return &DelayHandler{
		cli:   cli,
		jobID: jobID,
	}
}

func (p *DelayHandler) Ack(ctx context.Context) error {
	return p.cli.Ack(ctx, p.jobID)
}

// 延迟队列客户端
type DelayClient struct {
	endpoint string

	httpCli http.Client

	namespace string
	queue     string

	ttlSeconds uint32
	tries      uint16
	ttrSeconds uint32
}

// 延迟队列任务
type Job struct {
	Namespace string `json:"namespace"`
	Queue     string `json:"queue"`
	Body      []byte `json:"body"` // 任务具体实体
	ID        string `json:"id"`
	TTL       uint32 `json:"ttl"`        // 任务过期时间 单位：s
	Delay     uint32 `json:"delay"`      // 任务延迟时间 单位：s
	ElapsedMS int64  `json:"elapsed_ms"` // 任务从产生到消费时间 单位：ms
}

type writeRes struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg,omitempty"`
	Data struct {
		Ent processor.PublishRes `json:"ent"`
	} `json:"data,omitempty"`
}

type readRes struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg,omitempty"`
	Data struct {
		Ent struct {
			Job *Job `json:"job"`
		} `json:"ent"`
	} `json:"data,omitempty"`
}

type ackRes struct {
	Ret  int      `json:"ret"`
	Msg  string   `json:"msg,omitempty"`
	Data struct{} `json:"data,omitempty"`
}

func NewDelayClient(endpoint, namespace, queue string, ttlSeconds, ttrSeconds uint32, tries uint16) *DelayClient {
	return &DelayClient{
		endpoint:   endpoint,
		namespace:  namespace,
		queue:      queue,
		ttlSeconds: ttlSeconds,
		ttrSeconds: ttrSeconds,
		tries:      tries,
	}
}

// NewDefaultDelayClient 通过topic创建默认客户端
func NewDefaultDelayClient(ctx context.Context, topic string) (*DelayClient, error) {
	Config, err := DefaultConfiger.GetConfig(ctx, topic, MQTypeDelay)
	if err != nil {
		return nil, err
	}
	namespace, queue, err := parseTopic(topic)
	if err != nil {
		return nil, err
	}
	client := NewDelayClient(Config.MQAddr[0], namespace, queue, Config.TTL, Config.TTR, Config.Tries)
	return client, nil
}

// Write 发布任务
func (p *DelayClient) Write(ctx context.Context, value interface{}, ttlSeconds, delaySeconds uint32, tries uint16) (jobID string, err error) {
	fun := "DelayClient.Write --> "
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		p.logConfigToSpan(span)
	}

	msg, err1 := json.Marshal(value)
	if err1 != nil {
		err = fmt.Errorf("%s json marshal, value = %v", err1, value)
		return
	}
	res := new(writeRes)
	req := &processor.Publish{
		Queue:        p.queue,
		Body:         msg,
		TTLSeconds:   ttlSeconds,
		DelaySeconds: delaySeconds,
		Tries:        tries,
	}
	path := fmt.Sprintf("/base/delayqueue/%s/job/publish", p.namespace)
	err = p.httpInvoke(ctx, path, req, res)
	if err != nil {
		return
	}
	if res.Ret == -1 {
		err = fmt.Errorf("%s http invoke, path = %s, err = %s", fun, path, res.Msg)
		return
	}
	jobID = res.Data.Ent.JobID
	return
}

// Read 消费任务
func (p *DelayClient) Read(ctx context.Context, ttrSeconds uint32) (job *Job, err error) {

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		p.logConfigToSpan(span)
	}
	res := new(readRes)
	req := &processor.Consume{
		Queue:      p.queue,
		TTRSeconds: ttrSeconds,
	}
	path := fmt.Sprintf("/base/delayqueue/%s/job/consume", p.namespace)
	for {
		time.Sleep(defaultRequestSleep)
		err = p.httpInvoke(ctx, path, req, res)
		if err != nil {
			break
		}
		if res.Msg == model.ErrNotFound.Error() {
			continue
		}
		if res.Ret == -1 {
			break
		}
		if res.Data.Ent.Job == nil {
			continue
		}
		job = res.Data.Ent.Job
		break
	}
	return
}

// Ack 确认消费
func (p *DelayClient) Ack(ctx context.Context, jobID string) error {
	fun := "DelayClient.Ack -->"
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		p.logConfigToSpan(span)
	}

	res := new(ackRes)
	req := &processor.DeleteJob{
		Queue: p.queue,
		JobID: jobID,
	}
	path := fmt.Sprintf("/base/delayqueue/%s/job/delete", p.namespace)
	err := p.httpInvoke(ctx, path, req, res)
	if err != nil {
		return err
	}
	if res.Ret == -1 {
		return fmt.Errorf("%s http invoke, path = %s, err = %s", fun, path, res.Msg)
	}
	return nil
}

func (p *DelayClient) httpInvoke(ctx context.Context, path string, req interface{}, res interface{}) error {
	url := fmt.Sprintf("%s%s?token=%s", p.endpoint, path, defaultToken)
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resData, code, err := snetutil.HttpReqPost(url, data, time.Minute)
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return fmt.Errorf("http request, url = %s, code = %d", url, code)
	}
	err = json.Unmarshal(resData, &res)
	return err
}

func (p *DelayClient) logConfigToSpan(span opentracing.Span) {
	span.LogFields(
		log.String(spanLogKeyMQType, fmt.Sprint(MQTypeDelay)),
		log.String(spanLogKeyKafkaBrokers, p.endpoint),
	)
}

// topic : group.service.module ==>  namespace: group.service queue: module
func parseTopic(topic string) (namespace, queue string, err error) {
	index := strings.LastIndex(topic, ".")
	if index == -1 {
		err = fmt.Errorf("topic format, topic = %s", topic)
		return
	}
	namespace = topic[:index]
	queue = topic[index+1:]
	return
}

func init()  {
	setHttpDefaultClient()
}

func setHttpDefaultClient() {
	snetutil.DefaultClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 128,
			MaxConnsPerHost:     1024,
			IdleConnTimeout:     600 * time.Second,
		},
		Timeout: 0,
	}
}
