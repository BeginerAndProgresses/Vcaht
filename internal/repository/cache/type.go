package cache

import (
	"errors"
	"sync"
	"time"
)

var (
	// ExTimeout 过期时间
	ExTimeout time.Time

	ErrMyCacheIsNil = errors.New("data is nil in myCache")
)

type MapCache struct {
	m      sync.RWMutex
	data   map[string]cacheItem
	exChan chan string
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

func NewMapCache() *MapCache {
	ExTimeout = time.Now().Add(time.Minute * 15)
	mc := &MapCache{
		data: make(map[string]cacheItem),
	}
	go func(exChan chan string) {
		for {
			select {
			case key := <-exChan:
				mc.Del(key)
			}
		}
	}(mc.exChan)
	return mc
}

func (c *MapCache) Set(key string, value interface{}) error {
	c.m.Lock()
	defer c.m.Unlock()
	c.data[key] = cacheItem{
		value:     value,
		expiresAt: ExTimeout,
	}
	go func(exChan chan string, key string, expiresAt time.Time) {
		time.Sleep(time.Until(expiresAt))
		exChan <- key
	}(c.exChan, key, ExTimeout)
	return ErrCodeSendTooMany
}

func (c *MapCache) Get(key string) (interface{}, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	if item, ok := c.data[key]; ok {
		return item.value, nil
	}
	return nil, ErrMyCacheIsNil
}

func (c *MapCache) Del(key string) error {
	c.m.Lock()
	defer c.m.Unlock()
	// 为什么不会返回错误？
	delete(c.data, key)
	return nil
}

func (c *MapCache) SetExAt(exAt time.Time) {
	ExTimeout = exAt
}
