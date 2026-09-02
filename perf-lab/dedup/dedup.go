// Package dedup holds the deduplication rule from the unigo commands, written
// five ways, so they can be compared side by side and benchmarked against each
// other.
//
// Each function returns the distinct elements of lines, in order of first
// appearance — the same job the unigo commands do, with the I/O removed:
//
//   - [MapStruct] is ../../unigo/unigo-map: map[string]struct{}, lookup then store.
//   - [MapBool]   is ../../unigo/unigo-mapbool: the other legacy idiom, map[string]bool.
//   - [MapLen]    is ../../unigo/unigo-maplen: store, then check whether len changed.
//   - [Mapset]    is ../../unigo/unigo-mapset: a plain map, one call to mapset.Insert.
//   - [Set]       is ../../unigo/unigo-set: the proposed container/set type.
//
// The five differ only in how they ask "is this line new?", and they allocate
// exactly the same: set.Set[E] is defined as map[E]struct{}, so no wrapper
// type is involved anywhere.
//
// They do not, however, run at the same speed, and which one wins depends on
// the input. The two-step versions test first and write only on a miss. The
// one-step versions call mapset.Insert, which always assigns — it reports
// whether the set changed by comparing len before and after (see
// vendored/mapset/mapset.go). So:
//
//   - when most lines are new, one unconditional write beats a lookup
//     followed by a write, and Mapset/Set are the faster pair;
//   - when most lines are duplicates, the two-step versions take a read-only
//     path while Mapset/Set pay for a wasted write every time, and the
//     ranking reverses.
//
// See README.md for measurements. This is the same trade-off recorded in
// ../README.md for gen_haystack.go, where "just put it there" won because
// every generated number was new.
package dedup

import (
	"github.com/ramalho/sets-in-go/vendored/mapset"
	"github.com/ramalho/sets-in-go/vendored/set"
)

// MapStruct deduplicates using map[string]struct{}: the comma-ok lookup and
// the store are separate steps.
func MapStruct(lines []string) []string {
	seen := make(map[string]struct{})
	distinct := make([]string, 0, len(lines))
	for _, line := range lines {
		if _, dup := seen[line]; !dup {
			seen[line] = struct{}{}
			distinct = append(distinct, line)
		}
	}
	return distinct
}

// MapBool deduplicates using map[string]bool, the other common way to spell a
// set in existing Go code. Still two steps, but no struct{}{} and no comma-ok:
// the zero value of bool already means "absent".
func MapBool(lines []string) []string {
	seen := make(map[string]bool)
	distinct := make([]string, 0, len(lines))
	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			distinct = append(distinct, line)
		}
	}
	return distinct
}

// MapLen deduplicates by storing every line unconditionally and asking whether
// the map grew — mapset.Insert's own strategy, written by hand and with no
// dependency. Compare its timings with [Mapset] to see what the generic call
// costs.
func MapLen(lines []string) []string {
	seen := make(map[string]struct{})
	distinct := make([]string, 0, len(lines))
	for _, line := range lines {
		before := len(seen)
		seen[line] = struct{}{}
		if len(seen) != before {
			distinct = append(distinct, line)
		}
	}
	return distinct
}

// Mapset deduplicates using the same map[string]struct{} as [MapStruct], but
// asks the question with mapset.Insert, which adds the element and reports
// whether the set changed. No change to the data type is needed.
func Mapset(lines []string) []string {
	seen := make(map[string]struct{})
	distinct := make([]string, 0, len(lines))
	for _, line := range lines {
		if mapset.Insert(seen, line) {
			distinct = append(distinct, line)
		}
	}
	return distinct
}

// Set deduplicates using the proposed container/set type, whose Insert method
// delegates to the same mapset.Insert used by [Mapset].
func Set(lines []string) []string {
	seen := make(set.Set[string])
	distinct := make([]string, 0, len(lines))
	for _, line := range lines {
		if seen.Insert(line) {
			distinct = append(distinct, line)
		}
	}
	return distinct
}
