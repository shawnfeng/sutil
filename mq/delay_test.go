package mq

import (
	"context"
	"fmt"
	"testing"
)

var delayCli *DelayClient

func init() {
	delayCli = NewDelayClient("http://0.0.0.0:7777", "base.test", "test", 100, 50, 1)
}

func Test_parseTopic(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name          string
		args          args
		wantNamespace string
		wantQueue     string
		wantErr       bool
	}{
		{
			name:          "success",
			args:          args{topic: "base.changeboard.event"},
			wantNamespace: "base.changeboard",
			wantQueue:     "event",
			wantErr:       false,
		},
		{
			name:          "success",
			args:          args{topic: "base.changeboard"},
			wantNamespace: "base",
			wantQueue:     "changeboard",
			wantErr:       false,
		},
		{
			name:          "format error",
			args:          args{topic: "base"},
			wantNamespace: "",
			wantQueue:     "",
			wantErr:       true,
		},
		{
			name:   "group topic",
			args:args{topic: "base.changeboard.event_t1"},
			wantNamespace: "base.changeboard",
			wantQueue:     "event_t1",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNamespace, gotQueue, err := parseTopic(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTopic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNamespace != tt.wantNamespace {
				t.Errorf("parseTopic() gotNamespace = %v, want %v", gotNamespace, tt.wantNamespace)
			}
			if gotQueue != tt.wantQueue {
				t.Errorf("parseTopic() gotQueue = %v, want %v", gotQueue, tt.wantQueue)
			}
		})
	}
}

func TestDelayClient_Write(t *testing.T) {
	ctx := context.Background()
	jobID, err := delayCli.Write(ctx, "test", delayCli.ttlSeconds, 3, 1)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(jobID)

}

// 01E0WTH8TCWT6Q2J2HQ0500000  01E0WTPYNE33SHMHSYW0500000 01E0WTT12VCRFWXXWNS4500000  01E0WVCSVE7TZZHFGRWC000000
func TestDelayClient_Read(t *testing.T) {
	ctx := context.Background()
	job, err := delayCli.Read(ctx, 5)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(job.Body))
	fmt.Println(job.ElapsedMS)
	fmt.Println("success")
	//return

}
