package worker

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"go-learning/practise/gin-practise/workqueue"
)

type WorkerInterface interface {
	Run(workers int, stopCh <-chan struct{})
}

type Worker struct {
	queue workqueue.Interface
}

var Queue = workqueue.NewQueue()

func NewWorker() WorkerInterface {
	return &Worker{
		queue: Queue,
	}
}

func (w *Worker) Run(workers int, stopCh <-chan struct{}) {
	for i := 0; i < workers; i++ {
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
