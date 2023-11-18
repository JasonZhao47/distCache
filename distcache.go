package distCache

import (
	"errors"
	"github.com/jasonzhao47/distCache/internal/lru"
	"log"
	"sync"
)

type Getter interface {
	Get(string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (fn GetterFunc) Get(key string) ([]byte, error) {
	return fn(key)
}

type GroupCache struct {
	getter Getter
	// getter func - callback function for data source
	// get the real data out from it
	name  string
	cache cache[string, ByteView]
}

func (c *GroupCache) loadFromLocal(key string) (ByteView, error) {
	val, err := c.getter.Get(key)
	if err != nil {
		return NewByteView([]byte{}), err
	}
	buf := NewByteView(val)
	c.cache.add(key, buf)
	return buf, nil
}

func (c *GroupCache) getFromPeer(peerName string, key byte) (ByteView, error) {
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
	Get(name string, key string) (ByteView, error)
}

func NewGroups(getter Getter, name string, cacheBytes int) DistCache {
	if getter == nil {
		panic("empty data getter")
	}
	mu.Lock()
	defer mu.Unlock()
	gc := &GroupCache{
		getter: getter,
		name:   name,
		cache:  cache[string, ByteView]{cacheBytes: cacheBytes},
	}
	groups[name] = gc
	return gc
}

func GetGroup(name string) (*GroupCache, error) {
	mu.RLock()
	defer mu.RUnlock()
	group, ok := groups[name]
	if !ok {
		return nil, ErrInternalError
	}
	return group, nil
}

func (c *GroupCache) Get(name string, key string) (ByteView, error) {
	g, err := GetGroup(name)
	if err != nil {
		return NewByteView([]byte{}), err
	}
	val, ok := g.cache.get(key)
	if !ok {
		log.Printf("[Cache] not hit %v\n", val)
		// val, peerOk := GetFromPeer()
		// if !peerOk
		// How to use this getter?
		data, err := c.loadFromLocal(key)
		if err != nil {
			return nil, err
		}
		return data, nil
		// problem: determine which node to emplace this data?
	}
	log.Printf("[Cache] hit %v\n", val)
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
