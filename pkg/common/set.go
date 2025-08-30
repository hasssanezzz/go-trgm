package common

type Set struct {
	values map[DocumentID]struct{}
}

func NewSet(values ...DocumentID) Set {
	s := Set{
		values: map[DocumentID]struct{}{},
	}

	if len(values) > 0 {
		s.UnionSlice(values)
	}

	return s
}

func (s Set) Add(v DocumentID) {
	s.values[v] = struct{}{}
}

func (s Set) Contains(v DocumentID) bool {
	_, ok := s.values[v]
	return ok
}

func (s Set) Remove(v DocumentID) {
	delete(s.values, v)
}

func (s Set) Union(a Set) {
	for value := range a.values {
		s.values[value] = struct{}{}
	}
}

func (s Set) UnionSlice(slice []DocumentID) {
	for _, value := range slice {
		s.values[value] = struct{}{}
	}
}

func (s Set) ToSlice() []DocumentID {
	result := make([]DocumentID, 0, len(s.values))
	for value, _ := range s.values {
		result = append(result, value)
	}
	return result
}

func (s Set) Size() int {
	return len(s.values)
}
