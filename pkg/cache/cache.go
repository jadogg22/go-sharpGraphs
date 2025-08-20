package cache

import (
	"sync"
	"time"
)

// CacheItem holds the value and expiry time for a cached item.
type CacheItem struct {
	Value  interface{}
	TypeID string
	Expiry time.Time
}

// Cache is a thread-safe cache that can store items of any type.
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

var MyCache = NewCache()

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// Set adds an item to the cache with a given time-to-live (TTL).
func (c *Cache) Set(key string, value interface{}, typeID string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:  value,
		TypeID: typeID,
		Expiry: time.Now().Add(ttl),
	}
}

// Get retrieves an item from the cache. It returns the value as an interface{} and the typeID.
func (c *Cache) Get(key string) (interface{}, string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || time.Now().After(item.Expiry) {
		return nil, "", false
	}
	return item.Value, item.TypeID, true
}

// Delete removes an item from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}
