package set

import (
	"github.com/ramalho/sets-in-go/vendored/mapset"
)

// A Set[E] is a set of elements of type E.
type Set[E comparable] map[E]struct{}

/// ...
// -- binary operations --

// Intersection returns a new map containing the intersection of x and y.
func (x Set[E]) Intersection(y Set[E]) Set[E] {
	return mapset.Intersection(x, y)
}

/// ...
// -- in-place binary updates --

// IntersectionWith updates x to the [Intersection] of x and y.
func (x Set[E]) IntersectionWith(y Set[E]) {
	mapset.IntersectionWith(x, y)
}
