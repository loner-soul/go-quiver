package ego

import (
	"fmt"
	"runtime/debug"
)

type OptionFunc func(e *Ego)

const (
	PANIC_SET_SIZE_VALUE = "ego size must > 0"
)

func WithSize(size int64) OptionFunc {
	if size < 1 {
		panic(PANIC_SET_SIZE_VALUE)
	}
	return func(e *Ego) {
		e.size = size
	}
}

func WithRecover(recoverFunc func(any)) OptionFunc {
	return func(e *Ego) {
		e.recoverFunc = recoverFunc
	}
}

func defaultRecover(v any) {
	fmt.Println("Recovered from panic:", v)
	debug.PrintStack()
}
