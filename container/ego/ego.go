package ego

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/loner-soul/ego/container/echan"
)

const (
	DEFAULT_EGO_SIZE = 1000
)

type FuncArgs func(ctx context.Context, args ...any)

type Func func(ctx context.Context)

type Ego struct {
	wg sync.WaitGroup
	// 处理异常函数
	recoverFunc func()
	// 缓存任务队列
	jobs JobQueue
	// goroutine计数器
	count atomic.Int64
	// 最大goroutine数量
	size int64
}

func New(opt ...OptionFunc) *Ego {
	eg := &Ego{}
	for _, o := range opt {
		o(eg)
	}
	if eg.size == 0 {
		eg.size = DEFAULT_EGO_SIZE
	}
	if eg.jobs == nil {
		eg.jobs = echan.newJobChan()
	}
	if eg.recoverFunc == nil {
		eg.recoverFunc = defaultRecover
	}
	go eg.loopQueue()
	return eg
}

// Runf 当任务队列满了会阻塞
func (e *Ego) Runf(ctx context.Context, task FuncArgs, args ...any) {
	e.wg.Add(1) // 提前加1避免Wait时候没加上
	job := NewJob(ctx, task, args...)
	for {
		count := e.count.Load()
		if count >= e.size {
			e.jobs.Push(job)
			return
		}
		// 计数器+1
		if e.count.CompareAndSwap(count, count+1) {
			e.runJob(job)
			break
		}
	}
}

// Run 等价于 Runf 不传参数
func (e *Ego) Run(ctx context.Context, task Func) {
	e.Runf(ctx, func(ctx context.Context, args ...any) {
		task(ctx)
	})
}

// Close 等待所有任务执行完成，需要确保所有任务都调用后才执行
// 在http服务中使用时，应在server.Close()之后调用
func (e *Ego) Close() {
	e.wg.Wait()
}

func (e *Ego) runJob(job Job) {
	go func() {
		defer func() {
			// 顺序问题
			e.recoverFunc()
			e.wg.Done()
			e.count.Add(-1)
		}()
		job.f(job.ctx, job.args...)
	}()
}

func (e *Ego) loopQueue() {
	for {
		// 获取job
		job, ack, ok := e.jobs.Pop()
		if !ok {
			return
		}
		// 处理job
		for {
			state := e.count.Load()
			if state >= e.size {
				continue
			}
			// 计数器+1
			if e.count.CompareAndSwap(state, state+1) {
				e.runJob(job)
				break
			}
		}
		// 确认消费
		if ack != nil {
			ack()
		}
	}
}

func (e *Ego) greaterThanOrEqualToSize() bool {
	state := e.count.Load()
	return state >= e.size
}
