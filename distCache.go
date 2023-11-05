package distCache

import (
	"github.com/jasonzhao47/distCache/internal/lru"
	"sync"
)

// question: how to test concurrency?

// cache: wrapper for LRUCache, thread safe
type cache[T comparable, V interface{}] struct {
	lru        lru.LRUCache[T, V]
	mu         sync.Mutex
	cacheBytes int
}

func (c *cache[T, V]) add(key T, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		lruCache, err := lru.New[T, V](c.cacheBytes)
		if err != nil {
			panic(err)
		}

		c.lru = lruCache
	}

	c.lru.Add(key, value)
	return
}

func (c *cache[T, V]) get(key T) (value V, evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v, ok
	}
	return
}

func (c *cache[T, V]) removeOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.RemoveOldest()
}
