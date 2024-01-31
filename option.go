package ego

import (
	"fmt"
	"runtime/debug"
)

type OptionFunc func(e *Ego)

const (
	PANIC_SET_SIZE_VALUE = "size must > 0"
)

func WithSize(size int64) OptionFunc {
	if size < 1 {
		panic(PANIC_SET_SIZE_VALUE)
	}
	return func(e *Ego) {
		e.size = size
	}
}

func WithQueue(queue JobQueue) OptionFunc {
	return func(e *Ego) {
		e.jobs = queue
	}
}

func WithRecover(recoverFunc func()) OptionFunc {
	return func(e *Ego) {
		e.recoverFunc = recoverFunc
	}
}

func defaultRecover() {
	if err := recover(); err != nil {
		fmt.Println("Recovered from panic:", err)
		debug.PrintStack()
	}
}
