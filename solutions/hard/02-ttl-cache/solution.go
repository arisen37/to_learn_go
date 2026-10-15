package lab

import (
	"sync"
	"time"
)

type ttlEntry struct {
	value   string
	expires time.Time
}
type Cache struct {
	mu    sync.Mutex
	items map[string]ttlEntry
	now   func() time.Time
}

func NewCache(now func() time.Time) *Cache {
	if now == nil {
		now = time.Now
	}
	return &Cache{items: make(map[string]ttlEntry), now: now}
}
func (c *Cache) Set(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ttl <= 0 {
		delete(c.items, key)
		return
	}
	c.items[key] = ttlEntry{value, c.now().Add(ttl)}
}
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[key]
	if !ok {
		return "", false
	}
	if !c.now().Before(e.expires) {
		delete(c.items, key)
		return "", false
	}
	return e.value, true
}
func (c *Cache) Prune() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	n := 0
	for k, e := range c.items {
		if !now.Before(e.expires) {
			delete(c.items, k)
			n++
		}
	}
	return n
}
