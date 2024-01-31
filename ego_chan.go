package ego

import (
	"sync/atomic"
)

// 阻塞式队列
type jobChan struct {
	jobs  chan job
	count atomic.Int64
}

func newJobChan() *jobChan {
	return &jobChan{
		jobs: make(chan job),
	}
}

func (c *jobChan) EnQueue(job job) {
	c.count.Add(1)
	c.jobs <- job
}

func (c *jobChan) DeQueue() (job, bool) {
	j, ok := <-c.jobs
	c.count.Add(-1)
	return j, ok
}

func (c *jobChan) Close() {
	close(c.jobs)
}

func (c *jobChan) Len() int64 {
	return c.count.Load()
}
