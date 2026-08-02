package setops_test

import (
	"fmt"
	"slices"

	setops "github.com/ramalho/modern-sets-examples"
)

// Membership: x ∈ A
//
// For a slice, slices.Contains does a linear scan.
func Example_contains() {
	primes := []int{2, 3, 5, 7, 11}
	fmt.Println(slices.Contains(primes, 7))
	fmt.Println(slices.Contains(primes, 8))
	// Output:
	// true
	// false
}

// For a map-set, membership is a direct key lookup — O(1) and the
// most idiomatic form.
func ExampleSet_membership() {
	primes := setops.NewSet(2, 3, 5, 7, 11)
	_, ok := primes[7]
	fmt.Println(ok)
	_, ok = primes[8]
	fmt.Println(ok)
	// Output:
	// true
	// false
}

// Subset / superset (A ⊆ B) is exactly the "ContainsAll" the stdlib does
// NOT provide. With slices it is a scan-per-element (O(len(a)*len(b))).
func containsAllSlice[E comparable](super, sub []E) bool {
	for _, e := range sub {
		if !slices.Contains(super, e) {
			return false
		}
	}
	return true
}

// With map-sets the same test is O(len(sub)) key lookups — this is why
// sets are the right tool for the job.
func isSubset[E comparable](sub, super setops.Set[E]) bool {
	if len(sub) > len(super) {
		return false
	}
	for e := range sub {
		if _, ok := super[e]; !ok {
			return false
		}
	}
	return true
}

func ExampleSet_subset() {
	letters := setops.NewSet('a', 'b', 'c', 'd')
	vowels := setops.NewSet('a')

	fmt.Println(isSubset(vowels, letters)) // vowels ⊆ letters
	fmt.Println(isSubset(letters, vowels)) // letters ⊆ vowels ?
	fmt.Println(containsAllSlice([]int{1, 2, 3}, []int{1, 3}))
	// Output:
	// true
	// false
	// true
}
