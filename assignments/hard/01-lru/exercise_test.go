package lab

import "testing"

func TestLRU(t *testing.T) {
	c := NewLRU(2)
	c.Put("a", 1)
	c.Put("b", 2)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatal(v, ok)
	}
	c.Put("c", 3)
	if _, ok := c.Get("b"); ok {
		t.Fatal("wrong eviction")
	}
	c.Put("a", 9)
	c.Put("d", 4)
	if _, ok := c.Get("c"); ok {
		t.Fatal("update did not refresh")
	}
	if v, ok := c.Get("a"); !ok || v != 9 || c.Len() != 2 {
		t.Fatal(v, ok, c.Len())
	}
	for _, cap := range []int{0, -2} {
		c = NewLRU(cap)
		c.Put("x", 1)
		if _, ok := c.Get("x"); ok || c.Len() != 0 {
			t.Fatal("disabled cache")
		}
	}
	c = NewLRU(1)
	c.Put("x", 0)
	c.Put("x", 2)
	if c.Len() != 1 {
		t.Fatal("update duplicated entry")
	}
}
