package query

import (
	"container/list"
	"strings"
	"sync"

	"github.com/itchyny/gojq"
)

type cachedProgram struct {
	code *gojq.Code
	err  error
}

type cacheItem struct {
	key   string
	value cachedProgram
}

type lruProgramCache struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
}

func newLRUProgramCache(capacity int) *lruProgramCache {
	if capacity < 1 {
		capacity = 1
	}

	return &lruProgramCache{
		capacity: capacity,
		items:    make(map[string]*list.Element, capacity),
		order:    list.New(),
	}
}

func (c *lruProgramCache) get(key string) (cachedProgram, bool) {
	elem, ok := c.items[key]
	if !ok {
		return cachedProgram{}, false
	}

	c.order.MoveToFront(elem)
	item := elem.Value.(cacheItem)
	return item.value, true
}

func (c *lruProgramCache) put(key string, value cachedProgram) {
	if elem, ok := c.items[key]; ok {
		elem.Value = cacheItem{key: key, value: value}
		c.order.MoveToFront(elem)
		return
	}

	elem := c.order.PushFront(cacheItem{key: key, value: value})
	c.items[key] = elem

	if len(c.items) > c.capacity {
		c.evictOldest()
	}
}

func (c *lruProgramCache) evictOldest() {
	back := c.order.Back()
	if back == nil {
		return
	}

	item := back.Value.(cacheItem)
	delete(c.items, item.key)
	c.order.Remove(back)
}

func (c *lruProgramCache) len() int {
	return len(c.items)
}

var (
	programCacheMu sync.Mutex
	programCache   = newLRUProgramCache(defaultProgramCacheCapacity)
)

const defaultProgramCacheCapacity = 256

func Execute(query string, data []any) ([]any, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return data, nil
	}

	code, err := getProgram(trimmed)
	if err != nil {
		return nil, err
	}

	result := make([]any, 0)
	for _, d := range data {
		iter := code.Run(d)
		for {
			v, ok := iter.Next()
			if !ok {
				// No more objects
				break
			}
			if err, ok := v.(error); ok {
				if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
					// clean stop
					break
				}
				// real error
				return nil, err
			}
			result = append(result, v)
		}
	}

	return result, nil
}

func getProgram(query string) (*gojq.Code, error) {
	programCacheMu.Lock()
	if cached, ok := programCache.get(query); ok {
		programCacheMu.Unlock()
		return cached.code, cached.err
	}
	programCacheMu.Unlock()

	q, err := gojq.Parse(query)
	if err != nil {
		programCacheMu.Lock()
		programCache.put(query, cachedProgram{err: err})
		programCacheMu.Unlock()
		return nil, err
	}

	code, err := gojq.Compile(q)
	programCacheMu.Lock()
	if cached, ok := programCache.get(query); ok {
		programCacheMu.Unlock()
		return cached.code, cached.err
	}
	programCache.put(query, cachedProgram{code: code, err: err})
	programCacheMu.Unlock()

	return code, err
}
