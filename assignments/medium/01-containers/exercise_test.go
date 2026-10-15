package lab

import "testing"

func TestContainers(t *testing.T) {
	var s Stack[string]
	if v, ok := s.Pop(); ok || v != "" {
		t.Fatal("empty stack")
	}
	s.Push("a")
	s.Push("b")
	if s.Len() != 2 {
		t.Fatal("stack len")
	}
	for _, want := range []string{"b", "a"} {
		v, ok := s.Pop()
		if !ok || v != want {
			t.Fatal(v, ok)
		}
	}
	var q Queue[int]
	if v, ok := q.Pop(); ok || v != 0 {
		t.Fatal("empty queue")
	}
	for i := 0; i < 1000; i++ {
		q.Push(i)
	}
	for i := 0; i < 500; i++ {
		v, ok := q.Pop()
		if !ok || v != i {
			t.Fatal(v, ok)
		}
	}
	for i := 1000; i < 1500; i++ {
		q.Push(i)
	}
	if q.Len() != 1000 {
		t.Fatal(q.Len())
	}
	for i := 500; i < 1500; i++ {
		v, ok := q.Pop()
		if !ok || v != i {
			t.Fatal(v, ok)
		}
	}
	if q.Len() != 0 {
		t.Fatal("queue len")
	}
	q.Push(7)
	if v, ok := q.Pop(); !ok || v != 7 {
		t.Fatal("reuse")
	}
}
