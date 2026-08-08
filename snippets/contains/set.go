package set

import (
	"iter"
	"maps"

	"github.com/ramalho/sets-in-go/vendored/mapset"
)

// A Set[E] is a set of elements of type E.
type Set[E comparable] map[E]struct{}

/// ...

// -- queries --  (reordered to suit this presentation)

// Contains reports whether set x contains element e.
func (x Set[E]) Contains(e E) bool {
	return mapset.Contains(x, e)
}

// ContainsAll reports whether set x contains all the elements of the sequence.
func (x Set[E]) ContainsAll(seq iter.Seq[E]) bool {
	return mapset.ContainsAll(x, seq)
}

// Equal reports whether sets x and y contain the same set of elements.
func (x Set[E]) Equal(y Set[E]) bool {
	return mapset.Equal(x, y)
}

// All returns an iterator over the sequence of elements of x.
// The sequence order is undefined and likely random.
func (x Set[E]) All() iter.Seq[E] {
	return maps.Keys(x)
}

// Len returns the number of set elements.
func (x Set[E]) Len() int {
	// (This method is provided so that the method
	// set forms a complete abstract data type.)
	return len(x)
}
