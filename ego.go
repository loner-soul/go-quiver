package ego

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

const (
	DEFAULT_EGO_SIZE = 1000
)

const (
	PANIC_SET_SIZE_AFTER_RUN = "can not set size after run"

	PANIC_ADD_TASK_AFTER_WAIT = "can not add task after wait"
)

// JobQueue 任务缓冲队列
type JobQueue interface {
	EnQueue(job job)
	DeQueue() (job, bool)
	Len() int64
	Close()
}

type job struct {
	f    FuncArgs
	ctx  context.Context
	args []any
}

type FuncArgs func(ctx context.Context, args ...any)

type Func func(ctx context.Context)

type Ego struct {
	wg   sync.WaitGroup
	jobs JobQueue // 任务队列
	// 状态管理
	run  bool // is running
	wait bool // 等待结束

	count atomic.Int64 // goroutine counter
	size  int64        // max goroutine num
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
		eg.jobs = newJobChan()
	}
	go eg.loopChan()
	return eg
}

func (e *Ego) SetSize(size int64) {
	if size < 1 {
		panic(PANIC_SET_SIZE_VALUE)
	}
	if e.run {
		panic(PANIC_SET_SIZE_AFTER_RUN)
	}
	e.size = size
}

// Runf 当任务队列满了会阻塞
func (e *Ego) Runf(ctx context.Context, task FuncArgs, args ...any) {
	// 调用eg.Wait之后不能再添加任务
	if e.wait {
		panic(PANIC_ADD_TASK_AFTER_WAIT)
	}
	// 标记运行状态
	e.run = true

	for {
		state := e.count.Load()
		if state >= e.size {
			e.jobs <- job{f: task, ctx: ctx, args: args}
			return
		}
		// 计数器+1
		if e.count.CompareAndSwap(state, state+1) {
			break
		}
	}
	e.goRun(ctx, task, args...)
}

// Run 等价于 Runf 不传参数
func (e *Ego) Run(ctx context.Context, task Func) {
	e.Runf(ctx, func(ctx context.Context, args ...any) {
		task(ctx)
	})
}

func (e *Ego) Wait() {
	e.wait = true

	// 等待所有chanel写入

	e.wg.Wait()
	close(e.jobs)
}

func (e *Ego) goRun(ctx context.Context, task FuncArgs, args ...any) {
	e.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Recovered from panic:", err)
				debug.PrintStack()
			}
			e.wg.Done()
			e.count.Add(-1)
		}()
		// TODO 自定义recover

		task(ctx, args...)
	}()
}

func (e *Ego) loopChan() {
	for {
		task, ok := <-e.jobs
		if !ok {
			// close chan
			return
		}
		for {
			state := e.count.Load()
			if state >= e.size {
				continue
			}
			// 计数器+1
			if e.count.CompareAndSwap(state, state+1) {
				break
			}
		}
		e.goRun(task.ctx, task.f, task.args...)
	}
}
