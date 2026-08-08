package mapset

import (
	"iter"

	"github.com/ramalho/sets-in-go/vendored/internal/maps"
)

/// ...

// -- queries -- (reordered to suit this presentation)

// Contains reports whether set x contains key k.
func Contains[M ~map[K]V, K comparable, V bool | struct{}](x M, k K) bool {
	_, ok := x[k]
	return ok
}

// ContainsAll reports whether set x contains all the keys in the sequence.
// The map values are ignored.
func ContainsAll[M ~map[K]V, K comparable, V bool | struct{}](x M, keys iter.Seq[K]) bool {
	for k := range keys {
		if !Contains(x, k) {
			return false
		}
	}
	return true
}

// Equal reports whether sets x and y contain the same set of keys.
// The map values are ignored.
func Equal[M ~map[K]V, K comparable, V bool | struct{}](x, y M) bool {
	return maps.Identical(x, y) ||
		len(x) == len(y) && ContainsAll(x, maps.Keys(y))
}
