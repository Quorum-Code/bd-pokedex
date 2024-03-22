package cli

import (
	"errors"
	"time"
)

type Cache struct {
	refreshInterval time.Duration
	entries         map[string]*cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(refreshInterval time.Duration) Cache {
	c := Cache{refreshInterval, map[string]*cacheEntry{}}
	c.reapLoop(500 * time.Millisecond)
	return c
}

func (c *Cache) reapLoop(checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	check := func() {
		for {
			<-ticker.C
			for k, v := range c.entries {
				if time.Since(v.createdAt) > c.refreshInterval {
					delete(c.entries, k)
				}
			}
		}
	}
	go check()
}

func (c *Cache) Add(s string, v []byte) {
	e := cacheEntry{
		time.Now(),
		v,
	}
	c.entries[s] = &e
}

func (c *Cache) Get(s string) ([]byte, error) {
	v, ok := c.entries[s]
	if !ok {
		return nil, errors.New("no entry")
	}
	return v.val, nil
}
