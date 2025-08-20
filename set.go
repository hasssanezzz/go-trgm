package trgm

type set struct {
	values map[IndexEntry]struct{}
}

func newSet(values ...IndexEntry) set {
	s := set{
		values: map[IndexEntry]struct{}{},
	}

	if len(values) > 0 {
		s.unionSlice(values)
	}

	return s
}

func (s set) add(v IndexEntry) {
	s.values[v] = struct{}{}
}

func (s set) contains(v IndexEntry) bool {
	_, ok := s.values[v]
	return ok
}

func (s set) remove(v IndexEntry) {
	delete(s.values, v)
}

func (s set) union(a set) {
	for value := range a.values {
		s.values[value] = struct{}{}
	}
}

func (s set) unionSlice(slice []IndexEntry) {
	for _, value := range slice {
		s.values[value] = struct{}{}
	}
}

func (s set) toSlice() []IndexEntry {
	result := make([]IndexEntry, 0, len(s.values))
	for value, _ := range s.values {
		result = append(result, value)
	}
	return result
}

func (s set) size() int {
	return len(s.values)
}
