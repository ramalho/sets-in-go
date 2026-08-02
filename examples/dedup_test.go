package setops_test

import (
	"fmt"
	"slices"
)

// Deduplication: turning a slice (a multiset) into a canonical set.
//
// slices.Compact removes only ADJACENT runs of equal elements, so you
// almost always sort first. This is the idiomatic "slice → set" recipe.
func Example_compact() {
	nums := []int{3, 1, 2, 3, 1, 1, 2}

	slices.Sort(nums)           // [1 1 1 2 2 3 3]
	nums = slices.Compact(nums) // [1 2 3]

	fmt.Println(nums)
	// Output: [1 2 3]
}

// The gotcha, shown deliberately: Compact WITHOUT sorting leaves the
// non-adjacent duplicate 3 in place.
func Example_compactGotcha() {
	nums := []int{3, 1, 2, 3}
	fmt.Println(slices.Compact(nums)) // 3 and 3 are not adjacent
	// Output: [3 1 2 3]
}
