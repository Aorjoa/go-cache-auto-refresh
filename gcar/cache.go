package gcar

import "sync"

// Cache use for construct model to store cache
type Cache struct {
	mu    sync.RWMutex
	items map[string]interface{}
}

// Set should be set cache to memory
func Set() (string, bool) {
	return "Hello, World", true
}

// Get should be get cache data from memory
func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	return item, true
}
