package gcar

import "sync"

// Cache use for construct model to store cache
type Cache struct {
	*cache
}

// CallFunction is type for function that we suppose to cache when process done
type CallFunction func() (interface{}, error)

type cache struct {
	mu    *sync.RWMutex
	items map[string]interface{}
}

// New is using for inital cache
func New() Cache {
	var mrw sync.RWMutex
	var items = make(map[string]interface{})
	return Cache{
		&cache{
			mu:    &mrw,
			items: items,
		},
	}
}

// Set should be set cache to memory
func (c *Cache) Set(k string, v interface{}) {
	c.mu.Lock()
	c.items[k] = v
	c.mu.Unlock()
}

// Get should be get cache data from memory
func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	c.mu.RUnlock()
	return item, true
}

// CallFunctionThenCache is using for call function that return interface{} and error
// then cache it
// !!! its NOT concurrecy safe for now !!!
func (c *Cache) CallFunctionThenCache(k string, cf CallFunction) {
	resp, err := cf()
	if err != nil {
		return
	}
	c.Set(k, resp)
}
