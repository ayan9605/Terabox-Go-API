package utils

import (
	"sync"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration int64
}

var (
	cache     = make(map[string]cacheItem)
	cacheMu   sync.RWMutex
	cacheTTL  time.Duration
	cleanupMu sync.Once
)

// InitCache initializes the cache with a TTL
func InitCache(ttl time.Duration) {
	cacheTTL = ttl
	cleanupMu.Do(func() {
		go cleanupCache()
	})
}

// SetCache stores a value in cache
func SetCache(key string, value interface{}) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cache[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(cacheTTL).Unix(),
	}
}

// GetFromCache retrieves a value from cache
func GetFromCache(key string) (interface{}, bool) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()

	item, found := cache[key]
	if !found {
		return nil, false
	}

	// Check expiration
	if time.Now().Unix() > item.expiration {
		return nil, false
	}

	return item.value, true
}

// cleanupCache removes expired items every minute
func cleanupCache() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cacheMu.Lock()
		now := time.Now().Unix()
		for key, item := range cache {
			if now > item.expiration {
				delete(cache, key)
			}
		}
		cacheMu.Unlock()
	}
}
