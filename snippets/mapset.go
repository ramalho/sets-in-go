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
//    a := mapset.Of(1, 2, 3)
//    b := mapset.Of(3, 4, 5)
//    fmt.Println(mapset.String(mapset.Union(a, b))) // {1, 2, 3, 4, 5}
package mapset

import (
    "fmt"
    "github.com/ramalho/sets-in-go/vendored/internal/fmtsort"
    "github.com/ramalho/sets-in-go/vendored/internal/maps"
    "iter"
    "reflect"
    "strings"
)

// Collect returns a new set containing the elements of the sequence.
func Collect[K comparable](seq iter.Seq[K]) map[K]struct{} {
    return collect[K, struct{}](seq)
}

// CollectBool returns a new set containing the elements of the sequence.
// The map values are all "true".
func CollectBool[K comparable](seq iter.Seq[K]) map[K]bool {
    return collect[K, bool](seq)
}

func collect[K comparable, V bool | struct{}](seq iter.Seq[K]) map[K]V {
    x := make(map[K]V)
    InsertAll(x, seq)
    return x
}

// Of creates a new set containing the elements of the sequence.
func Of[K comparable](elems ...K) map[K]struct{} {
    return of[K, struct{}](elems...)
}

// OfBool creates a new set containing the elements of the sequence.
// The map values are all "true".
func OfBool[K comparable](elems ...K) map[K]bool {
    return of[K, bool](elems...)
}

func of[K comparable, V bool | struct{}](elems ...K) map[K]V {
    x := make(map[K]V, len(elems))
    for _, elem := range elems {
        insert(x, elem)
    }
    return x
}

// String returns a representation of the map as a string of the form
// "{a, ..., z}". Keys are sorted in the same manner as [fmt.Sprint]
// of a map, and each is printed as if by fmt.Sprint.
// The map values are ignored.
func String[M ~map[K]V, K comparable, V bool | struct{}](x M) string {
    var buf strings.Builder
    buf.WriteByte('{')
    for i, entry := range fmtsort.Sort(reflect.ValueOf(x)) {
        if i > 0 {
            buf.WriteString(", ")
        }
        fmt.Fprint(&buf, entry.Key.Interface())
    }
    buf.WriteByte('}')
    return buf.String()
}

// -- queries --

// Equal reports whether sets x and y contain the same set of keys.
// The map values are ignored.
func Equal[M ~map[K]V, K comparable, V bool | struct{}](x, y M) bool {
    return maps.Identical(x, y) ||
        len(x) == len(y) && ContainsAll(x, maps.Keys(y))
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

// Contains reports whether set x contains key k.
func Contains[M ~map[K]V, K comparable, V bool | struct{}](x M, k K) bool {
    _, ok := x[k]
    return ok
}

// All returns an iterator over the sequence of elements of x.
// The sequence order is undefined and likely random.
func All[M ~map[K]V, K comparable, V bool | struct{}](x M) iter.Seq[K] {
    return maps.Keys(x)
}

// Intersects reports whether Intersection(x, y) is non-empty.
func Intersects[MX, MY ~map[K]V, K comparable, V bool | struct{}](x MX, y MY) bool {
    if maps.Identical(x, y) {
        return len(x) > 0
    }

    // Iterate over the smaller of the two maps.
    if len(y) < len(x) {
        for k := range y {
            if Contains(x, k) {
                return true
            }
        }
    } else {
        for k := range x {
            if Contains(y, k) {
                return true
            }
        }
    }
    return false
}

// -- binary operations --

// Union returns a new map containing the union of x and y.
func Union[MX ~map[K]VX, MY ~map[K]VY, K comparable, VX, VY bool | struct{}](x MX, y MY) MX {
    z := make(MX)
    copy(z, x)
    if !maps.Identical(x, y) {
        copy(z, y)
    }
    return z
}

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

// Difference returns a new map containing the elements of x that are
// not present in y.
func Difference[MX ~map[K]VX, MY ~map[K]VY, K comparable, VX, VY bool | struct{}](x MX, y MY) MX {
    z := make(MX)
    diff(z, x, y)
    return z
}

// SymmetricDifference returns a new map containing the elements of x
// that are not present in y, and the elements of y that are not
// present in x.
func SymmetricDifference[MX ~map[K]VX, MY ~map[K]VY, K comparable, VX, VY bool | struct{}](x MX, y MY) MX {
    z := make(MX)
    diff(z, x, y)
    diff(z, y, x)
    return z
}

func diff[MX ~map[K]VX, MY ~map[K]VY, MZ ~map[K]VZ, K comparable, VX, VY, VZ bool | struct{}](z MZ, x MX, y MY) {
    if !maps.Identical(x, y) {
        for k := range x {
            if !Contains(y, k) {
                insert(z, k)
            }
        }
    }
}

// -- mutations --

// Insert adds elem element to the set.
// If the set values are boolean, the value 'true' is used.
// It reports whether len(x) changed.
func Insert[M ~map[K]V, K comparable, V bool | struct{}](x M, elem K) bool {
    pre := len(x)
    insert(x, elem)
    return len(x) != pre
}

func insert[M ~map[K]V, K comparable, V bool | struct{}](m M, k K) {
    // Choose the distinguished "present" value (true or struct{}{}).
    // This compiles to a load from .rodata.
    var present V
    if _, ok := any(present).(bool); ok {
        present = any(true).(V)
    }

    // This is the canonical insertion operation.
    // All maps created by this API use only the
    // distinguished 'present' value for the result type.
    m[k] = present
}

// InsertAll adds each element of the addenda sequence to the set.
// If the set values are boolean, the value 'true' is used.
// It reports whether len(x) changed.
func InsertAll[M ~map[K]V, K comparable, V bool | struct{}](x M, addenda iter.Seq[K]) bool {
    pre := len(x)
    for k := range addenda {
        insert(x, k)
    }
    return len(x) != pre
}

// Delete removes an element from the set.
// It reports whether len(x) changed.
func Delete[M ~map[K]V, K comparable, V bool | struct{}](x M, k K) bool {
    pre := len(x)
    delete(x, k)
    return len(x) != pre
}

// DeleteAll removes from set x each element of the delenda sequence.
// It reports whether len(x) changed.
func DeleteAll[M ~map[K]V, K comparable, V bool | struct{}](x M, delenda iter.Seq[K]) bool {
    pre := len(x)
    for k := range delenda {
        delete(x, k)
    }
    return len(x) != pre
}

// DeleteFunc removes each element k of set x such that f(k).
// It reports whether len(x) changed.
func DeleteFunc[M ~map[K]V, K comparable, V bool | struct{}](x M, f func(K) bool) bool {
    pre := len(x)
    for k := range x {
        if f(k) {
            delete(x, k)
        }
    }
    return len(x) != pre
}

// -- in-place binary updates --

// UnionWith updates x to the [Union] of x and y.
//
// (If V is boolean, the values of y "win".)
func UnionWith[MX ~map[K]VX, MY ~map[K]VY, K comparable, VX, VY bool | struct{}](x MX, y MY) {
    if maps.Identical(x, y) {
        return // x ∪ x = x
    }
    copy(x, y)
}

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

// DifferenceWith updates x to the [Difference] of x and y.
// In other words, it removes from x all the elements of y.
func DifferenceWith[M ~map[K]V, K comparable, V bool | struct{}](x, y M) {
    if maps.Identical(x, y) {
        clear(x) // x \ x = Ø
        return
    }
    for k := range y {
        delete(x, k)
    }
}

// SymmetricDifferenceWith updates x to the [SymmetricDifference] of x and y.
func SymmetricDifferenceWith[M ~map[K]V, K comparable, V bool | struct{}](x, y M) {
    if maps.Identical(x, y) {
        clear(x) // x ∆ x = Ø
        return
    }
    for k := range y {
        if Contains(x, k) {
            delete(x, k)
        } else {
            insert(x, k)
        }
    }
}
