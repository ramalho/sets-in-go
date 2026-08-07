// Package setgm reimplements the proposed [container/set.Set] using the
// method-level type parameters added in Go 1.27.
//
// The proposed set.Set wraps container/mapset, but loses generality doing so.
// mapset's binary operations are heterogeneous by design — mapset.Union has
// signature
//
//	Union[MX ~map[K]VX, MY ~map[K]VY, K comparable, VX, VY bool | struct{}](x MX, y MY) MX
//
// so a modern map[K]struct{} may be unioned with a legacy map[K]bool. Every
// set.Set method throws that away:
//
//	func (x Set[E]) Union(y Set[E]) Set[E]
//
// Not because the design wants it, but because before Go 1.27 a method could
// not introduce type parameters of its own. The second operand had to be
// spelled in terms of the receiver's E alone, so it could only ever be
// another Set[E].
//
// SetGM restores the generality. Its binary methods take any map-shaped set
// — Set[E], SetGM[E], map[E]struct{}, or a legacy map[E]bool:
//
//	perms := setgm.Of("read", "write")
//	granted := map[string]bool{"read": true, "admin": true} // legacy set
//
//	perms.Union(granted)     // {admin, read, write}
//	perms.Superset(granted)  // false
//	perms.Intersects(granted) // true
//
// # Differences from the proposal
//
// Beyond the generic methods, SetGM makes three deliberate changes:
//
//   - [Of] returns SetGM[E]. The proposal's set.Of is declared as returning
//     map[E]struct{} even though its body returns Set[E], which is why the
//     examples in ../examples pre-declare their variables.
//   - [SetGM.Superset], [SetGM.Subset] and their proper variants are provided
//     directly. The proposal has no such predicate; ContainsAll is the superset
//     test, but it takes an iter.Seq[E], so callers must reach for the other
//     set's iterator. These also carry the identity fast path that ContainsAll
//     lacks.
//   - [SetGM.Equal], [SetGM.Intersects] and the three narrow *With operations
//     are implemented here rather than delegated, because the corresponding
//     mapset functions bind both operands to a single map type M and so cannot
//     express the heterogeneous case.
//
// # Cost of generic methods
//
// An interface may not declare a method with type parameters. Every generic
// method below is therefore invisible to interface satisfaction: SetGM cannot
// be reached through an interface exposing Union, Equal, Superset, and so on.
// The non-generic methods are unaffected, so SetGM still satisfies
// [fmt.Stringer]. A generic method also cannot be used as a method value
// without being instantiated first. See setgm_test.go for both.
package setgm

import (
	"iter"
	"maps"
	"unsafe"

	"github.com/ramalho/sets-in-go/vendored/mapset"
)

// A SetGM[E] is a set of elements of type E, represented transparently as a
// map[E]struct{} so that it may be used as the operand of a range loop.
type SetGM[E comparable] map[E]struct{}

// MapSet constrains the operand of SetGM's generic methods to any map-shaped
// set of E: map[E]struct{} and map[E]bool, and any named type whose underlying
// type is one of those, including set.Set[E] and SetGM[E] itself.
//
// As in mapset, a bool-valued operand is read as a set of keys; the boolean
// values are ignored, and sets built by this package always store true.
type MapSet[E comparable, V bool | struct{}] interface {
	~map[E]V
}

// identical reports whether two maps refer to the same data structure.
//
// This duplicates maps.Identical from the still-open CL 760800
// (golang/go#78456), which mapset depends on. The vendored module supplies it
// as ../vendored/internal/maps, but an internal package cannot be imported
// from another module, so the three lines are repeated here.
//
// Beware that shortcuts based on identical(x, y) may behave surprisingly for
// maps containing floating-point NaNs, since NaN != NaN.
func identical[MX ~map[E]VX, MY ~map[E]VY, E comparable, VX, VY any](x MX, y MY) bool {
	// Maps in Go are references yet the core language
	// provides no safe way to ask whether they alias.
	type pointer = unsafe.Pointer
	return *(*pointer)(pointer(&x)) == *(*pointer)(pointer(&y))
}

// -- construction --

// Of creates a new set containing the given elements.
func Of[E comparable](elems ...E) SetGM[E] {
	return SetGM[E](mapset.Of(elems...))
}

// Collect returns a new set containing the elements of the sequence.
func Collect[E comparable](seq iter.Seq[E]) SetGM[E] {
	return SetGM[E](mapset.Collect(seq))
}

// From returns a new set containing the keys of any map-shaped set, copying
// them. The map values are ignored.
//
// This is the conversion a plain type conversion cannot express: SetGM[E](m)
// is legal only when m is already a map[E]struct{}.
func From[M ~map[E]V, E comparable, V bool | struct{}](m M) SetGM[E] {
	return SetGM[E](mapset.Collect(mapset.All(m)))
}

// -- queries --

// Len returns the number of set elements.
func (x SetGM[E]) Len() int { return len(x) }

// Contains reports whether the set contains element e.
func (x SetGM[E]) Contains(e E) bool { return mapset.Contains(x, e) }

// All returns an iterator over the elements of x.
// The sequence order is undefined and likely random.
func (x SetGM[E]) All() iter.Seq[E] { return maps.Keys(x) }

// ContainsAll reports whether x contains every element of the sequence.
//
// This retains the proposal's signature. To test against another set, prefer
// [SetGM.Superset], which adds an identity check and a length guard.
func (x SetGM[E]) ContainsAll(seq iter.Seq[E]) bool { return mapset.ContainsAll(x, seq) }

// String returns a representation of the set as a string of the form
// "{a, ..., z}". Elements are sorted in the same manner as [fmt.Sprint] of a
// map, and each is printed as if by fmt.Sprint.
func (x SetGM[E]) String() string { return mapset.String(x) }

// Equal reports whether x and y contain the same elements.
// The values of y are ignored.
func (x SetGM[E]) Equal[M MapSet[E, V], V bool | struct{}](y M) bool {
	if identical(x, y) {
		return true
	}
	if len(x) != len(y) {
		return false
	}
	for k := range y {
		if _, ok := x[k]; !ok {
			return false
		}
	}
	return true
}

// Superset reports whether x ⊇ y: every element of y is in x.
// The values of y are ignored.
func (x SetGM[E]) Superset[M MapSet[E, V], V bool | struct{}](y M) bool {
	if identical(x, y) {
		return true // ⊇ is reflexive
	}
	// The length guard is an optimization, not correctness: the loop below is
	// already right, since it walks y and fails on the first element x lacks.
	if len(x) < len(y) {
		return false
	}
	for k := range y {
		if _, ok := x[k]; !ok {
			return false
		}
	}
	return true
}

// ProperSuperset reports whether x ⊃ y: x ⊇ y and x has an element y lacks.
func (x SetGM[E]) ProperSuperset[M MapSet[E, V], V bool | struct{}](y M) bool {
	// Once x ⊇ y holds, len(x) > len(y) is exactly "x has an element y lacks",
	// so no second pass over the elements is needed.
	return len(x) > len(y) && x.Superset(y)
}

// Subset reports whether x ⊆ y: every element of x is in y.
// The values of y are ignored.
func (x SetGM[E]) Subset[M MapSet[E, V], V bool | struct{}](y M) bool {
	if identical(x, y) {
		return true // ⊆ is reflexive
	}
	if len(x) > len(y) {
		return false
	}
	// A subset test is directional, so it must iterate x and probe y;
	// len(x) lookups is the floor.
	for k := range x {
		if _, ok := y[k]; !ok {
			return false
		}
	}
	return true
}

// ProperSubset reports whether x ⊂ y: x ⊆ y and y has an element x lacks.
func (x SetGM[E]) ProperSubset[M MapSet[E, V], V bool | struct{}](y M) bool {
	return len(x) < len(y) && x.Subset(y)
}

// Intersects reports whether x.Intersection(y) is non-empty.
// The values of y are ignored.
func (x SetGM[E]) Intersects[M MapSet[E, V], V bool | struct{}](y M) bool {
	if identical(x, y) {
		return len(x) > 0
	}
	// Iterate over the smaller of the two maps.
	if len(y) < len(x) {
		for k := range y {
			if _, ok := x[k]; ok {
				return true
			}
		}
	} else {
		for k := range x {
			if _, ok := y[k]; ok {
				return true
			}
		}
	}
	return false
}

// -- binary operations --

// Union returns a new set containing the elements of x and y.
func (x SetGM[E]) Union[M MapSet[E, V], V bool | struct{}](y M) SetGM[E] {
	return mapset.Union(x, y)
}

// Intersection returns a new set containing the elements present in both
// x and y.
func (x SetGM[E]) Intersection[M MapSet[E, V], V bool | struct{}](y M) SetGM[E] {
	return mapset.Intersection(x, y)
}

// Difference returns a new set containing the elements of x that are not
// present in y.
func (x SetGM[E]) Difference[M MapSet[E, V], V bool | struct{}](y M) SetGM[E] {
	return mapset.Difference(x, y)
}

// SymmetricDifference returns a new set containing the elements of x that are
// not present in y, and the elements of y that are not present in x.
func (x SetGM[E]) SymmetricDifference[M MapSet[E, V], V bool | struct{}](y M) SetGM[E] {
	return mapset.SymmetricDifference(x, y)
}

// -- mutations --

// Clear removes all elements from the set.
func (x SetGM[E]) Clear() { clear(x) }

// Clone returns a new non-nil set with the same elements as x.
func (x SetGM[E]) Clone() SetGM[E] {
	if x == nil {
		return make(SetGM[E])
	}
	return maps.Clone(x)
}

// Insert adds the specified element to the set.
// It reports whether len(x) changed.
func (x SetGM[E]) Insert(elem E) bool { return mapset.Insert(x, elem) }

// InsertAll adds each element of the addenda sequence to the set.
// It reports whether len(x) changed.
func (x SetGM[E]) InsertAll(addenda iter.Seq[E]) bool { return mapset.InsertAll(x, addenda) }

// Delete removes element e from the set.
// It reports whether len(x) changed.
func (x SetGM[E]) Delete(e E) bool { return mapset.Delete(x, e) }

// DeleteAll removes from the set each element of the delenda sequence.
// It reports whether len(x) changed.
func (x SetGM[E]) DeleteAll(delenda iter.Seq[E]) bool { return mapset.DeleteAll(x, delenda) }

// DeleteFunc removes each element e of the set such that f(e).
// It reports whether len(x) changed.
func (x SetGM[E]) DeleteFunc(f func(E) bool) bool { return mapset.DeleteFunc(x, f) }

// -- in-place binary updates --

// UnionWith updates x to the [SetGM.Union] of x and y.
func (x SetGM[E]) UnionWith[M MapSet[E, V], V bool | struct{}](y M) {
	mapset.UnionWith(x, y)
}

// IntersectionWith updates x to the [SetGM.Intersection] of x and y.
func (x SetGM[E]) IntersectionWith[M MapSet[E, V], V bool | struct{}](y M) {
	if identical(x, y) {
		return // x ∩ x = x
	}
	for k := range x {
		if _, ok := y[k]; !ok {
			delete(x, k)
		}
	}
}

// DifferenceWith updates x to the [SetGM.Difference] of x and y.
// In other words, it removes from x all the elements of y.
func (x SetGM[E]) DifferenceWith[M MapSet[E, V], V bool | struct{}](y M) {
	if identical(x, y) {
		clear(x) // x ∖ x = Ø
		return
	}
	for k := range y {
		delete(x, k)
	}
}

// SymmetricDifferenceWith updates x to the [SetGM.SymmetricDifference] of
// x and y.
func (x SetGM[E]) SymmetricDifferenceWith[M MapSet[E, V], V bool | struct{}](y M) {
	if identical(x, y) {
		clear(x) // x ∆ x = Ø
		return
	}
	for k := range y {
		if _, ok := x[k]; ok {
			delete(x, k)
		} else {
			x[k] = struct{}{}
		}
	}
}
