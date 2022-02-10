package getter

import (
	"context"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"requester/getter/md5"
	"requester/test"
)

func TestGetter_GetResponseHash(t *testing.T) {
	type fields struct {
		timeout time.Duration
		hash    iHash
	}
	type args struct {
		ctx context.Context
		to  string
	}
	timeout := time.Millisecond * 100

	srv := httptest.NewServer(test.FastAndSlowHandlers(timeout * 2))
	defer srv.Close()

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantHash []byte
		wantErr  bool
	}{
		{
			name: "success",
			fields: fields{
				timeout: timeout,
				hash:    md5.NewCalculator(),
			},
			args: args{
				ctx: context.Background(),
				to:  fmt.Sprintf("%s/fast", srv.URL),
			},
			wantHash: []byte{0x0c, 0xc1, 0x75, 0xb9, 0xc0, 0xf1, 0xb6, 0xa8,
				0x31, 0xc3, 0x99, 0xe2, 0x69, 0x77, 0x26, 0x61},
			wantErr: false,
		},
		{
			name: "failed: wrong URL",
			fields: fields{
				timeout: timeout,
				hash:    md5.NewCalculator(),
			},
			args: args{
				ctx: context.Background(),
				to:  "wrong url",
			},
			wantHash: nil,
			wantErr:  true,
		},
		{
			name: "failed: slow server context cancelled",
			fields: fields{
				timeout: timeout,
				hash:    md5.NewCalculator(),
			},
			args: args{
				ctx: context.Background(),
				to:  fmt.Sprintf("%s/slow", srv.URL),
			},
			wantHash: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Getter{
				timeout: tt.fields.timeout,
				hash:    tt.fields.hash,
			}
			gotHash, err := g.GetResponseHash(tt.args.ctx, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResponseHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("GetResponseHash() gotHash = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}
