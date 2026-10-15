package lab

import "container/list"

type entry struct {
	key   string
	value int
}
type LRU struct {
	capacity int
	order    *list.List
	index    map[string]*list.Element
}

func NewLRU(capacity int) *LRU {
	return &LRU{capacity: capacity, order: list.New(), index: make(map[string]*list.Element)}
}
func (c *LRU) Get(key string) (int, bool) {
	e, ok := c.index[key]
	if !ok {
		return 0, false
	}
	c.order.MoveToFront(e)
	return e.Value.(entry).value, true
}
func (c *LRU) Put(key string, value int) {
	if c.capacity <= 0 {
		return
	}
	if e, ok := c.index[key]; ok {
		e.Value = entry{key, value}
		c.order.MoveToFront(e)
		return
	}
	c.index[key] = c.order.PushFront(entry{key, value})
	if len(c.index) > c.capacity {
		e := c.order.Back()
		delete(c.index, e.Value.(entry).key)
		c.order.Remove(e)
	}
}
func (c *LRU) Len() int { return len(c.index) }
