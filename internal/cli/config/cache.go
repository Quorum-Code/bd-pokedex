package config

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

type Cache struct {
	refreshInterval time.Duration
	Entries         map[string]*cacheEntry
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
			for k, v := range c.Entries {
				if time.Since(v.createdAt) > c.refreshInterval {
					delete(c.Entries, k)
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
	c.Entries[s] = &e
}

func (c *Cache) Get(url string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.Entries[url]
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode > 299 {
			return nil, errors.New("failed response")
		} else if err != nil {
			return nil, err
		}

		go c.Add(url, body)
		return body, nil
	}
	return v.val, nil
}
