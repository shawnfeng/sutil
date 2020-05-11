package redis

import (
	"reflect"
	"testing"
)

func Test_instanceConfFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		wantConf *InstanceConf
		wantErr  bool
	}{
		{
			name: "test no prefix key",
			args: args{
				s: "default-base/test-e-false",
			},
			wantConf:  &InstanceConf{
				Group:     "default",
				Namespace: "base/test",
				Wrapper:   "e",
				NoFixKey:  false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConf, err := instanceConfFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("instanceConfFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConf, tt.wantConf) {
				t.Errorf("instanceConfFromString() gotConf = %v, want %v", gotConf, tt.wantConf)
			}
		})
	}
}
