package ego

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
)

const (
	BenchParam = 10
	n          = 1000_0000
)

func Test_EgoUse(t *testing.T) {
	pool := New()

	pool.Runf(context.Background(), func(ctx context.Context, args ...any) {
		fmt.Println("call in goroutine")
	})
	pool.Run(context.Background(), func(ctx context.Context) {
		fmt.Println("call in goroutine")
	})

	pool.Done()
	fmt.Println("done")
}

func demoFunc() {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}

var curMem uint64

func TestNoPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			demoFunc()
			wg.Done()
		}()
	}

	wg.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestAntsPool(t *testing.T) {
	p, _ := ants.NewPool(10000)
	defer p.Release()
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		err := p.Submit(func() {
			demoFunc()
			wg.Done()
		})
		if err != nil {
			t.Logf("submit error: %v", err)
		}
	}
	wg.Wait()

	t.Logf("pool, capacity:%d", p.Cap())
	t.Logf("pool, running workers number:%d", p.Running())
	t.Logf("pool, free workers number:%d", p.Free())

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestEgoPool(t *testing.T) {
	// defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	for i := 0; i < n; i++ {
		Run(context.Background(), func(ctx context.Context) {
			demoFunc()
		})
	}
	Done()

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}
