package lab

type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	i := len(s.items) - 1
	v := s.items[i]
	s.items[i] = zero
	s.items = s.items[:i]
	return v, true
}
func (s *Stack[T]) Len() int { return len(s.items) }

type Queue[T any] struct {
	items []T
	head  int
}

func (q *Queue[T]) Push(v T) {
	if q.head > 0 && q.head*2 >= len(q.items) {
		copy(q.items, q.items[q.head:])
		n := len(q.items) - q.head
		clear(q.items[n:])
		q.items = q.items[:n]
		q.head = 0
	}
	q.items = append(q.items, v)
}
func (q *Queue[T]) Pop() (T, bool) {
	var zero T
	if q.head == len(q.items) {
		return zero, false
	}
	v := q.items[q.head]
	q.items[q.head] = zero
	q.head++
	if q.head == len(q.items) {
		q.items = nil
		q.head = 0
	}
	return v, true
}
func (q *Queue[T]) Len() int { return len(q.items) - q.head }
