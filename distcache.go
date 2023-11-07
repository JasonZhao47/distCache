package distCache

import (
	"errors"
	"github.com/jasonzhao47/distCache/internal/lru"
	"sync"
)

type GroupCache struct {
	getter func(byte, []byte) error
	// getter func - callback function for data source
	// get the real data out from it
	name  string
	cache cache[byte, ByteView]
}

func (c *GroupCache) GetFromLocal(key byte) (ByteView, error) {
	//TODO implement me
	panic("implement me")
}

func (c *GroupCache) GetFromPeer(peerName string, key byte) (ByteView, error) {
	//TODO implement me
	panic("implement me")
}

var (
	ErrInternalError = errors.New("Internal error")
	// global variable - shared memory!
	groups map[string]*GroupCache
	mu     sync.RWMutex
)

// DistCache question: how to test concurrency?
type DistCache interface {
	GetGroup(name string) (*GroupCache, error)
	Get(name string, key byte) (ByteView, error)
	GetFromPeer(peerName string, key byte) (ByteView, error)
	GetFromLocal(key byte) (ByteView, error)
}

func NewGroups(getter func(byte, []byte) error,
	name string, cacheBytes int) DistCache {
	if getter == nil {
		panic("empty data getter")
	}
	mu.Lock()
	defer mu.Unlock()
	gc := &GroupCache{
		getter: getter,
		name:   name,
		cache:  cache[byte, ByteView]{cacheBytes: cacheBytes},
	}
	groups[name] = gc
	return gc
}

func (c *GroupCache) GetGroup(name string) (*GroupCache, error) {
	mu.RLock()
	defer mu.RUnlock()
	group, ok := groups[name]
	if !ok {
		return nil, ErrInternalError
	}
	return group, nil
}

func (c *GroupCache) Get(name string, key byte) (ByteView, error) {
	g, err := c.GetGroup(name)
	if err != nil {
		return NewByteView([]byte{}), err
	}
	val, ok := g.cache.get(key)
	if !ok {
		// val, peerOk := GetFromPeer()
		// if !peerOk
		// How to use this getter?
		//err := c.getter(key, data)
		//if err != nil {
		//	return nil, err
		//}
		//c.cache.add(key, data)
		// problem: determine which node to emplace this data?
	}
	return val, nil
}

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
