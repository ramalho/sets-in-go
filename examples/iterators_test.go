package setops_test

import (
	"fmt"
	"maps"
	"slices"

	setops "github.com/ramalho/modern-sets-examples"
)

// Go 1.23 iterators (iter.Seq) let you compose set-shaped pipelines.
// maps.Keys / maps.Values now return iterators, and slices.Sorted /
// slices.Collect consume them.

// The canonical "give me a set's elements in order" one-liner.
func Example_sortedKeys() {
	s := setops.NewSet("delta", "alpha", "charlie", "bravo")
	fmt.Println(slices.Sorted(maps.Keys(s)))
	// Output: [alpha bravo charlie delta]
}

// slices.Collect materializes any iterator into a slice — here, building
// a set straight from a filtered iterator without an intermediate slice.
func Example_collect() {
	src := setops.NewSet(1, 2, 3, 4, 5, 6)

	// keep even elements, as a sorted slice
	evens := slices.Sorted(func(yield func(int) bool) {
		for e := range src {
			if e%2 == 0 && !yield(e) {
				return
			}
		}
	})
	fmt.Println(evens)
	// Output: [2 4 6]
}

// Because maps.Keys is an iter.Seq, membership over the union of two sets
// can be expressed without building the union at all.
func Example_unionLazy() {
	a := setops.NewSet(1, 2, 3)
	b := setops.NewSet(3, 4, 5)

	all := slices.Sorted(chained(maps.Keys(a), maps.Keys(b)))
	fmt.Println(slices.Compact(all))
	// Output: [1 2 3 4 5]
}

// chained yields from each iterator in turn (like an "or" of sequences).
func chained[E any](seqs ...func(func(E) bool)) func(func(E) bool) {
	return func(yield func(E) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}
