// Package setops collects small, runnable examples showing how set
// operations from mathematics map onto the modern Go standard library
// (the slices and maps packages from Go 1.21, and the iterator
// additions from Go 1.23).
//
// The standard library has no Set type and no set-algebra functions.
// The idiom is to represent a set as map[T]struct{} (a "map-set"), and
// to build the operations from primitives. Each file here focuses on
// one theme:
//
//   - membership_test.go — element membership and the missing ContainsAll
//   - dedup_test.go       — turning a slice (multiset) into a set
//   - equality_test.go    — set equality
//   - union_test.go       — union via maps.Copy
//   - algebra_test.go     — intersection / difference / symmetric difference (DIY)
//   - iterators_test.go   — composing sets with Go 1.23 iterators
//
// Run everything with:
//
//	go test ./...
package setops

// Set is the idiomatic Go set: a map whose keys are the elements and
// whose values carry no information (struct{} uses zero bytes).
type Set[E comparable] map[E]struct{}

// NewSet builds a Set from the given elements.
func NewSet[E comparable](elems ...E) Set[E] {
	s := make(Set[E], len(elems))
	for _, e := range elems {
		s[e] = struct{}{}
	}
	return s
}
