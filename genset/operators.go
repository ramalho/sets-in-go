package genset

// Intersection returns a new set with elements that are in s AND in other.
// Math: S ∩ Z.
func (s Set[E]) Intersection(other Set[E]) Set[E] {
	var longer, shorter Set[E]
	if s.Len() > other.Len() {
		longer = s
		shorter = other
	} else {
		longer = other
		shorter = s
	}
	res := Make[E]()
	for elem := range shorter.store {
		if longer.Contains(elem) {
			res.store[elem] = struct{}{}
		}
	}
	return res
}

// Union returns a new Set: with elements that are in s OR in other.
// Math: S ∪ Z.
func (s Set[E]) Union(other Set[E]) Set[E] {
	res := s.Copy()
	for elem := range other.store {
		res.store[elem] = struct{}{}
	}
	return res
}

// Difference returns a new Set with the elements of s minus the elements of other.
// Math: S \ Z.
func (s Set[E]) Difference(other Set[E]) Set[E] {
	res := Make[E]()
	for elem := range s.store {
		if !other.Contains(elem) {
			res.store[elem] = struct{}{}
		}
	}
	return res
}

// SymmetricDifference returns a new Set with elements present
// in either set but not on both. Think boolean XOR.
// Math: S ∆ Z.
func (s Set[E]) SymmetricDifference(other Set[E]) Set[E] {
	all := s.Union(other)
	common := s.Intersection(other)
	return all.Difference(common)
}
