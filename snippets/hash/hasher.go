package maphash

// A Hasher defines the interface between a hash-based container and its elements.
// It provides a hash function and an equivalence relation over values
// of type T, enabling those values to be inserted in hash tables
// and similar data structures.
//
// ...more than 100 lines of comments...

type Hasher[T any] interface {
    Hash(*Hash, T)
    Equal(x, y T) bool
}

// ComparableHasher is an implementation of [Hasher] whose
// Equal(x, y) method is consistent with x == y.
//
// ComparableHasher is defined only for comparable types.
// The type system will not prevent you from instantiating a type
// such as ComparableHasher[any]; nonetheless you must not pass
// non-comparable argument values to its Hash or Equal methods.
type ComparableHasher[T comparable] struct {
    _ [0]func(T) // disallow comparison, and conversion between ComparableHasher[X] and ComparableHasher[Y]
}

func (ComparableHasher[T]) Hash(h *Hash, v T) { WriteComparable(h, v) }
func (ComparableHasher[T]) Equal(x, y T) bool { return x == y }
