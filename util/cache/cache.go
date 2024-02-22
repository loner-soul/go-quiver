package cache

import (
	"sync"
)

type Cache[T comparable, K any] interface {
	Get(T) (K, error)
	MustGet(T) K

	Delete(T) // 删除一个缓存
	Clear()   // 清空所有缓存
}

func NewUnsafeCache[T comparable, K any](fn func(T) (K, error)) Cache[T, K] {
	return &unsafeCache[T, K]{
		data: make(map[T]K),
		fn:   fn,
	}
}

func NewSafeCache[T comparable, K any](fn func(T) (K, error)) Cache[T, K] {
	return &safeCache[T, K]{
		data: sync.Map{},
		fn:   fn,
	}
}
