package cache

import (
	"time"

	"github.com/bluele/gcache"
)

type Cache interface {
	Set(key, value interface{}) error
	SetWithExpire(key, value interface{}, expiration time.Duration) error
	Get(key interface{}) (interface{}, error)
	Remove(key interface{}) bool
}

type LRU struct {
	cache   gcache.Cache
	enabled bool
}

func NewCache(capacity int) Cache {
	return &LRU{
		cache: gcache.New(capacity).LRU().Build(),
	}
}

func (l *LRU) Set(key, value interface{}) error {
	if !l.enabled {
		return nil
	}
	return l.cache.Set(key, value)
}

func (l *LRU) SetWithExpire(key, value interface{}, expiration time.Duration) error {
	if !l.enabled {
		return nil
	}
	return l.cache.SetWithExpire(key, value, expiration)
}

func (l LRU) Get(key interface{}) (interface{}, error) {
	return l.cache.Get(key)
}

func (l *LRU) Remove(key interface{}) bool {
	if !l.enabled {
		return false
	}
	return l.cache.Remove(key)
}

func (l *LRU) Expire(key interface{}) bool {
	return l.Remove(key)
}

func (l *LRU) Enable() {
	l.enabled = true
}

func (l *LRU) Disable() {
	l.enabled = false
}
