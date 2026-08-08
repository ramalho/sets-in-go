package set

import (
    "iter"

    "github.com/ramalho/sets-in-go/vendored/mapset"
)

// A Set[E] is a set of elements of type E.
type Set[E comparable] map[E]struct{}

// Collect creates a new set containing the elements of the sequence.
func Collect[E comparable](seq iter.Seq[E]) Set[E] {
    return Set[E](mapset.Collect(seq))
}
