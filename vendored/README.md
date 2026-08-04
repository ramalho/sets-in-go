# Vendored `container/set` and `container/mapset`

Alan Donovan's proposed set packages for the Go standard library, extracted
from Gerrit so they can be built and demoed without a patched Go toolchain.

**Neither package is in Go.** Both CLs were open (unmerged) when this was
extracted; see `PROVENANCE.txt` for the exact patchsets and date.

| Upstream | CL | Issue |
| --- | --- | --- |
| `container/set` | [745441](https://go-review.googlesource.com/c/go/+/745441) | — |
| `container/mapset` | [724420](https://go-review.googlesource.com/c/go/+/724420) | [#77052](https://github.com/golang/go/issues/77052) |
| `maps.Identical` | [760800](https://go-review.googlesource.com/c/go/+/760800) | [#78456](https://github.com/golang/go/issues/78456) |

## Re-extracting

```sh
./refresh.sh                                  # pinned patchsets
SET_PS=current MAPSET_PS=current ./refresh.sh # latest patchsets
```

The script downloads each file from Gerrit's REST API, copies `internal/fmtsort`
out of your local GOROOT, rewrites imports, and then builds and tests the result.
It overwrites `set/`, `mapset/`, `internal/fmtsort/`, and `PROVENANCE.txt`; it
does not touch `internal/maps/`.

## How much was changed

Four import lines across three files. Every function body is verbatim upstream,
so anything you show on a slide is the real proposed code.

Three imports had to be redirected because the originals are unreachable from
outside GOROOT:

- `internal/fmtsort` → vendored verbatim from your GOROOT. It is used at exactly
  one place, `mapset.String`, to sort keys the way `fmt` prints maps. It only
  imports `cmp`, `reflect`, and `slices`, so it needs no changes.
- `container/mapset` → the local `mapset/`. `set.Set` is a thin façade; nearly
  every method delegates here, which is why CL 745441 cannot be vendored alone.
- `maps` → the local `internal/maps` shim, which re-exports `Clone` and `Keys`
  from the standard library and adds `Identical`, copied from CL 760800.
  `mapset` uses `Identical` for aliasing shortcuts (`Intersection(s, s)` and
  friends). This is the one file here that is not upstream code.

## Using it from another module in this repo

The module is not published, so point at it with a `replace`:

```sh
cd examples
go mod edit -require=github.com/ramalho/sets-in-go/vendored@v0.0.0 \
            -replace=github.com/ramalho/sets-in-go/vendored=../vendored
```

## Known rough edges in the current patchsets

Worth knowing before you put these on a slide — both are artifacts of code still
under review, not of the extraction.

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

- **`set` has no tests.** CL 745441 adds only `set.go`; its one test line lives
  in `src/container/container_test.go`, which depends on scaffolding in the
  wider `container` package and is not included here. `mapset` brings its full
  258-line test suite, which passes.

## Licensing

Go's BSD-3-Clause license, copied to `LICENSE`. The `container/set` and
`container/mapset` sources are Copyright The Go Authors.
