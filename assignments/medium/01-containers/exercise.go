package lab

type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(v T)       { panic("TODO") }
func (s *Stack[T]) Pop() (T, bool) { panic("TODO") }
func (s *Stack[T]) Len() int       { panic("TODO") }

type Queue[T any] struct {
	items []T
	head  int
}

func (q *Queue[T]) Push(v T)       { panic("TODO") }
func (q *Queue[T]) Pop() (T, bool) { panic("TODO") }
func (q *Queue[T]) Len() int       { panic("TODO") }
