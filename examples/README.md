# Set operations in the modern Go standard library

Runnable examples mapping mathematical set operations onto the Go
standard library (`slices` and `maps` from Go 1.21, plus the iterator
additions from Go 1.23).

The standard library has **no `Set` type and no set-algebra functions**.
The idiom is to represent a set as `map[E]struct{}` (a "map-set") and
build the operations from primitives. Each file focuses on one theme and
is written as `Example` tests, so `go test` verifies every output.

| File | Set concept | Key stdlib functions shown |
|---|---|---|
| `doc.go` | package docs + `Set[E]` / `NewSet` helper (the `map[E]struct{}` idiom) | — |
| `membership_test.go` | `x ∈ A`, and the missing `ContainsAll` (subset) done by hand | `slices.Contains`, map lookup |
| `dedup_test.go` | slice → set, incl. the adjacent-only **gotcha** | `slices.Sort` + `slices.Compact` |
| `equality_test.go` | set equality vs. sequence equality | `maps.Equal`, `slices.Equal` |
| `union_test.go` | `A ∪ B`, destructive and cloned | `maps.Copy`, `maps.Clone` |
| `algebra_test.go` | `∩`, `∖`, `△` — the algebra the stdlib omits | `maps.DeleteFunc`, `maps.Clone`, `slices.Sorted` |
| `iterators_test.go` | Go 1.23 composable pipelines | `maps.Keys`, `slices.Sorted`, `slices.Collect`, `iter.Seq` |
| `proposed_test.go` | the same operations under the **proposed** `container/set` | `set.Set`, `mapset` on legacy `map[K]bool` |

## Running

```sh
cd examples
go test ./...        # verify every example's output
go test -v ./...     # see each example run
go doc .             # read the package overview
```

## The proposed `container/set`

`proposed_test.go` reruns the themes above against Alan Donovan's proposed
`container/set` and `container/mapset` packages, so each handwritten idiom sits
next to the one-liner that would replace it.

**Neither package is in Go.** Both are open CLs, extracted into `../vendored`
and pulled in with a `replace` directive so the examples actually run:

```
require github.com/ramalho/sets-in-go/vendored v0.0.0
replace github.com/ramalho/sets-in-go/vendored => ../vendored
```

See `../vendored/README.md` for provenance, how to re-extract at a newer
patchset, and the rough edges in the current revisions.

## Notes

- The iterator examples require **Go 1.23+**, where `maps.Keys` returns
  an `iter.Seq`. Before 1.23 that lived in `golang.org/x/exp/maps` and
  returned a slice.
- Shared helpers (`union`, `sorted`, `chained`) live in the `_test`
  files where they are first used.
