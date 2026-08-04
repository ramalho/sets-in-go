package setops_test

// The same operations as the rest of this package, written against Alan
// Donovan's proposed container/set and container/mapset packages.
//
// NEITHER PACKAGE IS IN GO. Both are open CLs, vendored into ../vendored so
// these examples can run. See ../vendored/README.md for provenance.
//
// Compare each example here with its handwritten counterpart:
//
//	Example_containsAll     ->  membership_test.go
//	Example_algebra         ->  algebra_test.go
//	Example_union           ->  union_test.go
//	Example_equal           ->  equality_test.go
//	Example_iterators       ->  iterators_test.go
//	Example_dedup           ->  dedup_test.go

import (
	"fmt"
	"maps"
	"slices"

	"github.com/ramalho/sets-in-go/vendored/mapset"
	"github.com/ramalho/sets-in-go/vendored/set"
)

// Membership is a method rather than the two-value map lookup, so it composes
// into larger expressions. Subset (A ⊆ B) — the "ContainsAll" the stdlib does
// not provide — takes an iter.Seq, so it works against any sequence, not just
// another set.
func Example_containsAll() {
	var primes set.Set[int] = set.Of(2, 3, 5, 7, 11)

	fmt.Println(primes.Contains(7))
	fmt.Println(primes.ContainsAll(slices.Values([]int{2, 5, 11})))
	fmt.Println(primes.ContainsAll(slices.Values([]int{2, 4})))
	// Output:
	// true
	// true
	// false
}

// The algebra that algebra_test.go builds by hand from maps.Clone and
// maps.DeleteFunc. Note String: elements are sorted the way fmt prints map
// keys, so the output is deterministic and slide-friendly.
func Example_algebra() {
	var a, b set.Set[int] = set.Of(1, 2, 3), set.Of(3, 4, 5)

	fmt.Println("a                      =", a)
	fmt.Println("b                      =", b)
	fmt.Println("a.Union(b)             =", a.Union(b))
	fmt.Println("a.Intersection(b)      =", a.Intersection(b))
	fmt.Println("a.Difference(b)        =", a.Difference(b))
	fmt.Println("a.SymmetricDiff(b)     =", a.SymmetricDifference(b))
	fmt.Println("a.Intersects(b)        =", a.Intersects(b))
	// Output:
	// a                      = {1, 2, 3}
	// b                      = {3, 4, 5}
	// a.Union(b)             = {1, 2, 3, 4, 5}
	// a.Intersection(b)      = {3}
	// a.Difference(b)        = {1, 2}
	// a.SymmetricDiff(b)     = {1, 2, 4, 5}
	// a.Intersects(b)        = true
}

// Binary operations return a new set; the "With" variants update in place and
// report nothing. The mutators return whether the set actually changed, which
// is the piece hardest to get right by hand.
func Example_union() {
	var a, b set.Set[int] = set.Of(1, 2), set.Of(2, 3)

	u := a.Union(b)
	fmt.Println(a.Len(), b.Len(), u) // inputs untouched

	a.UnionWith(b) // in-place
	fmt.Println(a)

	fmt.Println(a.Insert(9), a.Insert(9)) // true, then false: no change
	fmt.Println(a.Delete(42))             // false: absent
	// Output:
	// 2 2 {1, 2, 3}
	// {1, 2, 3}
	// true false
	// false
}

// Set equality, without maps.Equal's requirement that values be comparable.
func Example_equal() {
	var a, b set.Set[string] = set.Of("x", "y"), set.Of("y", "x")

	fmt.Println(a.Equal(b))
	fmt.Println(a.Equal(set.Of("x")))
	// Output:
	// true
	// false
}

// Set is *transparently* a map[E]struct{}, which is the central design claim:
// no wrapper, no accessor, so a set still ranges like the map it is and
// interoperates with everything that already takes a map.
func Example_transparent() {
	var s set.Set[int] = set.Of(7, 42)

	total := 0
	for elem := range s { // ranges directly, like the map it is
		total += elem
	}
	fmt.Println(total)
	fmt.Println(len(s))                      // the built-in, not a method
	fmt.Println(slices.Sorted(maps.Keys(s))) // and the stdlib maps package

	var m map[int]struct{} = s // assignable in both directions
	fmt.Println(len(m))
	// Output:
	// 49
	// 2
	// [7 42]
	// 2
}

// Sets are built from and turned back into iterators, so they drop into the
// Go 1.23 pipelines that iterators_test.go composes by hand.
func Example_iterators() {
	words := []string{"the", "quick", "brown", "the", "fox"}

	s := set.Collect(slices.Values(words))
	fmt.Println(slices.Sorted(s.All()))

	s.DeleteFunc(func(w string) bool { return len(w) < 4 })
	fmt.Println(slices.Sorted(s.All()))
	// Output:
	// [brown fox quick the]
	// [brown quick]
}

// Deduplication is Collect — no sort-then-Compact, and no requirement that
// elements be ordered.
func Example_dedup() {
	s := set.Collect(slices.Values([]int{3, 1, 3, 1, 2}))
	fmt.Println(s.Len(), s)
	// Output:
	// 3 {1, 2, 3}
}

// The migration story. Existing Go code overwhelmingly spells a set
// map[K]bool; container/mapset operates on those directly, so legacy sets get
// the algebra without being rewritten. Operands may mix representations — the
// result takes the type of the first.
func Example_legacyMaps() {
	legacy := map[string]bool{"read": true, "write": true}
	modern := set.Of("write", "admin")

	fmt.Println(mapset.String(legacy))
	fmt.Println(mapset.String(mapset.Intersection(legacy, modern)))
	fmt.Println(mapset.String(mapset.Difference(legacy, modern)))

	granted := mapset.Union(legacy, modern)
	fmt.Printf("%T %v\n", granted, mapset.String(granted))
	// Output:
	// {read, write}
	// {write}
	// {read}
	// map[string]bool {admin, read, write}
}

// A rough edge in the current patchset, worth showing rather than hiding:
// set.Of is declared to return map[E]struct{} instead of Set[E], so its result
// carries no methods and every example above has to convert or pre-declare.
//
//	var a set.Set[int] = set.Of(1, 2, 3) // ok: assignable to the named type
//	set.Of(1, 2, 3).Intersection(b)      // compile error: no method Intersection
//
// set.Collect, on the next line of set.go, returns Set[E] correctly — so this
// reads as an oversight in code still under review, not a design decision.
