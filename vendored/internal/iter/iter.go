// Package iter stands in for the standard library's iter package when
// building the vendored container/hash sources.
//
// container/hash calls seq.Every and seq.Some, methods on iter.Seq proposed
// in CL 745440 and absent from every released Go. Methods cannot be attached
// to a type declared in another package, and redeclaring Seq here as a
// defined type would make it non-assignable to the standard iter.Seq that
// slices.Collect and the rest of the standard library expect. So Seq and Seq2
// are aliases for the real types, and Every and Some are free functions
// carrying the bodies from CL 745440. refresh.sh rewrites the three call
// sites in map.go and set.go from method to function form.
//
// This is the only file under vendored/hash's dependencies that is not
// upstream code.
package iter

import stditer "iter"

// Seq is an alias for the standard [iter.Seq].
type Seq[V any] = stditer.Seq[V]

// Seq2 is an alias for the standard [iter.Seq2].
type Seq2[K, V any] = stditer.Seq2[K, V]

// The bodies below are verbatim from CL 745440, with the receiver
// (seq Seq[V]) moved into the parameter list. The upstream comment
// explaining the desugared implementation is kept intact.

// Using the desugared implementations instead of 'range seq' runs
// about 2x faster than the code in the comments based on range loops.
// The cost of range loops comes from the compiler heap-allocating a
// control variable, and additional checks of well-behavedness w.r.t.
// concurrency, panics, stopping when yield returns false, etc.
// There's no value to these "middleware" iterators repeating these
// checks: the user's range loop (the ultimate consumer) will keep
// the underlying iterator 'seq' (the ultimate producer) honest.

// Every reports whether f(v) is true for every element in seq.
// It stops at the first element where f returns false.
func Every[V any](seq Seq[V], f func(V) bool) bool {
	every := true
	seq(func(v V) bool {
		if !f(v) {
			every = false
			return false
		}
		return true
	})
	return every

	// for v := range seq {
	// 	if !f(v) {
	// 		return false
	// 	}
	// }
	// return true
}

// Some reports whether f(v) is true for some element in seq.
// It stops at the first element where f returns true.
func Some[V any](seq Seq[V], f func(V) bool) bool {
	some := false
	seq(func(v V) bool {
		if f(v) {
			some = true
			return false
		}
		return true
	})
	return some

	// for v := range seq {
	// 	if f(v) {
	// 		return true
	// 	}
	// }
	// return false
}
