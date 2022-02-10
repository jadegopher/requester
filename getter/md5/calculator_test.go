package md5

import (
	"reflect"
	"testing"
)

func TestCalculator_Calculate(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "nil data",
			args: args{
				data: nil,
			},
			want: []byte{0xd4, 0x1d, 0x8c, 0xd9, 0x8f, 0x00, 0xb2, 0x04,
				0xe9, 0x80, 0x09, 0x98, 0xec, 0xf8, 0x42, 0x7e},
		},
		{
			name: "empty data",
			args: args{
				data: []byte{},
			},
			want: []byte{0xd4, 0x1d, 0x8c, 0xd9, 0x8f, 0x00, 0xb2, 0x04,
				0xe9, 0x80, 0x09, 0x98, 0xec, 0xf8, 0x42, 0x7e},
		},
		{
			name: "hash check",
			args: args{
				data: []byte{'a'},
			},
			want: []byte{0x0c, 0xc1, 0x75, 0xb9, 0xc0, 0xf1, 0xb6, 0xa8,
				0x31, 0xc3, 0x99, 0xe2, 0x69, 0x77, 0x26, 0x61},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Calculator{}
			if got := c.Calculate(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
