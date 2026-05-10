package genset

// SubsetOf reports whether s is a subset of other
// Math: S ⊆ Z.
func (s Set[E]) SubsetOf(other Set[E]) bool {
	if s.Len() > other.Len() {
		return false
	}
	return s.allIn(other)
}

// SupersetOf reports whether s is a superset of other
// Math: S ⊇ Z.
func (s Set[E]) SupersetOf(other Set[E]) bool {
	return other.SubsetOf(s)
}
