package cache_test

import (
	"github.com/jasonzhao47/distCache/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLRU_Peek(t *testing.T) {
	testCases := []struct {
		name         string
		mock         func() cache.LRUCache[int, int]
		key          int
		value        int
		isKeyPresent bool
	}{
		{
			name: "Should get value",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				for i := 1; i <= 2; i++ {
					lruCache.Add(i, 8777)
				}
				for i := 0; i <= 5; i++ {
					lruCache.Get(1)
					lruCache.Get(2)
				}
				return lruCache
			},
			key:          1,
			value:        8777,
			isKeyPresent: true,
		},
		{
			name: "Should get value without moving internal order",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				for i := 1; i <= 2; i++ {
					lruCache.Add(i, 8777)
				}
				for i := 0; i <= 5; i++ {
					lruCache.Get(1)
					lruCache.Get(2)
				}
				lruCache.Add(3, 29)
				return lruCache
			},
			key:          1,
			value:        0,
			isKeyPresent: false,
		},
		{
			name: "Should not get value when not present",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				for i := 1; i <= 2; i++ {
					lruCache.Add(i, 8777)
				}
				for i := 0; i <= 5; i++ {
					lruCache.Get(1)
					lruCache.Get(2)
				}
				return lruCache
			},
			key:          3,
			value:        0,
			isKeyPresent: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lruCache := tc.mock()
			val, ok := lruCache.Peek(tc.key)
			assert.Equal(t, tc.isKeyPresent, ok)
			assert.Equal(t, tc.value, val)
		})
	}
}

func TestLRU_Remove(t *testing.T) {
	testCases := []struct {
		name         string
		mock         func() cache.LRUCache[int, int]
		key          int
		value        int
		isKeyPresent bool
	}{
		{
			name: "Should remove key & value",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				for i := 1; i <= 2; i++ {
					lruCache.Add(i, 8777)
				}
				for i := 0; i <= 5; i++ {
					lruCache.Get(1)
					lruCache.Get(2)
				}
				return lruCache
			},
			key:          1,
			value:        0,
			isKeyPresent: false,
		},
		{
			name: "Should not remove value when not present",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				for i := 1; i <= 2; i++ {
					lruCache.Add(i, 8777)
				}
				for i := 0; i <= 5; i++ {
					lruCache.Get(1)
					lruCache.Get(2)
				}
				return lruCache
			},
			key:          3,
			value:        0,
			isKeyPresent: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lruCache := tc.mock()
			lruCache.Remove(tc.key)
			val, ok := lruCache.Peek(tc.key)
			assert.Equal(t, tc.isKeyPresent, ok)
			assert.Equal(t, tc.value, val)
		})
	}
}

func TestLRUCache_RemoveOldest(t *testing.T) {
	testCases := []struct {
		name string
		// no need to mock anything here
		mock         func() cache.LRUCache[int, int]
		key          int
		isKeyPresent bool
	}{
		{
			name: "Should not remove anything when cache is empty",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				return lruCache
			},
			key:          1,
			isKeyPresent: false,
		},
		{
			name: "Pop out element",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				for i := 1; i <= 2; i++ {
					lruCache.Add(i, 8777)
				}
				for i := 0; i <= 5; i++ {
					lruCache.Get(1)
					lruCache.Get(2)
				}
				return lruCache
			},
			key:          1,
			isKeyPresent: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lruCache := tc.mock()
			lruCache.RemoveOldest()
			_, ok := lruCache.Peek(tc.key)
			assert.Equal(t, tc.isKeyPresent, ok)
		})
	}

}

func TestLRUCache_Get(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func() cache.LRUCache[int, int]
		key     int
		value   int
		present bool
	}{
		{
			name: "Should get an existing element and set to recently used",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				lruCache.Add(1235, 1928)
				return lruCache
			},
			key:     1235,
			value:   1928,
			present: true,
		},
		{
			name: "Should not get element when not existing",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				lruCache.Add(1235, 1928)
				return lruCache
			},
			key:     1,
			value:   0,
			present: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lruCache := tc.mock()
			val, ok := lruCache.Get(tc.key)
			assert.Equal(t, tc.value, val)
			assert.Equal(t, tc.present, ok)
		})
	}

}

func TestLRUCache_Add(t *testing.T) {
	testCases := []struct {
		name         string
		mock         func() cache.LRUCache[int, int]
		recentlyUsed int
		isInserted   bool
	}{
		{
			name: "Should add element as most recently used",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				lruCache.Add(1123, 1928)
				return lruCache
			},
			recentlyUsed: 29385,
			isInserted:   true,
		},
		{
			name: "Should update element used when key exists",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				for i := 1; i <= 2; i++ {
					lruCache.Add(i, 917)
				}
				return lruCache
			},
			recentlyUsed: 29385,
			isInserted:   true,
		},
		{
			name: "Should add and remove oldest if maxEntry exceeds",
			mock: func() cache.LRUCache[int, int] {
				lruCache, err := cache.New[int, int](2)
				require.NoError(t, err)
				for i := 3; i <= 4; i++ {
					lruCache.Add(i, 917)
				}
				return lruCache
			},
			recentlyUsed: 29385,
			isInserted:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lruCache := tc.mock()
			lruCache.Add(1, 29385)
			val, ok := lruCache.Peek(1)
			assert.Equal(t, tc.isInserted, ok)
			assert.Equal(t, tc.recentlyUsed, val)
		})
	}
}
