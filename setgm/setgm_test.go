package setgm_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ramalho/sets-in-go/setgm"
	"github.com/ramalho/sets-in-go/vendored/set"
)

// SetGM still satisfies interfaces built from its non-generic methods.
var _ fmt.Stringer = setgm.SetGM[int]{}

// The point of the package: one receiver, four different operand types.
//
// Each of these is a distinct instantiation of the method's own type
// parameters. set.Set[string] and SetGM[string] are named types over
// map[string]struct{}; the bool-valued map is the legacy representation, whose
// values are ignored.
func Example_heterogeneousOperands() {
	x := setgm.Of("a", "b")

	fmt.Println(x.Union(setgm.Of("c")))                // SetGM[string]
	fmt.Println(x.Union(set.Of("d")))                  // the proposal's set.Set
	fmt.Println(x.Union(map[string]struct{}{"e": {}})) // plain map-set
	fmt.Println(x.Union(map[string]bool{"f": true}))   // legacy bool map-set
	fmt.Println(x.Union(map[string]bool{"g": false}))  // values ignored, not keys
	// Output:
	// {a, b, c}
	// {a, b, d}
	// {a, b, e}
	// {a, b, f}
	// {a, b, g}
}

// The practical case: checking granted permissions held as a legacy
// map[string]bool against a modern set, with no conversion at the call site.
func Example_legacyPermissions() {
	granted := map[string]bool{"read": true, "write": true, "admin": true}
	required := setgm.Of("read", "write")

	fmt.Println(required.Subset(granted))   // are all required perms granted?
	fmt.Println(required.Superset(granted)) // ...and nothing more?
	fmt.Println(required.Difference(granted))
	fmt.Println(setgm.From(granted).Difference(required))
	// Output:
	// true
	// false
	// {}
	// {admin}
}

// A generic method cannot be used as a method value without instantiation
// ("cannot use generic function x.Union without instantiation"), but an
// explicitly instantiated one is an ordinary func value.
func Example_instantiatedMethodValue() {
	x := setgm.Of(1, 2)
	unionBool := x.Union[map[int]bool]

	fmt.Println(unionBool(map[int]bool{3: true}))
	// Output:
	// {1, 2, 3}
}

func TestBinaryOperations(t *testing.T) {
	x := setgm.Of(1, 2, 3)
	y := map[int]bool{3: true, 4: true, 5: true} // legacy operand throughout

	for _, tc := range []struct {
		name string
		got  setgm.SetGM[int]
		want string
	}{
		{"Union", x.Union(y), "{1, 2, 3, 4, 5}"},
		{"Intersection", x.Intersection(y), "{3}"},
		{"Difference", x.Difference(y), "{1, 2}"},
		{"SymmetricDifference", x.SymmetricDifference(y), "{1, 2, 4, 5}"},
	} {
		if got := tc.got.String(); got != tc.want {
			t.Errorf("%s = %s, want %s", tc.name, got, tc.want)
		}
	}

	if got := x.String(); got != "{1, 2, 3}" {
		t.Errorf("receiver mutated: x = %s", got)
	}
}

func TestPredicates(t *testing.T) {
	letters := setgm.Of('a', 'b', 'c', 'e')
	vowels := map[rune]bool{'a': true, 'e': true}

	for _, tc := range []struct {
		name string
		got  bool
		want bool
	}{
		{"letters ⊇ vowels", letters.Superset(vowels), true},
		{"letters ⊃ vowels", letters.ProperSuperset(vowels), true},
		{"letters ⊆ vowels", letters.Subset(vowels), false},
		{"letters = vowels", letters.Equal(vowels), false},
		{"letters ∩ vowels ≠ Ø", letters.Intersects(vowels), true},

		// ⊇ and ⊆ are reflexive; ⊃ and ⊂ are not.
		{"letters ⊇ letters", letters.Superset(letters), true},
		{"letters ⊆ letters", letters.Subset(letters), true},
		{"letters ⊃ letters", letters.ProperSuperset(letters), false},
		{"letters ⊂ letters", letters.ProperSubset(letters), false},
		{"letters = letters", letters.Equal(letters), true},

		// Equality is insensitive to both the map type and the bool values.
		{"vowels = vowels'", setgm.From(vowels).Equal(map[rune]struct{}{'a': {}, 'e': {}}), true},
		{"false value still a member", setgm.Of('z').Equal(map[rune]bool{'z': false}), true},

		// Disjoint sets.
		{"letters ∩ {x,y} = Ø", letters.Intersects(map[rune]bool{'x': true, 'y': true}), false},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %v, want %v", tc.name, tc.got, tc.want)
		}
	}
}

// The identity fast paths must agree with the general algorithm. Aliasing is
// what makes these worth testing: x and y are the same map, so the *With
// operations would otherwise mutate the map they are ranging over.
func TestIdentityFastPaths(t *testing.T) {
	for _, tc := range []struct {
		name string
		op   func(setgm.SetGM[int])
		want string
	}{
		{"UnionWith self", func(s setgm.SetGM[int]) { s.UnionWith(s) }, "{1, 2, 3}"},
		{"IntersectionWith self", func(s setgm.SetGM[int]) { s.IntersectionWith(s) }, "{1, 2, 3}"},
		{"DifferenceWith self", func(s setgm.SetGM[int]) { s.DifferenceWith(s) }, "{}"},
		{"SymmetricDifferenceWith self", func(s setgm.SetGM[int]) { s.SymmetricDifferenceWith(s) }, "{}"},
	} {
		s := setgm.Of(1, 2, 3)
		tc.op(s)
		if got := s.String(); got != tc.want {
			t.Errorf("%s: got %s, want %s", tc.name, got, tc.want)
		}
	}

	// A distinct map with equal contents must take the general path and reach
	// the same answers.
	s := setgm.Of(1, 2, 3)
	clone := s.Clone()
	s.SymmetricDifferenceWith(clone)
	if got := s.String(); got != "{}" {
		t.Errorf("SymmetricDifferenceWith(clone) = %s, want {}", got)
	}
}

func TestInPlaceUpdates(t *testing.T) {
	y := map[int]bool{3: true, 4: true} // legacy operand throughout

	for _, tc := range []struct {
		name string
		op   func(setgm.SetGM[int])
		want string
	}{
		{"UnionWith", func(s setgm.SetGM[int]) { s.UnionWith(y) }, "{1, 2, 3, 4}"},
		{"IntersectionWith", func(s setgm.SetGM[int]) { s.IntersectionWith(y) }, "{3}"},
		{"DifferenceWith", func(s setgm.SetGM[int]) { s.DifferenceWith(y) }, "{1, 2}"},
		{"SymmetricDifferenceWith", func(s setgm.SetGM[int]) { s.SymmetricDifferenceWith(y) }, "{1, 2, 4}"},
	} {
		s := setgm.Of(1, 2, 3)
		tc.op(s)
		if got := s.String(); got != tc.want {
			t.Errorf("%s = %s, want %s", tc.name, got, tc.want)
		}
	}
}

func TestMutations(t *testing.T) {
	s := setgm.Of(1, 2, 3)

	if !s.Insert(4) || s.Insert(4) {
		t.Error("Insert should report whether len changed")
	}
	if !s.Delete(4) || s.Delete(4) {
		t.Error("Delete should report whether len changed")
	}
	if !s.DeleteFunc(func(e int) bool { return e%2 == 0 }) {
		t.Error("DeleteFunc should report a change")
	}
	if got := s.String(); got != "{1, 3}" {
		t.Errorf("after DeleteFunc: %s, want {1, 3}", got)
	}

	if !s.InsertAll(setgm.Of(5, 7).All()) {
		t.Error("InsertAll should report a change")
	}
	if got, want := slices.Sorted(s.All()), []int{1, 3, 5, 7}; !slices.Equal(got, want) {
		t.Errorf("All = %v, want %v", got, want)
	}
	if s.Len() != 4 {
		t.Errorf("Len = %d, want 4", s.Len())
	}

	s.Clear()
	if s.Len() != 0 {
		t.Errorf("after Clear: Len = %d, want 0", s.Len())
	}
}

// The nil set is the empty set for every read-only operation: ranging over a
// nil map yields nothing, and lookups in one miss.
func TestNil(t *testing.T) {
	var none setgm.SetGM[int]
	some := setgm.Of(1, 2)

	for _, tc := range []struct {
		name string
		got  bool
		want bool
	}{
		{"none ⊆ some", none.Subset(some), true},
		{"none ⊇ some", none.Superset(some), false},
		{"some ⊇ none", some.Superset(none), true},
		{"none ⊇ none", none.Superset(none), true},
		{"none = none", none.Equal(none), true},
		{"none = some", none.Equal(some), false},
		{"none ∩ some", none.Intersects(some), false},
		{"none contains 1", none.Contains(1), false},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %v, want %v", tc.name, tc.got, tc.want)
		}
	}

	if got := none.Union(some).String(); got != "{1, 2}" {
		t.Errorf("nil.Union = %s, want {1, 2}", got)
	}
	if got := none.Len(); got != 0 {
		t.Errorf("nil.Len = %d, want 0", got)
	}
	if c := none.Clone(); c == nil {
		t.Error("Clone of nil set should be non-nil")
	}
}

// Of returns SetGM[E], not map[E]struct{}. Unlike the proposal's set.Of, the
// result can be used directly as a receiver without a helper variable.
func TestOfReturnsNamedType(t *testing.T) {
	if !setgm.Of(1, 2, 3).Superset(setgm.Of(1, 2)) {
		t.Error("Of(...).Superset(...) should hold")
	}
	if got := setgm.Collect(slices.Values([]int{1, 1, 2})).String(); got != "{1, 2}" {
		t.Errorf("Collect = %s, want {1, 2}", got)
	}
}

// From copies, so later mutation of the source does not show through.
func TestFromCopies(t *testing.T) {
	src := map[string]bool{"a": true}
	got := setgm.From(src)
	src["b"] = true

	if want := "{a}"; got.String() != want {
		t.Errorf("From = %s, want %s", got, want)
	}
}
