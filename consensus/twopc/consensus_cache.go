package twopc

import (
	"github.com/hashicorp/golang-lru/simplelru"
	"sync"
)

const (
	Default2pcMsgCacheSize = 1024
)

type TwopcMsgCache struct {
	lru  *simplelru.LRU
	lock sync.RWMutex
}

func NewTwopcMsgCache(size int) (*TwopcMsgCache, error) {
	w := &TwopcMsgCache{}
	lru, err := simplelru.NewLRU(size, nil)

	if err != nil {
		return nil, err
	}
	w.lru = lru
	return w, nil
}


// Purge is used to completely clear the cache
func (c *TwopcMsgCache) Purge() {
	c.lock.Lock()
	c.lru.Purge()
	c.lock.Unlock()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *TwopcMsgCache) Add(key, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Add(key, value)
}

// Get looks up a key's value from the cache.
func (c *TwopcMsgCache) Get(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.lru.Get(key)
	if !ok {
		return nil, ok
	}
	return value, ok
}

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *TwopcMsgCache) Contains(key interface{}) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if !c.lru.Contains(key) {
		return false
	}
	return true
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *TwopcMsgCache) Peek(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.lru.Peek(key)
	return value, ok
}

// ContainsOrAdd checks if a key is in the cache  without updating the
// recent-ness or deleting it for being stale,  and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *TwopcMsgCache) ContainsOrAdd(key, value interface{}) (ok, evict bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lru.Contains(key) {
		return true, false
	} else {
		return false, c.lru.Add(key, value)
	}
}

// Remove removes the provided key from the cache.
func (c *TwopcMsgCache) Remove(key interface{}) {
	c.lock.Lock()
	c.lru.Remove(key)
	c.lock.Unlock()
}

// RemoveOldest removes the oldest item from the cache.
func (c *TwopcMsgCache) RemoveOldest() {
	c.lock.Lock()
	c.lru.RemoveOldest()
	c.lock.Unlock()
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *TwopcMsgCache) Keys() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Keys()
}

// Len returns the number of items in the cache.
func (c *TwopcMsgCache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Len()
}

