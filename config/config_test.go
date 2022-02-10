package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		goroutinesCount int
		urls            []string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "failed: goroutines count == 0",
			args: args{
				goroutinesCount: 0,
				urls:            []string{"some.url"},
			},
			want:    Config{},
			wantErr: true,
		},
		{
			name: "failed: goroutines count < 0",
			args: args{
				goroutinesCount: -1,
				urls:            []string{"some.url"},
			},
			want:    Config{},
			wantErr: true,
		},
		{
			name: "failed: urls nil",
			args: args{
				goroutinesCount: 12,
				urls:            nil,
			},
			want:    Config{},
			wantErr: true,
		},
		{
			name: "failed: urls empty",
			args: args{
				goroutinesCount: 12,
				urls:            []string{},
			},
			want:    Config{},
			wantErr: true,
		},
		{
			name: "success: urls without prefix",
			args: args{
				goroutinesCount: 12,
				urls:            []string{"some.url"},
			},
			want: Config{
				GoroutinesCount: 12,
				URLs:            []string{"http://some.url"},
			},
			wantErr: false,
		},
		{
			name: "success: urls with prefix",
			args: args{
				goroutinesCount: 12,
				urls:            []string{"http://some.url"},
			},
			want: Config{
				GoroutinesCount: 12,
				URLs:            []string{"http://some.url"},
			},
			wantErr: false,
		},
		{
			name: "success: urls",
			args: args{
				goroutinesCount: 12,
				urls:            []string{"http://some.url", "https://ssome.url", "sssome.url"},
			},
			want: Config{
				GoroutinesCount: 12,
				URLs:            []string{"http://some.url", "https://ssome.url", "http://sssome.url"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.goroutinesCount, tt.args.urls)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
