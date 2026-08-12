// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package container_test

import (
	"fmt"
	"hash/maphash"
	"slices"

	"github.com/ramalho/sets-in-go/vendored/hash"
	"github.com/ramalho/sets-in-go/vendored/set"
)

// Example_setOfVersusCollect compares set.Of and set.Collect on the same data.
// set.Of returns the bare map[int]struct{} type.
// set.Of is variadic, so a slice of elements can be passed to it too, spread
// with "...", instead of listing the arguments individually.
// set.Collect returns set.Set[int] itself, so it prints via Set's String method
// instead.
// The items are sorted in the String methods of each type,
// to make it easier to test.

func Example_setOfVersusCollect() {
	o := set.Of(3, 2, 5, 1, 5)
	fmt.Printf("%T %v\n", o, o)

	ints := []int{3, 2, 5, 1, 5}
	o2 := set.Of(ints...)
	fmt.Printf("%T %v\n", o2, o2)

	c := set.Collect(slices.Values(ints))
	fmt.Printf("%T %v\n", c, c)

	// Output:
	// map[int]struct {} map[1:{} 2:{} 3:{} 5:{}]
	// map[int]struct {} map[1:{} 2:{} 3:{} 5:{}]
	// set.Set[int] {1, 2, 3, 5}
}

// Example_unionOfSetSet declares a and b with the abstract
// _AbstractSet[int, set.Set[int]] interface type, and assigns each
// one a concrete set.Set[int] built by set.Collect.
//
// This shows that a concrete set.Set satisfies the F-bounded interface.
//
// Union's signature, x.Union(S) S, requires its argument as
// the concrete type S, not the interface, so b must be asserted back
// to set.Set[int] to call a.Union(b). The result u comes out as the
// concrete set.Set[int], not _AbstractSet, as shown by %T.
func Example_unionOfSetSet() {
	var a _AbstractSet[int, set.Set[int]] = set.Collect(slices.Values([]int{1, 2, 3}))
	var b _AbstractSet[int, set.Set[int]] = set.Collect(slices.Values([]int{2, 3, 4}))
	u := a.Union(b.(set.Set[int]))
	fmt.Printf("%T %v\n", u, u)

	// Output:
	// set.Set[int] {1, 2, 3, 4}
}

// Example_unionOfHashSet is the *hash.Set counterpart of
// Example_unionOfSetSet: same elements, but united with *hash.Set's
// own Union method, whose result type is *hash.Set[int] rather than
// set.Set[int].
func Example_unionOfHashSet() {
	var a set.Set[int] = set.Of(1, 2, 3)
	var b set.Set[int] = set.Of(2, 3, 4)

	hasher := maphash.ComparableHasher[int]{}
	c := hash.NewSet(hasher)
	c.InsertAll(a.All())
	d := hash.NewSet(hasher)
	d.InsertAll(b.All())
	u := c.Union(d)
	fmt.Printf("%T %v\n", u, u)

	// Output:
	// *hash.Set[int] {1, 2, 3, 4}
}
