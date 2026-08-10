# Abstract collection interfaces

Defined in `go/src/container/container_test.go`.

Some comments from that file:

```
// The following interfaces define the abstract data types for
// collections in Go. They are expressed using F-bounded polymorphic
// interfaces to achieve covariant parameter/result specialization,
// and may be used as constraint types in generic functions.
//
// These interfaces are not yet published, but may be included in a
// future Go release once we have experience of whether these methods
// are necessary and sufficient.
//
// ...
//
// Until such time as these interfaces are published, when defining
// your own generic functions, we suggest you make your own copy of
// these constraint types. For maximum generality, you may wish to
// delete any methods that are not required by a particular function.
```

The interfaces, with some lines reordered for presentation.

```go
// _AbstractCollection models a collection C of elements E,
// such as *hash.Map, *hash.Set, *ordered.Map, or set.Set.
type _AbstractCollection[E any, C _AbstractCollection[E, C]] interface {
	Clear()
	Clone() C
	Contains(E) bool
	ContainsAll(iter.Seq[E]) bool
	Len() int
	String() string
}

// _AbstractMap models a mapping M from keys K to values V,
// such as *hash.Map or *ordered.Map.
type _AbstractMap[K, V any, M _AbstractMap[K, V, M]] interface {
	_AbstractCollection[K, M]

	Set(K, V) (V, bool)
	SetAll(iter.Seq2[K, V]) bool
	Get(K) (V, bool)
	At(K) V
	All() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
	Delete(K) (V, bool)
	DeleteAll(iter.Seq[K]) bool
	DeleteFunc(func(K, V) bool) bool
}

// _AbstractSet models a set S of elements E,
// such as *hash.Set, or set.Set.
type _AbstractSet[E any, S _AbstractSet[E, S]] interface {
	_AbstractCollection[E, S]

	Insert(E) bool
	InsertAll(iter.Seq[E]) bool
	Equal(S) bool
	All() iter.Seq[E]
	Intersection(S) S
	IntersectionWith(S)
	Intersects(S) bool
	Union(S) S
	UnionWith(S)
	Difference(S) S
	DifferenceWith(S)
	SymmetricDifference(S) S
	SymmetricDifferenceWith(S)
	Delete(E) bool
	DeleteAll(iter.Seq[E]) bool
	DeleteFunc(func(E) bool) bool
}

// NOTES
//
// These types define the fundamental operations that need to be
// implemented by all set and map types. Their naming, signature, and
// semantic conventions should be followed wherever possible when
// defining new collection types.
//
// AbstractMap cannot have Equal since V values may be non-comparable.
//
// The iteration order of All is not specified by the interface. For
// some collections (e.g. ordered.Map) it is ascending key order; for
// others (e.g. set.Set) it is random; and for still others it may be
// unspecified.
//
// Map.Set should replace an existing entry with an equivalent key.
// This follows the built-in map: https://go.dev/play/p/pkH8kkFTuEg.
//
// The "plural" functions {Contains,Delete,Insert,Set}All are
// sufficiently important that they belong as methods; they
// also compute a convenient bool result.


```