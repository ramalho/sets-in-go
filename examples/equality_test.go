package setops_test

import (
	"fmt"
	"maps"
	"slices"

	setops "github.com/ramalho/modern-sets-examples"
)

// Set equality: A = B, order-independent.
//
// maps.Equal is the real "set equality" for map-sets: two sets are equal
// iff they have the same keys, regardless of insertion order.
func Example_equalSets() {
	a := setops.NewSet(3, 1, 2)
	b := setops.NewSet(1, 2, 3)
	fmt.Println(maps.Equal(a, b))
	// Output: true
}

// slices.Equal, by contrast, is sequence equality: it requires the same
// elements in the same order, so it is NOT set equality.
func Example_equalSlices() {
	a := []int{3, 1, 2}
	b := []int{1, 2, 3}
	fmt.Println(slices.Equal(a, b)) // different order

	// To compare slices AS sets, canonicalize both first.
	slices.Sort(a)
	slices.Sort(b)
	fmt.Println(slices.Equal(a, b))
	// Output:
	// false
	// true
}
