package distCache

import (
	"github.com/jasonzhao47/distCache/internal/lru"
	"sync"
)

// question: how to test concurrency?

type cache[T comparable, V interface{}] struct {
	lru        lru.LRUCache[T, V]
	mu         sync.Mutex
	cacheBytes int
}

func (c *cache[T, V]) add(key T, value V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		lruCache, err := lru.New[T, V](c.cacheBytes)
		if err != nil {
			return err
		}
		c.lru = lruCache
	}

	c.lru.Add(key, value)
	return nil
}

func (c *cache[int, interface{}]) get(key int) (value int, evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(int), ok
	}
	return
}
