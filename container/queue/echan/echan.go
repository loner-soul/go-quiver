package echan

import (
	"sync"
	"sync/atomic"

	"github.com/loner-soul/ego/container/ego"
)

type jobChan struct {
	rw    sync.RWMutex
	jobs  chan ego.Job
	count atomic.Int64
}

func newJobChan(size int64) ego.JobQueue {
	return &jobChan{
		jobs: make(chan ego.Job, size),
	}
}

// Push 添加一个
func (c *jobChan) Push(job ego.Job) {
	c.count.Add(1)
	c.jobs <- job
}

func (c *jobChan) Pop() (ego.Job, func(), bool) {
	j, ok := <-c.jobs
	if !ok {
		return j, func() {}, false
	}
	f := func() { c.count.Add(-1) }
	return j, f, true
}

func (c *jobChan) Len() int64 {
	return c.count.Load()
}

func (c *jobChan) Close() {
	close(c.jobs)
}
