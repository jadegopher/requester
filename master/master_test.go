package master

import (
	"context"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"requester/getter"
	"requester/getter/md5"
	"requester/test"
)

func TestMaster_ProcessTasks(t *testing.T) {
	type fields struct {
		goroutinesCount int
	}
	type args struct {
		ctx  func() (context.Context, context.CancelFunc)
		urls []string
	}
	timeout := time.Millisecond * 100

	srv := httptest.NewServer(test.FastAndSlowHandlers(timeout * 2))
	defer srv.Close()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "two fast without errors 2 goroutines",
			fields: fields{
				goroutinesCount: 2,
			},
			args: args{
				ctx: func() (context.Context, context.CancelFunc) {
					return context.WithTimeout(context.Background(), timeout*4)
				},
				urls: []string{
					fmt.Sprintf("%s/fast?hash=a", srv.URL),
					fmt.Sprintf("%s/fast?hash=b", srv.URL),
				},
			},
			want: map[string]string{
				fmt.Sprintf("%s/fast?hash=a", srv.URL): "0cc175b9c0f1b6a831c399e269772661",
				fmt.Sprintf("%s/fast?hash=b", srv.URL): "0cc175b9c0f1b6a831c399e269772661",
			},
		},
		{
			name: "two fast without errors 1 goroutine",
			fields: fields{
				goroutinesCount: 1,
			},
			args: args{
				ctx: func() (context.Context, context.CancelFunc) {
					return context.WithTimeout(context.Background(), timeout*4)
				},
				urls: []string{
					fmt.Sprintf("%s/fast?hash=a", srv.URL),
					fmt.Sprintf("%s/fast?hash=b", srv.URL),
				},
			},
			want: map[string]string{
				fmt.Sprintf("%s/fast?hash=a", srv.URL): "0cc175b9c0f1b6a831c399e269772661",
				fmt.Sprintf("%s/fast?hash=b", srv.URL): "0cc175b9c0f1b6a831c399e269772661",
			},
		},
		{
			name: "one fast and one slow with one error",
			fields: fields{
				goroutinesCount: 2,
			},
			args: args{
				ctx: func() (context.Context, context.CancelFunc) {
					return context.WithTimeout(context.Background(), timeout*4)
				},
				urls: []string{
					fmt.Sprintf("%s/fast?hash=a", srv.URL),
					fmt.Sprintf("%s/slow?hash=b", srv.URL),
				},
			},
			want: map[string]string{
				fmt.Sprintf("%s/fast?hash=a", srv.URL): "0cc175b9c0f1b6a831c399e269772661",
				fmt.Sprintf("%s/slow?hash=b", srv.URL): fmt.Sprintf(`failed to get response hash. Reason: Get "%s/slow?hash=b": context deadline exceeded`, srv.URL),
			},
		},
		{
			name: "context cancelled",
			fields: fields{
				goroutinesCount: 2,
			},
			args: args{
				ctx: func() (context.Context, context.CancelFunc) {
					return context.WithTimeout(context.Background(), timeout/2)
				},
				urls: []string{
					fmt.Sprintf("%s/fast?hash=a", srv.URL),
					fmt.Sprintf("%s/slow?hash=b", srv.URL),
				},
			},
			want: map[string]string{
				fmt.Sprintf("%s/fast?hash=a", srv.URL): "0cc175b9c0f1b6a831c399e269772661",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Master{
				goroutinesCount: tt.fields.goroutinesCount,
				hashGetter:      getter.NewGetter(timeout, md5.NewCalculator()),
			}
			ctx, cancelFunc := tt.args.ctx()
			if got := m.ProcessTasks(ctx, tt.args.urls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessTasks() = %v, want %v", got, tt.want)
			}
			cancelFunc()
		})
	}
}
