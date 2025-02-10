package generics

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(value T) {
	s.items = append(s.items, value)
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	index := len(s.items) - 1
	value := s.items[index]
	s.items = s.items[:index]
	return value, true
}

func (s *Stack[T]) IsEmpty() bool {
	if len(s.items) == 0 {
		return true
	} else {
		return false
	}
}
