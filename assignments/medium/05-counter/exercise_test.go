package lab

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	var c Counter
	if c.Value() != 0 {
		t.Fatal("zero value")
	}
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				c.Add(1)
				_ = c.Value()
			}
		}()
	}
	wg.Wait()
	if c.Value() != 32000 {
		t.Fatal(c.Value())
	}
	c.Add(-2)
	if c.Value() != 31998 {
		t.Fatal(c.Value())
	}
}
