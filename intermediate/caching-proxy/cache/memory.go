package cache

import (
	"fmt"
	"net/http"
	"sync"
)

type CachedResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

type Store struct {
	mu    sync.RWMutex
	store map[string]CachedResponse
}

func NewStore() *Store {
	return &Store{
		store: make(map[string]CachedResponse),
	}
}

func (c *Store) Get(key string) (CachedResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res, found := c.store[key]
	return res, found
}

func (c *Store) Set(key string, res CachedResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = res
}

func (c *Store) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]CachedResponse)
	fmt.Println("Cache cleared in memory")
}
