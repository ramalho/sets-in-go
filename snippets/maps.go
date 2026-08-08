// Package maps stands in for the standard library's maps package, adding
// Identical, which mapset depends on but which is not yet in any released Go.
//
// mapset.go imports this instead of "maps" so that its function bodies can be
// kept verbatim. Identical is copied from the proposal CL:
// https://go-review.googlesource.com/c/go/+/760800 (golang/go#78456).
//
// Once maps.Identical lands in the standard library, delete this package and
// drop the import rewrite for it in refresh.sh.
package maps

import (
	"iter"
	"maps"
	"unsafe"
)

// Identical reports whether two maps refer to the same data structure.
//
// Beware that some shortcuts based on Identical(x, y) may have surprising
// behavior for maps containing floating-point NaNs, since NaN != NaN.
func Identical[MX ~map[K]VX, MY ~map[K]VY, K comparable, VX, VY any](x MX, y MY) bool {
	// Maps in Go are references yet the core language
	// provides no safe way to ask whether they alias.
	type pointer = unsafe.Pointer
	return *(*pointer)(pointer(&x)) == *(*pointer)(pointer(&y))
}

// Clone is [maps.Clone].
func Clone[M ~map[K]V, K comparable, V any](m M) M { return maps.Clone(m) }

// Keys is [maps.Keys].
func Keys[M ~map[K]V, K comparable, V any](m M) iter.Seq[K] { return maps.Keys(m) }
