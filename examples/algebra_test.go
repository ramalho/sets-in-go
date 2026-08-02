package setops_test

import (
	"cmp"
	"fmt"
	"maps"
	"slices"

	setops "github.com/ramalho/modern-sets-examples"
)

// The set algebra the stdlib does NOT ship: intersection, difference,
// and symmetric difference. Each is a short loop over map-sets. Showing
// them by hand is a good motivation for part 3 (a proper Set type).

// Intersection: A ∩ B. Iterate the smaller set for efficiency.
func intersection[E comparable](a, b setops.Set[E]) setops.Set[E] {
	if len(b) < len(a) {
		a, b = b, a
	}
	out := setops.Set[E]{}
	for e := range a {
		if _, ok := b[e]; ok {
			out[e] = struct{}{}
		}
	}
	return out
}

// Difference: A ∖ B (elements in A but not in B).
func difference[E comparable](a, b setops.Set[E]) setops.Set[E] {
	out := maps.Clone(a)
	maps.DeleteFunc(out, func(e E, _ struct{}) bool {
		_, inB := b[e]
		return inB
	})
	return out
}

// Symmetric difference: (A ∖ B) ∪ (B ∖ A).
func symmetricDifference[E comparable](a, b setops.Set[E]) setops.Set[E] {
	return union(difference(a, b), difference(b, a))
}

func sorted[E cmp.Ordered](s setops.Set[E]) []E {
	return slices.Sorted(maps.Keys(s))
}

func ExampleSet_algebra() {
	a := setops.NewSet(1, 2, 3, 4)
	b := setops.NewSet(3, 4, 5, 6)

	fmt.Println("∩ ", sorted(intersection(a, b)))
	fmt.Println("∖ ", sorted(difference(a, b)))
	fmt.Println("△ ", sorted(symmetricDifference(a, b)))
	// Output:
	// ∩  [3 4]
	// ∖  [1 2]
	// △  [1 2 5 6]
}
