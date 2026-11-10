package pokecache

import (
	"context"
	"sync"
	"time"
)

type Cache struct {
	interval time.Duration
	Data     map[string]cacheEntry
	mu       sync.RWMutex
}

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

func NewCache(ctx context.Context, interval time.Duration) *Cache {
	c := &Cache{Data: make(map[string]cacheEntry), interval: interval}
	go c.reapLoop(ctx)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e := cacheEntry{val: val, createdAt: time.Now()}
	c.Data[key] = e
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if v, ok := c.Data[key]; ok {
		return v.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.evictExpired()
		case <-ctx.Done():
			return
		}
	}
}

func (c *Cache) evictExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.Data {
		if time.Since(v.createdAt) > c.interval {
			delete(c.Data, k)
		}
	}
}
