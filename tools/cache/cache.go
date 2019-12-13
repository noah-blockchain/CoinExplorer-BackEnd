package cache

import (
	"sync"
	"time"
)

type ExplorerCache struct {
	items       *sync.Map
}

// cache constructor
func NewCache() *ExplorerCache {
	cache := &ExplorerCache{
		items:       new(sync.Map),
	}

	return cache
}

// create new cache item
func (c *ExplorerCache) newCacheItem(value interface{}, ttl time.Duration) *Item {
	end := time.Now().Add(ttl * time.Second)
	return &Item{value: value, ttl: &end}
}

// get or store value from cache
func (c *ExplorerCache) Get(key interface{}, callback func() interface{}, ttl time.Duration) interface{} {
	v, ok := c.items.Load(key)
	if ok {
		item := v.(*Item)
		if !item.IsExpired() {
			return item.value
		}
	}

	return c.Store(key, callback(), ttl)
}

// save value to cache
func (c *ExplorerCache) Store(key interface{}, value interface{}, ttl time.Duration) interface{} {
	c.items.Store(key, c.newCacheItem(value, ttl))
	return value
}

// loop for checking items expiration
func (c *ExplorerCache) ExpirationCheck() {
	c.items.Range(func(key, value interface{}) bool {
		item := value.(*Item)
		if item.IsExpired() {
			c.items.Delete(key)
		}

		return true
	})
}

// set new last block id
//func (c *ExplorerCache) SetBlockId(id uint64) {
//	c.lastBlockId = id
//	// clean expired items
//	go c.ExpirationCheck()
//}
//
//// update last block id by ws data
//func (c *ExplorerCache) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
//	var block blocks.Resource
//	err := json.Unmarshal(e.Data, &block)
//	helpers.CheckErr(err)
//
//	// update last block id
//	c.SetBlockId(block.ID)
//}

