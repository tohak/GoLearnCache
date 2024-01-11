package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type item struct {
	value      interface{}
	expiration int64
}

func (i item) expired() bool {
	return time.Now().UnixNano() > i.expiration
}

type Cache struct {
	store     map[string]item
	storeTime int64
	mutex     sync.RWMutex
}

func NewCache(timeMinute int64) *Cache {
	if timeMinute == 0 {
		panic("error: time должен быть больше 1 минуты")
	}
	c := Cache{
		store:     make(map[string]item),
		storeTime: int64(time.Duration(timeMinute) * time.Minute),
	}
	go c.schedulerClean()
	return &c
}
func (c *Cache) AddCache(key string, value interface{}) {
	c.mutex.Lock()
	c.store[key] = item{
		value:      value,
		expiration: time.Now().Add(time.Duration(c.storeTime)).UnixNano(),
	}
	c.mutex.Unlock()
}
func (c *Cache) GetCacheItem(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, found := c.store[key]
	if value.expired() || !found {
		return nil, false
	}
	return value.value, true
}
func (c *Cache) DeleteCacheItem(key string) (bool, error) {
	_, found := c.store[key]
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if found {
		delete(c.store, key)
		return true, nil
	}
	err := errors.New("такого " + key + " не найдено в кеше")
	return false, err

}
func (c *Cache) Clean() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, item := range c.store {
		if item.expired() {
			delete(c.store, key)
		}
	}
}
func (c *Cache) schedulerClean() {
	for {
		c.Clean()
		fmt.Println("schedulerClean use")
		time.Sleep(time.Second)
	}
}
