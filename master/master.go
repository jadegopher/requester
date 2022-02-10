package master

import (
	"context"

	"requester/getter"
	"requester/master/worker"
)

type Master struct {
	goroutinesCount int
	hashGetter      getter.HashGetter
}

func NewMaster(goroutinesCount int, hash getter.HashGetter) *Master {
	return &Master{
		goroutinesCount: goroutinesCount,
		hashGetter:      hash,
	}
}

func (m *Master) ProcessTasks(ctx context.Context, urls []string) map[string]string {
	resultCh := make(chan []string, m.goroutinesCount)
	tasksCh := make(chan string, len(urls))

	for i := 0; i < m.goroutinesCount; i++ {
		w := worker.NewWorker(m.hashGetter)
		go w.Consume(ctx, tasksCh, resultCh)
	}

	go func(urls []string, taskCh chan string) {
		for _, url := range urls {
			tasksCh <- url
		}
		close(tasksCh)
	}(urls, tasksCh)

	result := make(map[string]string, len(urls))

	tasksCount := len(urls)
	for {
		select {
		case r := <-resultCh:
			if len(r) == 2 {
				result[r[0]] = r[1]
			}

			tasksCount--
			if tasksCount == 0 {
				return result
			}
		case <-ctx.Done():
			return result
		}
	}
}
