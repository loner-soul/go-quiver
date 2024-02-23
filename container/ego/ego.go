package ego

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/loner-soul/go-quiver/container/syncx"
)

type FuncArgs func(ctx context.Context, args ...any)

type Func func(ctx context.Context)

type Ego struct {
	wg   sync.WaitGroup
	cond *sync.Cond
	// 处理异常函数
	recoverFunc func(any)
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
	if eg.recoverFunc == nil {
		eg.recoverFunc = defaultRecover
	}
	eg.cond = sync.NewCond(syncx.NewSpinLock())
	return eg
}

// Runf 当任务队列满了会阻塞
func (e *Ego) Runf(ctx context.Context, task FuncArgs, args ...any) {
	if e.isDone {
		panic("can not run any task after done")
	}
	e.cond.L.Lock()
	e.wg.Add(1)
	for e.size > 0 && e.count >= e.size {
		e.cond.Wait()
	}
	e.count++
	e.cond.L.Unlock()

	go func() {
		defer func() {
			atomic.AddInt64(&e.count, -1)
			e.wg.Done()
			if v := recover(); v != nil {
				e.recoverFunc(v)
			}
			e.cond.Signal()
		}()
		task(ctx, args...)
	}()
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

func (e *Ego) SetEgoSize(size int64) {
	atomic.StoreInt64(&e.size, size)
}
