package cache

type unsafeCache[T comparable, K any] struct {
	data map[T]K
	fn   func(T) (K, error)
}

func (c *unsafeCache[T, K]) Get(key T) (K, error) {
	var err error
	v, ok := c.data[key]
	if !ok {
		v, err = c.fn(key)
		if err != nil {
			return v, err
		}
		c.data[key] = v
	}
	return v, nil
}

func (c *unsafeCache[T, K]) MustGet(key T) K {
	v, err := c.Get(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *unsafeCache[T, K]) Clear() {
	c.data = make(map[T]K)
}

func (c *unsafeCache[T, K]) Delete(key T) {
	delete(c.data, key)
}
