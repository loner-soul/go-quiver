package cache

import (
	"sync"
)

type safeCache[T comparable, K any] struct {
	data sync.Map
	fn   func(T) (K, error)
}

func (c *safeCache[T, K]) Get(key T) (K, error) {
	v, ok := c.data.Load(key)
	if !ok {
		newValue, err := c.fn(key)
		if err != nil {
			return newValue, err
		}
		v = newValue
		c.data.Store(key, v)
	}
	return v.(K), nil
}

func (c *safeCache[T, K]) MustGet(key T) K {
	v, err := c.Get(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *safeCache[T, K]) Clear() {
	c.data = sync.Map{}
}

func (c *safeCache[T, K]) Delete(key T) {
	c.data.Delete(key)
}
