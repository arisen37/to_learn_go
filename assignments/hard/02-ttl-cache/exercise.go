package lab

import "time"

type Cache struct{}

func NewCache(now func() time.Time) *Cache                { panic("TODO") }
func (c *Cache) Set(key, value string, ttl time.Duration) { panic("TODO") }
func (c *Cache) Get(key string) (string, bool)            { panic("TODO") }
func (c *Cache) Prune() int                               { panic("TODO") }
