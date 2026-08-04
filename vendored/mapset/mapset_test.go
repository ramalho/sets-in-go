// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapset_test

import (
	"fmt"
	"github.com/ramalho/sets-in-go/vendored/internal/maps"
	. "github.com/ramalho/sets-in-go/vendored/mapset"
	"iter"
	"slices"
	"strings"
	"testing"
)

func Example() {
	a := OfBool(1, 2, 3)
	b := CollectBool(seq(3, 4, 5))

	fmt.Printf("a                              = %v\n", String(a))
	fmt.Printf("b                              = %v\n", String(b))
	fmt.Printf("Intersection(a, b)             = %v\n", String(Intersection(a, b)))
	fmt.Printf("Union(a, b)                    = %v\n", String(Union(a, b)))
	fmt.Printf("Difference(a, b)               = %v\n", String(Difference(a, b)))
	fmt.Printf("Difference(b, a)               = %v\n", String(Difference(b, a)))
	fmt.Printf("SymmetricDifference(a, b)      = %v\n", String(SymmetricDifference(a, b)))
	fmt.Printf("Intersects(a, b)               = %t\n", Intersects(a, b))
	fmt.Printf("Equal(a, b)                    = %t\n", Equal(a, b))
	fmt.Printf("Equal(a, a)                    = %t\n", Equal(a, a))
	fmt.Printf("Contains(a, 3)                 = %t\n", Contains(a, 3))

	// Output:
	// a                              = {1, 2, 3}
	// b                              = {3, 4, 5}
	// Intersection(a, b)             = {3}
	// Union(a, b)                    = {1, 2, 3, 4, 5}
	// Difference(a, b)               = {1, 2}
	// Difference(b, a)               = {4, 5}
	// SymmetricDifference(a, b)      = {1, 2, 4, 5}
	// Intersects(a, b)               = true
	// Equal(a, b)                    = false
	// Equal(a, a)                    = true
	// Contains(a, 3)                 = true
}

func TestBinaryPredicates(t *testing.T) {
	var (
		s123     = OfBool(1, 2, 3)
		s123copy = OfBool(1, 2, 3)
		s345     = OfBool(3, 4, 5)
		s45      = OfBool(4, 5)
		s456     = OfBool(4, 5, 6)
	)

	for _, test := range []struct {
		x, y map[int]bool
		want string
	}{
		{s123, s123, "Identical Equal ContainsAll Intersects"},
		{s123, s123copy, "Equal ContainsAll Intersects"},
		{s123, s345, "Intersects"},
		{s345, s123, "Intersects"},
		{s123, nil, "ContainsAll"},
		{nil, s123, ""},
		{s123copy, s123, "Equal ContainsAll Intersects"},
		{s123copy, s123copy, "Identical Equal ContainsAll Intersects"},
		{s123, s456, ""},
		{s456, s123, ""},
		{s345, s456, "Intersects"},
		{s456, s345, "Intersects"},
		{s456, s45, "ContainsAll Intersects"},
		{s45, s456, "Intersects"},
	} {
		var attrs []string
		if maps.Identical(test.x, test.y) {
			attrs = append(attrs, "Identical")
		}
		if Equal(test.x, test.y) {
			attrs = append(attrs, "Equal")
		}
		if ContainsAll(test.x, maps.Keys(test.y)) {
			attrs = append(attrs, "ContainsAll")
		}
		if Intersects(test.x, test.y) {
			attrs = append(attrs, "Intersects")
		}
		if got := strings.Join(attrs, " "); got != test.want {
			t.Errorf("%s x %s: got %q, want %q", String(test.x), String(test.y), got, test.want)
		}
	}
}

func TestUnit(t *testing.T) {
	a := Of(1, 2, 3)
	b := Of(3, 4, 5)
	c := Intersection(a, b)
	if got, want := String(c), "{3}"; got != want {
		t.Errorf("Intersection(%v, %v) = %v, want %s", String(a), String(b), got, want)
	}

	if got := Collect(All(a)); !Equal(a, got) {
		t.Errorf("a=%s: CollectUnit(All(a)) = %s, want Equal(a, %s)", String(a), String(got), String(a))
	}
}

func TestMutations(t *testing.T) {
	t.Run("InsertAll", func(t *testing.T) {
		m := OfBool(1, 2, 3)
		if !InsertAll(m, seq(3, 4, 5)) {
			t.Errorf("InsertAll did not report change")
		}
		if want := OfBool(1, 2, 3, 4, 5); !Equal(m, want) {
			t.Errorf("got %v, want %v", String(m), String(want))
		}
		if InsertAll(m, seq(1, 2, 3)) {
			t.Errorf("InsertAll reported change for no-op")
		}
	})

	t.Run("DeleteAll", func(t *testing.T) {
		m := OfBool(1, 2, 3)
		if !DeleteAll(m, seq(1, 3)) {
			t.Errorf("DeleteAll did not report change")
		}
		if want := OfBool(2); !Equal(m, want) {
			t.Errorf("got %v, want %v", String(m), String(want))
		}
		if DeleteAll(m, seq(4, 5)) {
			t.Errorf("DeleteAll reported change for no-op")
		}
	})

	t.Run("DeleteFunc", func(t *testing.T) {
		m := OfBool(1, 2, 3, 4)
		isEven := func(k int) bool { return k%2 == 0 }
		if !DeleteFunc(m, isEven) {
			t.Errorf("DeleteFunc did not report change")
		}
		if want := OfBool(1, 3); !Equal(m, want) {
			t.Errorf("got %v, want %v", String(m), String(want))
		}
		if DeleteFunc(m, isEven) {
			t.Errorf("DeleteFunc reported change for no-op")
		}
	})
}

func TestInPlaceUpdates(t *testing.T) {
	check := func(t *testing.T, m map[int]bool, want ...int) {
		t.Helper()
		if w := OfBool(want...); !Equal(m, w) {
			t.Errorf("got %v, want %v", String(m), String(w))
		}
	}

	t.Run("UnionWith", func(t *testing.T) {
		x := OfBool(1, 2)
		y := OfBool(2, 3)
		UnionWith(x, y)
		check(t, x, 1, 2, 3)

		// Self-aliasing
		UnionWith(x, x)
		check(t, x, 1, 2, 3)
	})

	t.Run("IntersectionWith", func(t *testing.T) {
		x := OfBool(1, 2, 3)
		y := OfBool(2, 3, 4)
		IntersectionWith(x, y)
		check(t, x, 2, 3)

		// Self-aliasing
		IntersectionWith(x, x)
		check(t, x, 2, 3)
	})

	t.Run("DifferenceWith", func(t *testing.T) {
		x := OfBool(1, 2, 3)
		y := OfBool(2, 4)
		DifferenceWith(x, y)
		check(t, x, 1, 3)

		// Self-aliasing
		DifferenceWith(x, x)
		check(t, x)
	})

	t.Run("SymmetricDifferenceWith", func(t *testing.T) {
		x := OfBool(1, 2, 3)
		y := OfBool(2, 3, 4)
		SymmetricDifferenceWith(x, y)
		check(t, x, 1, 4)

		// Self-aliasing
		SymmetricDifferenceWith(x, x)
		check(t, x)
	})
}

func TestSecondOperandWins(t *testing.T) {
	// x := map[int]bool{1: true, 2: "x2"}
	// y := map[int]bool{2: "y2", 3: "y3"}

	// t.Run("Union", func(t *testing.T) {
	// 	got := Union(x, y)
	// 	want := map[int]string{1: "x1", 2: "y2", 3: "y3"}
	// 	if !Equal(got, want) || got[2] != "y2" {
	// 		t.Errorf("Union: got %v, want %v", got, want)
	// 	}
	// })

	// t.Run("Intersection", func(t *testing.T) {
	// 	got := Intersection(x, y)
	// 	want := map[int]string{2: "y2"}
	// 	if !Equal(got, want) || got[2] != "y2" {
	// 		t.Errorf("Intersection: got %v, want %v", got, want)
	// 	}
	// })

	// t.Run("UnionWith", func(t *testing.T) {
	// 	xc := maps.Clone(x)
	// 	UnionWith(xc, y)
	// 	want := map[int]string{1: "x1", 2: "y2", 3: "y3"}
	// 	if !Equal(xc, want) || xc[2] != "y2" {
	// 		t.Errorf("UnionWith: got %v, want %v", xc, want)
	// 	}
	// })
}

func TestNil(t *testing.T) {
	var n map[int]bool
	s := OfBool(1)

	if got := Union(n, n); len(got) > 0 {
		t.Errorf("Union(nil, nil) = %v, want empty", String(got))
	}
	if got := Union(n, s); !Equal(got, s) {
		t.Errorf("Union(nil, s) = %v, want %v", String(got), String(s))
	}
	if got := Union(s, n); !Equal(got, s) {
		t.Errorf("Union(s, nil) = %v, want %v", String(got), String(s))
	}

	if got := Intersection(n, s); len(got) > 0 {
		t.Errorf("Intersection(nil, s) = %v, want empty", String(got))
	}
	if got := Intersection(s, n); len(got) > 0 {
		t.Errorf("Intersection(s, nil) = %v, want empty", String(got))
	}

	if got, want := String[map[int]bool](nil), "{}"; got != want {
		t.Errorf("String(nil) = %q, want %q", got, want)
	}
}

func seq[T any](values ...T) iter.Seq[T] { return slices.Values(values) }
