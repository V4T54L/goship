package cache

import (
	"sync"
	"time"
)

// SimpleCache is a basic in-memory cache with expiration.
type SimpleCache struct {
	items map[string]item
	mu    sync.RWMutex
}

type item struct {
	value      string
	expiration int64
}

// NewSimpleCache creates a new in-memory cache.
func NewSimpleCache() *SimpleCache {
	return &SimpleCache{
		items: make(map[string]item),
	}
}

// Set adds a key-value pair with TTL (in seconds).
func (c *SimpleCache) Set(key, value string, ttl int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = item{
		value:      value,
		expiration: time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	}
}

// Get retrieves a value by key. Returns value and a boolean indicating existence.
func (c *SimpleCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	it, found := c.items[key]
	if !found || time.Now().Unix() > it.expiration {
		return "", false
	}
	return it.value, true
}
