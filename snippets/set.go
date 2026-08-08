// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package set defines [Set], the canonical representation of a set
// recommended for use in public APIs.
//
// A Set[E] is transparently represented as a map[E]struct{}, and so
// it may be used as the operand of a range loop.
//
// Example:
//
//    var x, y, z set.Set[int]
//    x = set.Of(1, 2, 3)        // "{1, 2, 3}"
//    y = set.Of(3, 4, 5)        // "{3, 4, 5"}
//    z = x.Intersection(y)        // "{3}"
//    for elem := range z {
//        println(elem)        // "3"
//    }
//
// See also [container/mapset], which provides helpers for working
// with legacy sets represented as unnamed map types.
package set

import (
    "github.com/ramalho/sets-in-go/vendored/mapset"
    "iter"
    "maps"
)

// A Set[E] is a set of elements of type E.
type Set[E comparable] map[E]struct{}

// Collect creates a new set containing the elements of the sequence.
func Collect[E comparable](seq iter.Seq[E]) Set[E] {
    return Set[E](mapset.Collect(seq))
}

// Of creates a new set containing the elements of the sequence.
func Of[E comparable](elems ...E) map[E]struct{} {
    return Set[E](mapset.Of(elems...))
}

// String returns a representation of the set as a string of the form
// "{a, ..., z}". Elements are sorted in the same manner as [fmt.Sprint]
// of a map, and each is printed as if by fmt.Sprint.
func (x Set[E]) String() string {
    return mapset.String(x)
}

// Clone returns a new non-nil set with the same elements as s.
func (s Set[E]) Clone() Set[E] {
    if s == nil {
        return make(Set[E])
    }
    return maps.Clone(s)
}

// Clear removes all elements from the map.
func (s Set[E]) Clear() {
    clear(s)
}

// -- queries --

// Len returns the number of set elements.
func (x Set[E]) Len() int {
    // (This method is provided so that the method
    // set forms a complete abstract data type.)
    return len(x)
}

// Equal reports whether sets x and y contain the same set of elements.
func (x Set[E]) Equal(y Set[E]) bool {
    return mapset.Equal(x, y)
}

// ContainsAll reports whether set x contains all the elements of the sequence.
func (x Set[E]) ContainsAll(seq iter.Seq[E]) bool {
    return mapset.ContainsAll(x, seq)
}

// Contains reports whether set x contains element e.
func (x Set[E]) Contains(e E) bool {
    return mapset.Contains(x, e)
}

// All returns an iterator over the sequence of elements of x.
// The sequence order is undefined and likely random.
func (x Set[E]) All() iter.Seq[E] {
    return maps.Keys(x)
}

// Intersects reports whether x.Intersection(y) is non-empty.
func (x Set[E]) Intersects(y Set[E]) bool {
    return mapset.Intersects(x, y)
}

// -- binary operations --

// Union returns a new map containing the union of x and y.
func (x Set[E]) Union(y Set[E]) Set[E] {
    return mapset.Union(x, y)
}

// Intersection returns a new map containing the intersection of x and y.
func (x Set[E]) Intersection(y Set[E]) Set[E] {
    return mapset.Intersection(x, y)
}

// Difference returns a new map containing the elements of x that are
// not present in y.
func (x Set[E]) Difference(y Set[E]) Set[E] {
    return mapset.Difference(x, y)
}

// SymmetricDifference returns a new map containing the elements of x
// that are not present in y, and the elements of y that are not
// present in x.
func (x Set[E]) SymmetricDifference(y Set[E]) Set[E] {
    return mapset.SymmetricDifference(x, y)
}

// -- mutations --

// Insert adds the specified element to the set.
// It reports whether len(x) changed.
func (x Set[E]) Insert(elem E) bool {
    return mapset.Insert(x, elem)
}

// InsertAll adds each element of the addenda sequence to the set.
// It reports whether len(x) changed.
func (x Set[E]) InsertAll(addenda iter.Seq[E]) bool {
    return mapset.InsertAll(x, addenda)
}

// Delete removes element e from the set.
// It reports whether len(x) changed.
func (x Set[E]) Delete(e E) bool {
    return mapset.Delete(x, e)
}

// DeleteAll removes from set x each element of the delenda sequence.
// It reports whether len(x) changed.
func (x Set[E]) DeleteAll(delenda iter.Seq[E]) bool {
    return mapset.DeleteAll(x, delenda)
}

// DeleteFunc removes each element e of set x such that f(e).
// It reports whether len(x) changed.
func (x Set[E]) DeleteFunc(f func(E) bool) bool {
    return mapset.DeleteFunc(x, f)
}

// -- in-place binary updates --

// UnionWith updates x to the [Union] of x and y.
func (x Set[E]) UnionWith(y Set[E]) {
    mapset.UnionWith(x, y)
}

// IntersectionWith updates x to the [Intersection] of x and y.
func (x Set[E]) IntersectionWith(y Set[E]) {
    mapset.IntersectionWith(x, y)
}

// DifferenceWith updates x to the [Difference] of x and y.
// In other words, it removes from x all the elements of y.
func (x Set[E]) DifferenceWith(y Set[E]) {
    mapset.DifferenceWith(x, y)
}

// SymmetricDifferenceWith updates x to the [SymmetricDifference] of x and y.
func (x Set[E]) SymmetricDifferenceWith(y Set[E]) {
    mapset.SymmetricDifferenceWith(x, y)
}
