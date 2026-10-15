package lab

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestTTL(t *testing.T) {
	now := time.Unix(0, 0)
	c := NewCache(func() time.Time { return now })
	c.Set("a", "one", time.Second)
	c.Set("b", "two", 2*time.Second)
	if v, ok := c.Get("a"); !ok || v != "one" {
		t.Fatal(v, ok)
	}
	now = now.Add(time.Second)
	if _, ok := c.Get("a"); ok {
		t.Fatal("exact expiry")
	}
	if c.Prune() != 0 {
		t.Fatal("unexpired entry removed")
	}
	now = now.Add(time.Second)
	if c.Prune() != 1 || c.Prune() != 0 {
		t.Fatal("prune count")
	}
	c.Set("z", "old", time.Hour)
	c.Set("z", "new", 0)
	if _, ok := c.Get("z"); ok {
		t.Fatal("zero TTL")
	}
	c.Set("z", "old", time.Hour)
	c.Set("z", "new", -1)
	if _, ok := c.Get("z"); ok {
		t.Fatal("negative TTL")
	}
	c.Set("r", "old", time.Second)
	now = now.Add(time.Second / 2)
	c.Set("r", "new", time.Hour)
	now = now.Add(time.Second)
	if v, ok := c.Get("r"); !ok || v != "new" {
		t.Fatal("replacement TTL")
	}
}
func TestTTLConcurrent(t *testing.T) {
	c := NewCache(nil)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := strconv.Itoa(i)
			for n := 0; n < 100; n++ {
				c.Set(k, "v", time.Hour)
				c.Get(k)
				c.Prune()
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < 20; i++ {
		if _, ok := c.Get(strconv.Itoa(i)); !ok {
			t.Fatal("lost item")
		}
	}
}
