package setops_test

// Superset and subset under the proposed container/set and container/mapset.
//
// Neither package has a Superset or Subset predicate. ContainsAll *is* the
// superset test — it just takes an iter.Seq[E] rather than a set, so the
// other set's iterator is what you feed it. See ../references/API-table.md.
//
// Compare isSubset in membership_test.go, the handwritten counterpart.

import (
	"fmt"
	"maps"

	"github.com/ramalho/sets-in-go/vendored/mapset"
	"github.com/ramalho/sets-in-go/vendored/set"
)

// supersetOf reports whether x ⊇ y.
//
// The length guard is an optimization, not correctness: ContainsAll alone is
// already right, since it walks y and fails on the first element x lacks.
// The guard settles mismatched sizes in O(1) instead.
//
// Note the asymmetry with Intersects and Intersection, which iterate whichever
// operand is smaller: a superset test is directional, so it must iterate y and
// probe x. O(len(y)) lookups is the floor.
func supersetOf[E comparable](x, y set.Set[E]) bool {
	return len(x) >= len(y) && x.ContainsAll(y.All())
}

func subsetOf[E comparable](x, y set.Set[E]) bool { return supersetOf(y, x) }

// x ⊃ y. Once x ⊇ y holds, len(x) > len(y) is exactly "x has an element y
// lacks", so no second pass over the elements is needed.
func properSupersetOf[E comparable](x, y set.Set[E]) bool {
	return supersetOf(x, y) && len(x) > len(y)
}

func Example_superset() {
	var letters, vowels set.Set[rune] = set.Of('a', 'b', 'c', 'e'), set.Of('a', 'e')

	fmt.Println(supersetOf(letters, vowels))       // letters ⊇ vowels
	fmt.Println(supersetOf(vowels, letters))       // vowels ⊇ letters?
	fmt.Println(subsetOf(vowels, letters))         // vowels ⊆ letters
	fmt.Println(supersetOf(letters, letters))      // ⊇ is reflexive
	fmt.Println(properSupersetOf(letters, vowels)) // but ⊃ is not
	fmt.Println(properSupersetOf(letters, letters))
	// Output:
	// true
	// false
	// true
	// true
	// true
	// false
}

// The nil set is the empty set here, and needs no special-casing: ranging over
// a nil map yields nothing, and lookups in one miss.
func Example_supersetNil() {
	var some set.Set[int] = set.Of(1, 2)
	var none set.Set[int] // nil

	fmt.Println(supersetOf(some, none))
	fmt.Println(supersetOf(none, some))
	fmt.Println(supersetOf(none, none))
	// Output:
	// true
	// false
	// true
}

// mapset.ContainsAll takes the sequence as a separate parameter, so the two
// operands may carry *different* value types — a legacy map[K]bool tested
// against a modern map[K]struct{}. Equal and IntersectionWith cannot do this:
// their signatures bind both operands to a single map type M.
func mapsetSupersetOf[MX ~map[K]VX, MY ~map[K]VY, K comparable, VX, VY bool | struct{}](x MX, y MY) bool {
	return len(x) >= len(y) && mapset.ContainsAll(x, maps.Keys(y))
}

func Example_supersetLegacyMaps() {
	granted := map[string]bool{"read": true, "write": true, "admin": true}
	required := mapset.Of("read", "write")

	fmt.Println(mapsetSupersetOf(granted, required))
	fmt.Println(mapsetSupersetOf(required, granted))
	// Output:
	// true
	// false
}

// One rough edge worth knowing: unlike Equal, Intersects, and Union,
// ContainsAll carries no maps.Identical fast path, and the length guard above
// does not catch the self-comparison either. x.ContainsAll(x.All()) really does
// walk the whole map. Check identity yourself if that lands on a hot path.
//
// (The other rough edge — set.Of returning map[E]struct{} instead of Set[E],
// which is why the examples here pre-declare their variables — is documented at
// the end of proposed_test.go.)
