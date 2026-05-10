// Package genset provides a generic Set type.
package genset

import (
	"bytes"
	"fmt"
	"iter"
	"sort"
	"strings"
)

// Set of elements of any comparable type.
type Set[E comparable] struct {
	store map[E]struct{}
}

// Make creates and returns a new Set.
func Make[E comparable](elems ...E) Set[E] {
	s := Set[E]{store: make(map[E]struct{})}
	for _, elem := range elems {
		s.store[elem] = struct{}{}
	}
	return s
}

// MakeFromText creates and returns a new Set[string] from
// a string of elements separated by whitespace.
func MakeFromText(text string) Set[string] {
	return Make(strings.Fields(text)...)
}

// Len reports the number of elements in the set.
func (s Set[E]) Len() int {
	return len(s.store)
}

// Contains reports whether set contains the element.
// Math: S ∋ e.
func (s Set[E]) Contains(elem E) bool {
	_, found := s.store[elem]
	return found
}

// ContainsAll reports whether s contains all the given elements.
func (s Set[E]) ContainsAll(elems ...E) bool {
	for _, elem := range elems {
		if _, found := s.store[elem]; !found {
			return false
		}
	}
	return true
}

// ToSlice returns a new slice with the elements of s.
// The order of the elements is undefined.
func (s Set[E]) ToSlice() []E {
	elems := make([]E, 0, len(s.store))
	for k := range s.store {
		elems = append(elems, k)
	}
	return elems
}

// String returns a string representation of s with
// elements in lexicographic order.
func (s Set[E]) String() string {
	strs := make([]string, 0, len(s.store))
	for elem := range s.store {
		strs = append(strs, fmt.Sprint(elem))
	}
	sort.Strings(strs)
	var buf bytes.Buffer
	buf.WriteString("Set{")
	buf.WriteString(strings.Join(strs, " "))
	buf.WriteByte('}')
	return buf.String()
}

// allIn reports whether all elements of s exist in other.
func (s Set[E]) allIn(other Set[E]) bool {
	for elem := range s.store {
		if _, found := other.store[elem]; !found {
			return false
		}
	}
	return true
}

// Equal reports whether set is equal to other.
func (s Set[E]) Equal(other Set[E]) bool {
	return len(s.store) == len(other.store) && s.allIn(other)
}

// Copy returns a new Set: a copy of s.
func (s Set[E]) Copy() Set[E] {
	res := Make[E]()
	for elem := range s.store {
		res.store[elem] = struct{}{}
	}
	return res
}

// All returns an iterator over the elements of s,
// compatible with for/range in Go 1.23 and later.
func (s Set[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for elem := range s.store {
			if !yield(elem) {
				return
			}
		}
	}
}
