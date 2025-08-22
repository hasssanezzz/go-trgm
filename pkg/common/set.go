package common

type Set struct {
	values map[IndexEntry]struct{}
}

func NewSet(values ...IndexEntry) Set {
	s := Set{
		values: map[IndexEntry]struct{}{},
	}

	if len(values) > 0 {
		s.UnionSlice(values)
	}

	return s
}

func (s Set) Add(v IndexEntry) {
	s.values[v] = struct{}{}
}

func (s Set) Contains(v IndexEntry) bool {
	_, ok := s.values[v]
	return ok
}

func (s Set) Remove(v IndexEntry) {
	delete(s.values, v)
}

func (s Set) Union(a Set) {
	for value := range a.values {
		s.values[value] = struct{}{}
	}
}

func (s Set) UnionSlice(slice []IndexEntry) {
	for _, value := range slice {
		s.values[value] = struct{}{}
	}
}

func (s Set) ToSlice() []IndexEntry {
	result := make([]IndexEntry, 0, len(s.values))
	for value, _ := range s.values {
		result = append(result, value)
	}
	return result
}

func (s Set) Size() int {
	return len(s.values)
}
