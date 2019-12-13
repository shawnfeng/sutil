// Copyright 2014 The mqrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mq

import (
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/shawnfeng/sutil/scontext"
)

type Payload struct {
	Carrier opentracing.TextMapCarrier `json:"c"`
	Value   string                     `json:"v"`
	Head    interface{}                `json:"h"`
	Control interface{}                `json:"t"`
}

func generatePayload(ctx context.Context, value interface{}) (*Payload, error) {
	carrier := opentracing.TextMapCarrier(make(map[string]string))
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.TextMap,
			carrier)
	}

	msg, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	head := ctx.Value(scontext.ContextKeyHead)
	control := ctx.Value(scontext.ContextKeyControl)

	return &Payload{
		Carrier: carrier,
		Value:   string(msg),
		Head:    head,
		Control: control,
	}, nil
}

func generateMsgsPayload(ctx context.Context, msgs ...Message) ([]Message, error) {
	carrier := opentracing.TextMapCarrier(make(map[string]string))
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.TextMap,
			carrier)
	}
	head := ctx.Value(scontext.ContextKeyHead)
	control := ctx.Value(scontext.ContextKeyControl)

	var nmsgs []Message
	for _, msg := range msgs {
		body, err := json.Marshal(msg.Value)
		if err != nil {
			return nil, err
		}
		nmsgs = append(nmsgs, Message{
			Key: msg.Key,
			Value: &Payload{
				Carrier: carrier,
				Value:   string(body),
				Head:    head,
				Control: control,
			},
		})
	}

	return nmsgs, nil
}

func parsePayload(payload *Payload, opName string, value interface{}) (context.Context, error) {
	tracer := opentracing.GlobalTracer()
	spanCtx, err := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(payload.Carrier))
	var span opentracing.Span
	if err == nil {
		span = tracer.StartSpan(opName, ext.RPCServerOption(spanCtx))
	} else {
		span = tracer.StartSpan(opName)
	}
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = context.WithValue(ctx, scontext.ContextKeyHead, payload.Head)
	ctx = context.WithValue(ctx, scontext.ContextKeyControl, payload.Control)

	err = json.Unmarshal([]byte(payload.Value), value)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
