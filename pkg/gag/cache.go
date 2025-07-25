package gag

import "sync"

type UpdateCallback[T any] func(T)

type Cache[T any] struct {
	mu    sync.RWMutex
	cache T
	cb    UpdateCallback[T]
}

func NewCache[T any](cb UpdateCallback[T]) *Cache[T] {
	return &Cache[T]{
		cache: *new(T),
		cb:    cb,
	}
}

func (c *Cache[T]) Get() T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cache
}

func (c *Cache[T]) Set(value T) {
	func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.cache = value
	}()

	if c.cb != nil {
		c.cb(value)
	}
}
