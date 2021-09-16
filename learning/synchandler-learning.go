package main

import (
	"fmt"
	"log"
	"time"

	"k8s.io/utils/clock"
)

type TestController struct {
	queue []string

	syncHandler func(rsKey string) error
}

func (tc *TestController) syncTestHandler(key string) error {
	startTime := time.Now()
	defer func() {
		fmt.Println(key, time.Since(startTime))
	}()

	var item string
	if len(tc.queue) > 0 {
		item, tc.queue = tc.queue[0], tc.queue[1:]
		log.Println("Process item ", item)
	}
	return nil
}

func (tc *TestController) AddItems() {
	for {
		time.Sleep(1 * time.Second)
		ts := time.Now().String()
		tc.queue = append(tc.queue, ts)
		fmt.Println("AddItems", tc.queue)

	}
}

func (tc *TestController) processNextWorkItem() bool {
	tc.syncHandler("processNextWorkItem")
	return true
}

func (tc *TestController) worker() {
	for tc.processNextWorkItem() {
	}
}

func Until(f func(), stopCh <-chan struct{}) {
	var t clock.Timer
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		func() {
			f()
		}()

		select {
		case <-stopCh:
			return
		case <-t.C():
		}
	}
}

func main() {
	stopCh := make(chan struct{})

	tc := TestController{
		queue: make([]string, 0),
	}

	tc.syncHandler = tc.syncTestHandler

	go tc.AddItems()

	for i := 0; i < 5; i++ {
		go Until(tc.worker, stopCh)
	}

	<-stopCh
}
