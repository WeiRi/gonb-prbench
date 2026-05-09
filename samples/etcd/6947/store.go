package cache

import (
	"container/list"
	"sync"
)

type lru struct {
	items map[string]*list.Element
	order *list.List
	cap   int
}

type Cache struct {
	mu  sync.Mutex
	lru *lru
}

func NewCache(size int) *Cache {
	return &Cache{lru: &lru{items: map[string]*list.Element{}, order: list.New(), cap: size}}
}

// Add — writes c.lru.items under c.mu (line 51 area).
func (c *Cache) Add(k string, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.items[k] = c.lru.order.PushFront(v) // line 51 write
}

// Size — BUG (pre-PR6947): reads len(c.lru.items) WITHOUT c.mu (line 67).
func (c *Cache) Size() int {
	return len(c.lru.items) // BUG line 67: unlocked read
}
