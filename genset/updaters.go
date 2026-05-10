package genset

/* Implementation note: The only methods that change a Set
   after it is created are in this file.
*/

// Add element to set.
func (s Set[E]) Add(elem E) {
	s.store[elem] = struct{}{}
}

// AddAll adds elements to set.
func (s Set[E]) AddAll(elems ...E) {
	for _, elem := range elems {
		s.store[elem] = struct{}{}
	}
}

// Remove element from set, if it is present.
func (s Set[E]) Remove(elem E) {
	delete(s.store, elem)
}

// RemoveAll removes elements from set, if they are present.
func (s Set[E]) RemoveAll(elems ...E) {
	for _, elem := range elems {
		delete(s.store, elem)
	}
}

// Pop tries to return some element of s, deleting it. If there was an element,
// the pair (element, true) is returned. Otherwise, the result is (zero, false).
func (s Set[E]) Pop() (elem E, found bool) {
	for elem = range s.store {
		delete(s.store, elem)
		return elem, true
	}
	return
}

// Clear removes all elements from s.
func (s *Set[E]) Clear() {
	s.store = make(map[E]struct{})
}

/* Set operations that change the receiver */

// IntersectionUpdate changes s in-place, keeping only elements
// that are in s AND in other. Math: S ∩ Z.
func (s Set[E]) IntersectionUpdate(other Set[E]) {
	for elem := range s.store {
		if !other.Contains(elem) {
			delete(s.store, elem)
		}
	}
}

// UnionUpdate changes s in-place, gathering all elements that
// are in s OR in other. Math: S ∪ Z.
func (s Set[E]) UnionUpdate(other Set[E]) {
	for elem := range other.store {
		s.store[elem] = struct{}{}
	}
}

// DifferenceUpdate changes s in-place, removing all elements
// that appear in other. Math: S \ Z.
func (s Set[E]) DifferenceUpdate(other Set[E]) {
	for elem := range other.store {
		delete(s.store, elem)
	}
}

// SymmetricDifferenceUpdate changes s in-place, gathering only
// elements that are in either set but not on both.
// Think boolean XOR. Math: S ∆ Z.
func (s Set[E]) SymmetricDifferenceUpdate(other Set[E]) {
	common := s.Intersection(other)
	s.UnionUpdate(other)
	s.DifferenceUpdate(common)
}
