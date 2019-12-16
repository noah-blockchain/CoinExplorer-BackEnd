package cache

import (
	"time"
)

type Item struct {
	value interface{} // cached value
	ttl   *time.Time  // expiration time
}

func (c *Item) IsExpired() bool {
	if c.ttl != nil && time.Now().Before(*c.ttl) {
		return false
	}
	return true
}
