package cli

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	refreshInterval time.Duration
	entries         map[string]*cacheEntry
	mu              *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(refreshInterval time.Duration) Cache {
	c := Cache{refreshInterval, map[string]*cacheEntry{}, &sync.Mutex{}}
	c.reapLoop(500 * time.Millisecond)
	return c
}

func (c *Cache) reapLoop(checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	check := func() {
		for {
			<-ticker.C
			c.mu.Lock()
			for k, v := range c.entries {
				if time.Since(v.createdAt) > c.refreshInterval {
					delete(c.entries, k)
				}
			}
			c.mu.Unlock()
		}
	}
	go check()
}

func (c *Cache) Add(s string, v []byte) {
	e := cacheEntry{
		time.Now(),
		v,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[s] = &e
}

func (c *Cache) Get(s string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.entries[s]
	if !ok {
		return nil, errors.New("no entry")
	}
	return v.val, nil
}
