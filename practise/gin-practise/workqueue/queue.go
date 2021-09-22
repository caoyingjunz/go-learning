package workqueue

import (
	"sync"
)

type Interface interface {
	Add(item string)
	Len() int
	Get() string
}

type Type struct {
	// queue defines the order in which we will work on items.
	queue []string

	lock sync.RWMutex
}

// Add marks item as needing processing.
func (q *Type) Add(item string) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queue = append(q.queue, item)
}

// Len returns the current queue length, for informational purposes only. You
// shouldn't e.g. gate a call to Add() or Get() on Len() being a particular
// value, that can't be synchronized properly.
func (q *Type) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.queue)
}

func (q *Type) Get() string {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.queue) == 0 {
		return ""
	}

	var item string
	item, q.queue = q.queue[0], q.queue[1:]

	return item
}

func NewQueue() Interface {
	return &Type{
		queue: []string{},
	}
}
