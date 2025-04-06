package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]CacheEntry
	mu       sync.RWMutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheMap[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ce, ok := c.cacheMap[key]
	if !ok {
		return nil, false
	}
	return ce.val, true
}

func (c *Cache) ReapLoop(interval time.Duration) {
	ticker := time.NewTicker(time.Duration(100) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.cacheMap {
			if now.Sub(v.createdAt) > interval {
				delete(c.cacheMap, k)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cacheMap: make(map[string]CacheEntry),
	}
	go cache.ReapLoop(interval)
	return cache
}
