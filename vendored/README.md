# Vendored `container/set`, `container/mapset` and `container/hash`

The set and map types proposed for the Go standard library, extracted from
Gerrit so they can be built and demoed without a patched Go toolchain.

**None of these packages is in Go.** Every CL was open (unmerged) when this was
extracted; see `PROVENANCE.txt` for the exact patchsets and date.

| Upstream | CL | Issue |
| --- | --- | --- |
| `container/set` | [745441](https://go-review.googlesource.com/c/go/+/745441) | [#69230](https://github.com/golang/go/issues/69230) |
| `container/mapset` | [724420](https://go-review.googlesource.com/c/go/+/724420) | [#77052](https://github.com/golang/go/issues/77052) |
| `container/hash` — `Map` | [612217](https://go-review.googlesource.com/c/go/+/612217) | [#69559](https://github.com/golang/go/issues/69559) |
| `container/hash` — `Set` | [741160](https://go-review.googlesource.com/c/go/+/741160) | [#80584](https://github.com/golang/go/issues/80584) |
| `maps.Identical` | [760800](https://go-review.googlesource.com/c/go/+/760800) | [#78456](https://github.com/golang/go/issues/78456) |
| `iter.Seq.Every`, `.Some` | [745440](https://go-review.googlesource.com/c/go/+/745440) | — |

All of it hangs off the collections umbrella issue,
[#80590](https://github.com/golang/go/issues/80590).

## Layout

Each directory mirrors one directory of the Go source tree; the second column
is the upstream path, rooted at the top of the `go` repo.

```
set/        src/container/set/     Set[E comparable] — façade over mapset
mapset/     src/container/mapset/  free functions over map[K]V
hash/       src/container/hash/    Map and Set — open hashing, custom equivalence
container/  src/container/         container_test.go alone: the conformance test
                                   checking set.Set and *hash.Set against the
                                   same abstract interfaces
internal/   (mostly not upstream)  three shims: fmtsort, maps, iter
```

`src/container/container_test.go` sits directly in the `container` directory,
a sibling of `set/` and `hash/` rather than a member of either. It declares
`package container_test` and has no non-test package to attach to: upstream it
exists purely as a conformance harness across the sibling collection packages.
Of the three shims, only `internal/fmtsort/` corresponds to a real upstream
directory (`src/internal/fmtsort/`).

## Requirements

Go **1.27** or later: `hash.Set` is parameterised by
[`maphash.Hasher[T]`](https://pkg.go.dev/hash/maphash#Hasher), which is new in
1.27 ([#70471](https://github.com/golang/go/issues/70471)). Everything here was
built and tested with `go1.27rc2`.

## Re-extracting

```sh
./refresh.sh                       # pinned patchsets
SET_PS=current HASH_PS=current ./refresh.sh   # latest patchsets
```

`SET_PS`, `MAPSET_PS`, `MAP_PS` and `HASH_PS` each select a patchset; `GO`
selects the toolchain (default `go1.27rc2`). The script downloads each file
from Gerrit's REST API, rewrites imports, and then builds and tests the result.
It overwrites `set/`, `mapset/`, `hash/`, `container/`, `internal/fmtsort/` and
`PROVENANCE.txt`; it does not touch `internal/maps/` or `internal/iter/`.

## How much was changed

Fifteen lines across eight files: twelve import lines, plus three call sites in
`hash`. Every function body is otherwise verbatim upstream, so anything you show
on a slide is the real proposed code.

Four import paths had to be redirected because the originals are unreachable
from outside GOROOT, or do not exist in any released Go:

- `internal/fmtsort` → the local `internal/fmtsort/`, taken from CL 612217
  rather than from your GOROOT: `hash` needs `fmtsort.Compare`, which that CL
  adds. It is otherwise identical to the copy in Go 1.27, and only imports
  `cmp`, `reflect` and `slices`, so it needs no changes. `mapset.String`,
  `hash.Map.String` and `hash.Set.String` all use it to sort keys the way `fmt`
  prints maps.
- `container/mapset` → the local `mapset/`. `set.Set` is a thin façade; nearly
  every method delegates here, which is why CL 745441 cannot be vendored alone.
- `maps` → the local `internal/maps` shim, which re-exports `Clone` and `Keys`
  from the standard library and adds `Identical`, copied from CL 760800.
  `mapset` uses `Identical` for aliasing shortcuts (`Intersection(s, s)` and
  friends).
- `iter` → the local `internal/iter` shim. See below.

`internal/maps` and `internal/iter` are the only two files here that are not
upstream code.

### The `iter` shim, and the three rewritten calls

`hash.Map.ContainsAll`, `hash.Set.ContainsAll` and `hash.Set.Intersects` call
`seq.Every(...)` and `seq.Some(...)` — **methods on `iter.Seq`** proposed in
CL 745440 and absent from every released Go, including 1.27.

Unlike `maps.Identical`, this one cannot be shimmed transparently. A method
must be declared in the same package as its receiver type, so `Every` can only
be added to a `Seq` that the shim itself declares — and a *defined* type
`internal/iter.Seq[V]` is not assignable to the standard `iter.Seq[V]`, which
is what `slices.Collect`, `slices.Values` and the rest of the standard library
speak. `hash/set.go` itself calls `slices.Collect(s.All())`, so a defined type
would not even compile.

So `internal/iter` declares `Seq` and `Seq2` as **aliases** for the real types
and provides `Every`/`Some` as free functions, carrying the CL's bodies
verbatim. `refresh.sh` rewrites the three call sites:

```go
keys.Every(m.Contains)          →  iter.Every(keys, m.Contains)
seq.Every(s.Contains)           →  iter.Every(seq, s.Contains)
other.All().Some(s.Contains)    →  iter.Some(other.All(), s.Contains)
```

The script greps for any remaining `.Every(`/`.Some(` afterwards and fails if a
newer patchset introduces a call the rules do not cover.

## `hash.Set` in one minute

`hash.Set[E]` is for the cases `set.Set[E comparable]` cannot reach: element
types that are not comparable, or whose `==` is the wrong equivalence relation.
You supply the relation at construction, as a `maphash.Hasher[E]`:

```go
type bigIntHasher struct{}

func (bigIntHasher) Hash(h *maphash.Hash, b *big.Int) {
	maphash.WriteComparable(h, b.Sign())
	for _, word := range b.Bits() {
		maphash.WriteComparable(h, word)
	}
}
func (bigIntHasher) Equal(x, y *big.Int) bool { return x.Cmp(y) == 0 }

s := hash.NewSet[*big.Int](bigIntHasher{})
s.Insert(big.NewInt(100))
s.Contains(big.NewInt(100)) // true — a different pointer, an equal value
```

That is `hash/example_test.go`, verbatim. Its own doc comment tells you when
*not* to use it: with a comparable `E` and ordinary `==`, use the built-in map
(and `mapset`, or `set.Set`).

Structurally it is `Map[E, struct{}]` — the same relationship `set.Set` has to
`map[E]struct{}` — implemented as open hashing: `map[uint64][]entry`, with
deleted entries tombstoned rather than compacted so that live iterators are not
disturbed.

Note the API shape differs from `set.Set`:

- `*Set[E]`, not `Set[E]`. It must be created with `NewSet`; a nil pointer and
  the zero value both panic, which `set_test.go` checks method by method.
- There is no `Of` or `Collect`. Every set starts from `NewSet(hasher)`.
- `Clone` is mandatory for copying — `Map` carries a `noCopy` field.

## Known rough edges in the current patchsets

Worth knowing before you put these on a slide. All are artifacts of code still
under review, not of the extraction.

- **`hash.Set.InsertAll` returns the wrong bool.** It is documented to report
  whether the size changed, but the body is:

  ```go
  pre := s.m.len
  for e := range seq {
      s.m.Set(e, unit{})
  }
  return pre > 0   // should be: s.m.len > pre
  ```

  It reports whether the set was *non-empty beforehand*, which is close to the
  exact inverse of the documented result:

  ```go
  s := hash.NewSet[int](maphash.ComparableHasher[int]{})
  s.InsertAll(slices.Values([]int{1, 2, 3})) // false — but the set grew 0→3
  s.InsertAll(slices.Values([]int{1, 2, 3})) // true  — but nothing changed
  ```

  `DeleteAll` right below it gets the same idiom right (`s.m.len < pre`), and
  so does `Map.SetAll` (`m.len > pre`). The bundled `set_test.go` never checks
  `InsertAll`'s result, so the suite passes.

- **`set.Of` returns the wrong type.** It is declared
  `func Of[E comparable](elems ...E) map[E]struct{}`, returning the underlying
  map type rather than `Set[E]`. Assigning to a declared variable works, as in
  the package doc example, but chaining does not:

  ```go
  var x set.Set[int] = set.Of(1, 2, 3) // ok, map[int]struct{} is assignable
  set.Of(1, 2, 3).Intersection(y)      // compile error: no method Intersection
  ```

  `Collect` on the next line returns `Set[E]` correctly, so this looks like an
  oversight rather than a design decision. Written up in full, with the
  counter-argument and why it fails, in
  [`../references/set-of-wart.md`](../references/set-of-wart.md).

- **`iter.Seq2.Every` is inverted** in CL 745440 — it initialises `every :=
  false` and sets it `true` on the first failure. `Seq.Every` and both `Some`
  methods are correct. Only `Seq.Every` and `Seq.Some` are used by `hash`, and
  only those two are carried into `internal/iter`, so this does not affect
  anything here.

- **`set` still has no tests of its own.** CL 745441 adds only
  `src/container/set/set.go`. What it does bring is
  `src/container/container_test.go`, vendored here as `container/container_test.go`:
  it is mostly a *static* conformance check that `set.Set[int]` and
  `*hash.Set[int]` both satisfy the same F-bounded `_AbstractSet` interface, plus
  generic `Subset`/`Superset`/`Take` helpers defined over it. Its only runtime
  test is `func Test(*testing.T) {}`, a placeholder.

- **`hash`'s tests are candidly labelled.** `set_test.go` opens with
  `// TODO: improve these witless LLM-generated tests.`

## Using it from another module in this repo

The module is not published, so point at it with a `replace`:

```sh
cd examples
go mod edit -require=github.com/ramalho/sets-in-go/vendored@v0.0.0 \
            -replace=github.com/ramalho/sets-in-go/vendored=../vendored
```

## Licensing

Go's BSD-3-Clause license, copied to `LICENSE`. All extracted sources are
Copyright The Go Authors.
