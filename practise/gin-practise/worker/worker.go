package worker

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"go-learning/practise/gin-practise/workqueue"
)

type WorkerInterface interface {
	Run(workers int, stopCh <-chan struct{})

	DoTest(ctx context.Context, s string) error

	DoAfterTest(ctx context.Context, s string) error
}

type Worker struct {
	queue workqueue.DelayingInterface
}

func NewWorker() WorkerInterface {
	return &Worker{
		queue: workqueue.NewQueue(),
	}
}

func (w *Worker) Run(workers int, stopCh <-chan struct{}) {
	for i := 0; i < workers; i++ {
		// 可启动多个协程
		go wait.Until(w.worker, 2*time.Second, stopCh)
	}
}

func (w *Worker) worker() {
	item := w.queue.Get()
	if len(item) == 0 {
		return
	}

	fmt.Println("item:", item, "length:", w.queue.Len())
}

func (w *Worker) DoTest(ctx context.Context, s string) error {
	// TODO: do somethings
	w.queue.Add(s)
	return nil
}

// 延迟5秒
func (w *Worker) DoAfterTest(ctx context.Context, s string) error {

	w.queue.AddAfter(s, 5)
	return nil
}
