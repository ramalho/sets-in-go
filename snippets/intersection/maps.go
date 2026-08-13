package maps

import (
    "unsafe"
    /// ...
)

/// ...

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
