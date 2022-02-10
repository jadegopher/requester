package worker

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"requester/mocks"
)

func TestWorker_Consume(t *testing.T) {
	type data struct {
		urls []string
	}
	type args struct {
		ctx    context.Context
		urls   chan string
		result chan []string
	}
	type expected struct {
		value []string
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	hash := []byte{0x67}
	url := "url"
	err := errors.New("some error")
	tests := []struct {
		name     string
		args     args
		data     data
		expected expected
		mocks    func(hashGetter *mocks.MockHashGetter)
	}{
		{
			name: "success",
			args: args{
				ctx:    ctx,
				urls:   make(chan string, 1),
				result: make(chan []string, 1),
			},
			data: data{
				urls: []string{url},
			},
			expected: expected{
				value: []string{url, "67"},
			},
			mocks: func(hashGetter *mocks.MockHashGetter) {
				gomock.InOrder(
					hashGetter.EXPECT().GetResponseHash(ctx, url).Return(hash, nil),
				)
			},
		},
		{
			name: "get response hash error",
			args: args{
				ctx:    ctx,
				urls:   make(chan string, 1),
				result: make(chan []string, 1),
			},
			data: data{
				urls: []string{url},
			},
			expected: expected{
				value: []string{url, "failed to get response hash. Reason: some error"},
			},
			mocks: func(hashGetter *mocks.MockHashGetter) {
				gomock.InOrder(
					hashGetter.EXPECT().GetResponseHash(ctx, url).Return(nil, err),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockHashGetter := mocks.NewMockHashGetter(ctrl)
			defer ctrl.Finish()
			w := &Worker{
				sender: mockHashGetter,
			}

			tt.mocks(mockHashGetter)

			go func() {
				w.Consume(tt.args.ctx, tt.args.urls, tt.args.result)
			}()

			for _, u := range tt.data.urls {
				tt.args.urls <- u
			}
			close(tt.args.urls)

			result := <-tt.args.result
			if !reflect.DeepEqual(result, tt.expected.value) {
				t.Errorf("Got = %v, want %v", result, tt.expected.value)
			}
		})
	}
}
