package algorithms

type Stack[T any] struct {
	arr []T
}

func NewStack[T any](cap int) Stack[T] {
	arr := make([]T, 0, cap)
	return Stack[T]{
		arr: arr,
	}
}

func (s Stack[T]) Len() int {
	return s.Len()
}

func (s Stack[T]) Pop() T {
	v := s.arr[s.Len()-1]
	s.arr = s.arr[:s.Len()-1]
	return v
}

func (s Stack[T]) Push(val T) {
	s.arr = append(s.arr, val)
}
