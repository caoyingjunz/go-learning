package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// 封装在结构体中使用 https://github.com/kubernetes/kubernetes/blob/ea0764452222146c47ec826977f49d7001b0ea8c/test/e2e/framework/log_size_monitoring.go#L173
// 单独使用
// 参考1 https://github.com/kubernetes/kubernetes/blob/9c147baa70c31afc966329df73302e9b52d8e432/test/e2e/apimachinery/flowcontrol.go#L157
// 参考2 https://github.com/kubernetes/kubernetes/blob/9c147baa70c31afc966329df73302e9b52d8e432/pkg/controller/replicaset/replica_set.go#L616

func Get(i string) error {
	time.Sleep(2 * time.Second)
	fmt.Println("i", i)
	return fmt.Errorf("test error")
}

//https://github.com/kubernetes/kubernetes/blob/ea0764452222146c47ec826977f49d7001b0ea8c/pkg/scheduler/metrics/metric_recorder_test.go#L81

func TestClear(t *testing.T) {
	var wg sync.WaitGroup
	incLoops, decLoops := 100, 80
	wg.Add(incLoops + decLoops)
	for i := 0; i < incLoops; i++ {
		go func() {
			wg.Done()
		}()
	}
	for i := 0; i < decLoops; i++ {
		go func() {
			wg.Done()
		}()
	}
	wg.Wait()
}

func main() {
	Items := []string{"1", "2", "3", "4", "5"}

	diff := len(Items)

	errCh := make(chan error, diff)
	var wg sync.WaitGroup
	wg.Add(diff)
	for _, i := range Items {
		// 参数避免使用指针，否则会出现同一个值问题
		go func(i string) {
			defer wg.Done()
			if err := Get(i); err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()

	select {
	case err := <-errCh:
		if err != nil {
			fmt.Println(err.Error())
			//return err
		}
	default:
	}
}
