// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mapset defines operations on sets represented either as
// map[K]struct{}, or as map[K]bool where the boolean value is
// ignored.
//
// These types are common common choices for representing sets in
// existing Go code. In new code, we recommend using the
// [container/set.Set] type to indicate that a built-in map represents
// a set, as this avoids any possible ambiguity about boolean values.
//
// Beware that, when using a bool-valued map, it is the set of keys
// that matters, not the Boolean values. When an operation in this
// package creates a new bool-valued map, its values will always be
// true, regardless of any Boolean values among the operands.
//
// Binary operations such as [Union] may be applied to operators of
// the same or different types; the type of the result matches the
// first operand.
//
// This package should not be used with maps whose keys are (or
// contain) floating-point numbers whose value may be NaN, as it may
// take shortcuts based on algebraic identities that assume k == k for
// any key k. For example, Intersects(m, m) returns true for any
// non-empty map, without inspecting the map keys.
//
// Example:
//
//	a := mapset.Of(1, 2, 3)
//	b := mapset.Of(3, 4, 5)
//	fmt.Println(mapset.String(mapset.Union(a, b))) // {1, 2, 3, 4, 5}
package mapset

import (
	"github.com/ramalho/sets-in-go/vendored/internal/maps"
)

/// ...

// -- queries --

/// ...

// Contains reports whether set x contains key k.
func Contains[M ~map[K]V, K comparable, V bool | struct{}](x M, k K) bool {
	_, ok := x[k]
	return ok
}

// -- binary operations --

// / ...
// Intersection returns a new map containing the intersection of x and y.
func Intersection[MX ~map[K]VX, MY ~map[K]VY, K comparable, VX, VY bool | struct{}](x MX, y MY) MX {
	z := make(MX)

	if maps.Identical(x, y) {
		copy(z, x)
		return z
	}

	// Iterate over the smaller of the two maps.
	if len(x) < len(y) {
		for k := range x {
			if Contains(y, k) {
				insert(z, k)
			}
		}
	} else {
		for k := range y {
			if Contains(x, k) {
				insert(z, k)
			}
		}
	}
	return z
}

/// ...

func copy[MD ~map[K]VD, MS ~map[K]VS, K comparable, VD, VS bool | struct{}](dst MD, src MS) {
	// Avoid maps.Clone, which may return nil,
	// and may propagate 'false' values.
	for k := range src {
		insert(dst, k)
	}
}

// IntersectionWith updates x to the [Intersection] of x and y.
func IntersectionWith[M ~map[K]V, K comparable, V bool | struct{}](x, y M) {
	if maps.Identical(x, y) {
		return // x ∩ x = x
	}
	for k := range x {
		if !Contains(y, k) {
			delete(x, k)
		}
	}
}
