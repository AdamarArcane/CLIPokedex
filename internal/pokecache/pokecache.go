package pokecache

import (
	"sync"
	"time"
)

// cacheEntry represents a single cache entry with a creation timestamp and value.
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache provides a thread-safe cache with expiration functionality.
type Cache struct {
	entries  map[string]cacheEntry
	mutex    sync.Mutex
	interval time.Duration
}

// NewCache creates a new Cache instance with a configurable interval.
// The interval defines how often the cache cleans up expired entries.
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}
	go c.reapLoop()
	return c
}

// Add inserts a new entry into the cache with the given key and value.
func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves an entry from the cache by key.
// It returns the value and a boolean indicating whether the key was found.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}

// reapLoop periodically removes expired entries from the cache.
// It runs in a separate goroutine started by NewCache.
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.mutex.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.Sub(entry.createdAt) > c.interval {
				delete(c.entries, key)
			}
		}
		c.mutex.Unlock()
	}
}
