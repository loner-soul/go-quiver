package ego

import (
	"context"
	"sync"
	"sync/atomic"
)

const (
	DEFAULT_EGO_SIZE = 1
)

type FuncArgs func(ctx context.Context, args ...any)

type Func func(ctx context.Context)

type Ego struct {
	wg   sync.WaitGroup
	cond *sync.Cond
	// 处理异常函数
	recoverFunc func()
	// goroutine计数器
	count int64
	// 最大goroutine数量
	size   int64
	isDone bool
}

func New(opt ...OptionFunc) *Ego {
	eg := &Ego{}
	for _, o := range opt {
		o(eg)
	}
	if eg.size == 0 {
		eg.size = DEFAULT_EGO_SIZE
	}
	if eg.recoverFunc == nil {
		eg.recoverFunc = defaultRecover
	}
	eg.cond = sync.NewCond(&sync.Mutex{})
	return eg
}

// Runf 当任务队列满了会阻塞
func (e *Ego) Runf(ctx context.Context, task FuncArgs, args ...any) {
	if e.isDone {
		panic("can not run any task after done")
	}

	e.wg.Add(1)
	e.cond.L.Lock()
	for e.size >= e.count {
		e.cond.Wait()
	}
	e.count++
	e.cond.L.Unlock()

	go func(ctx context.Context, task FuncArgs, args ...any) {
		defer func() {
			e.recoverFunc()
			e.wg.Done()
			atomic.AddInt64(&e.count, -1)
		}()
		task(ctx, args...)
	}(ctx, task, args...)

}

// Run 等价于 Runf 不传参数
func (e *Ego) Run(ctx context.Context, task Func) {
	e.Runf(ctx, func(ctx context.Context, args ...any) {
		task(ctx)
	})
}

// Done 等待所有任务执行完成，需要确保所有任务都调用后才执行
// 在http服务中使用时，应在server.Close()之后调用
func (e *Ego) Done() {
	e.isDone = true
	e.wg.Wait()
}
