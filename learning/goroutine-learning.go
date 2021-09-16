package main

import (
	"fmt"
	"sync"

	"github.com/go-basic/uuid"
)

// sync + chan 控制并发

type Store struct {
	items map[string]int
	lock  sync.RWMutex
}

func (s *Store) Add(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	value, exist := s.items[key]
	if !exist {
		s.items[key] = 1
		return
	}

	s.items[key] = value + 1
}

func main() {
	count := 1000

	var wg sync.WaitGroup
	s := &Store{
		items: make(map[string]int),
	}

	// 控制并发的 chan
	ch := make(chan struct{}, 200)

	wg.Add(count)
	for i := 0; i < count; i++ {
		ch <- struct{}{}
		go func(i int, s *Store) {
			defer wg.Done()
			s.Add(uuid.New())
			<-ch
		}(i, s)
	}
	wg.Wait()

	for k, v := range s.items {
		// k 和 v 分别是 生成的 uuid 和 出现的次数
		fmt.Println(k, v)
	}
}
