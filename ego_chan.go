package ego

import (
	"sync/atomic"
)

// 阻塞式队列
type jobChan struct {
	jobs  chan Job
	close bool
	count atomic.Int64
}

func newJobChan() *jobChan {
	return &jobChan{
		jobs: make(chan Job),
	}
}

// EnQueue close后不能再写入
func (c *jobChan) EnQueue(job Job) {
	if c.close {
		// TODO logs
		return
	}
	c.count.Add(1)
	c.jobs <- job
}

func (c *jobChan) DeQueue() (Job, func(), bool) {
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
	c.close = true
	for {
		if c.Len() == 0 {
			break
		}
	}
	close(c.jobs)
}
