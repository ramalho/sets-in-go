package setops_test

import (
	"fmt"
	"maps"
	"slices"

	setops "github.com/ramalho/modern-sets-examples"
)

// Union: A ∪ B.
//
// maps.Copy(dst, src) is the one standard-library function that IS a set
// operation: it merges src's keys into dst. For map-sets that is union.
// (For general maps it is right-biased: src values win on key conflicts.)
func Example_copyUnion() {
	a := setops.NewSet(1, 2, 3)
	b := setops.NewSet(3, 4, 5)

	maps.Copy(a, b) // a becomes a ∪ b

	got := slices.Sorted(maps.Keys(a))
	fmt.Println(got)
	// Output: [1 2 3 4 5]
}

// If you must not mutate the inputs, clone first. maps.Clone + maps.Copy
// gives a non-destructive union.
func union[E comparable](a, b setops.Set[E]) setops.Set[E] {
	out := maps.Clone(a)
	maps.Copy(out, b)
	return out
}

func Example_copyNonDestructive() {
	a := setops.NewSet(1, 2)
	b := setops.NewSet(2, 3)

	u := union(a, b)

	fmt.Println(len(a), len(b)) // inputs untouched
	fmt.Println(slices.Sorted(maps.Keys(u)))
	// Output:
	// 2 2
	// [1 2 3]
}
