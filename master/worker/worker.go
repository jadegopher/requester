package worker

import (
	"context"
	"fmt"

	"requester/getter"
)

type Worker struct {
	sender getter.HashGetter
}

func NewWorker(sender getter.HashGetter) *Worker {
	return &Worker{
		sender: sender,
	}
}

func (w *Worker) Consume(ctx context.Context, urls <-chan string, result chan<- []string) {
	for {
		select {
		case url, isOpen := <-urls:
			if !isOpen {
				return
			}
			hash, err := w.sender.GetResponseHash(ctx, url)
			if err != nil {
				result <- []string{url, fmt.Sprintf("failed to get response hash. Reason: %s", err.Error())}
				continue
			}
			result <- []string{url, fmt.Sprintf("%x", hash)}
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}
