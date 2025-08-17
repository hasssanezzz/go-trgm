package trgm

type set[T comparable] struct {
	values map[T]struct{}
}

func (s *set[T]) add(v T) {
	s.values[v] = struct{}{}
}

func (s *set[T]) contains(v T) bool {
	_, ok := s.values[v]
	return ok
}

func (s *set[T]) remove(v T) {
	delete(s.values, v)
}

func (s *set[T]) toSlice() []T {
	result := make([]T, 0, len(s.values))
	for value, _ := range s.values {
		result = append(result, value)
	}
	return result
}
