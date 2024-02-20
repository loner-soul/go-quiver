package echan

import (
	"sync/atomic"

	"github.com/loner-soul/go-quiver/container/queue"
)

type chanQueue[T any] struct {
	c     chan T
	count atomic.Int64
}

func NewEChan[T any](size int64) queue.Queue[T] {
	return &chanQueue[T]{
		c: make(chan T, size),
	}
}

// Push 添加一个
func (c *chanQueue[T]) Push(job T) {
	c.count.Add(1)
	c.c <- job
}

func (c *chanQueue[T]) Pop() (T, bool) {
	j, ok := <-c.c
	if !ok {
		return j, false
	}
	c.count.Add(-1)
	return j, true
}

func (c *chanQueue[T]) Len() int64 {
	return c.count.Load()
}
