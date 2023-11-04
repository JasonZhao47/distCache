package lru

import (
	"container/list"
)

// LRUCache not safe for concurrency
type LRUCache[T comparable, V interface{}] interface {
	Add(T, V)
	Get(T) (V, bool)
	Remove(T)
	RemoveOldest()
	Peek(T) (V, bool)
	Len() int
}

type Key interface{}

// LRU implementation
type LRU[T comparable, V interface{}] struct {
	ll         *list.List
	cache      map[T]*list.Element
	MaxEntries int
}

func New[T comparable, V interface{}](maxEntries int) (LRUCache[T, V], error) {
	return &LRU[T, V]{
		ll:         list.New(),
		cache:      make(map[T]*list.Element),
		MaxEntries: maxEntries,
	}, nil
}

func (L *LRU[T, V]) Len() int {
	return len(L.cache)
}

func (L *LRU[T, V]) Add(key T, value V) {
	// check if lru exists
	if L.cache == nil {
		return
	}
	// check if key is inside lru
	elem, ok := L.cache[key]
	if !ok {
		// if not, add and move to front
		e := &entry{
			key:   key,
			value: value,
		}
		if L.Len()+1 > L.MaxEntries {
			L.RemoveOldest()
		}
		elem := L.ll.PushFront(e)
		L.cache[key] = elem

	} else {
		elem.Value.(*entry).value = value
		L.ll.MoveToFront(elem)
	}
	return
}

func (L *LRU[T, V]) Get(key T) (value V, evicted bool) {
	// check if lru exists
	if L.cache == nil {
		return
	}
	// check if t is in lru
	elem, ok := L.cache[key]
	if ok {
		L.ll.MoveToFront(elem)
		return elem.Value.(*entry).value.(V), ok
	}
	return
}

func (L *LRU[T, V]) Remove(key T) {
	// check if lru exists
	if L.cache == nil {
		return
	}
	// if element doesn't exist,
	// nothing happens
	// if exists, removes it,
	// no need to change internal order
	elem, ok := L.cache[key]
	if ok {
		L.removeElement(elem)
	}
	return
}

func (L *LRU[T, V]) RemoveOldest() {
	if L.cache == nil || L.ll.Len() == 0 {
		return
	}
	// remove the oldest elem from both hashtable and linked-list.
	// if both are empty just return
	elem := L.ll.Back()
	L.removeElement(elem)
}

func (L *LRU[T, V]) Peek(key T) (value V, evicted bool) {
	// if lru is nil, return default values
	if L.cache == nil {
		return
	}
	// peek without changing lru's order
	elem, ok := L.cache[key]
	if ok {
		return elem.Value.(*entry).value.(V), ok
	}
	return
}

type entry struct {
	key   Key
	value interface{}
}

func (L *LRU[T, V]) removeElement(elem *list.Element) {
	rmKey := elem.Value.(*entry).key
	delete(L.cache, rmKey.(T))
	L.ll.Remove(elem)
}
