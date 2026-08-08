# Where the vendored code will live in the Go source tree

`vendored/` holds code extracted from unmerged Gerrit CLs (see
[`../vendored/PROVENANCE.txt`](../vendored/PROVENANCE.txt)). This is the
future layout of the same code inside the Go repository, under `go/src/`,
plus the other packages named in the umbrella proposal
[ggo#80590 — *proposal: container/...: generic collection types*](https://github.com/golang/go/issues/80590).

Legend:

| Mark | Meaning |
| --- | --- |
| `★` | vendored in this repo under `vendored/` |
| `✚` | new file/package proposed for Go 1.28, CL uploaded |
| `?` | package proposed but no CL on Gerrit yet, so file names are a guess |
| `=` | already in the standard library, unchanged by the proposal |
| (none) | a directory shown only to give the path to a marked file |

## Simplified proposed layout

Test files omitted.

```
go/src/
├── container/
│   ├── hash/
│   │   ├── map.go           ✚  hash.Map, keyed by maphash.Hasher (#69559)
│   │   └── set.go           ✚  hash.Set, same approach (#80584)
│   ├── heap/
│   │   ├── heap.go          =  existing, pre-generics
│   │   └── v2/
│   │       └── heap.go      ?  generic Heap replacing the callback API (#77397)
│   ├── list/
│   │   └── list.go          =  existing
│   ├── mapset/
│   │   └── mapset.go        ★✚ helpers for legacy map-based sets (#77052)
│   ├── ordered/
│   │   └── ordered.go       ?  tree-based ordered.Map/Set (#60630)
│   ├── ring/
│   │   └── ring.go          =  existing
│   └── set/
│       └── set.go           ★✚ the canonical set, transparently map[T]struct{} (#69230)
├── hash/
│   └── maphash/
│       └── maphash.go       =  Hasher, released in go1.27 (#70471)
├── internal/
│   └── fmtsort/
│       └── sort.go          ★  used by mapset.String to sort keys the way fmt does
└── maps/
    └── maps.go              ★  gains Identical (#78456), used by mapset for aliasing shortcuts
```

## Proposed layout of `go/src`

```
go/src/
├── container/
│   ├── container_test.go                ✚  CL 761460 — unexported abstract
│   │                                       _AbstractCollection/_Set/_Map constraints
│   │                                       + symmetry tests; also touched by
│   │                                       CLs 745441, 612217, 741160
│   ├── hash/                            ✚  hash.Map / hash.Set, keyed by
│   │                                       maphash.Hasher (#69559, #80584)
│   │   ├── example_test.go              ✚  CL 612217, CL 741160
│   │   ├── iter_test.go                 ✚  CL 612217
│   │   ├── map.go                       ✚  CL 612217
│   │   ├── map_test.go                  ✚  CL 612217
│   │   ├── set.go                       ✚  CL 741160
│   │   └── set_test.go                  ✚  CL 741160
│   ├── heap/                            =  existing, pre-generics
│   │   ├── example_intheap_test.go      =
│   │   ├── example_pq_test.go           =
│   │   ├── heap.go                      =
│   │   ├── heap_test.go                 =
│   │   └── v2/                          ?  generic Heap replacing the
│   │       └── heap.go                  ?  callback API (#77397)
│   ├── list/                            =  existing
│   │   ├── example_test.go              =
│   │   ├── list.go                      =
│   │   └── list_test.go                 =
│   ├── mapset/                          ✚  package-level helpers for legacy
│   │                                       map[T]bool / map[T]struct{} sets (#77052)
│   │   ├── mapset.go                    ★✚ CL 724420 ps 27
│   │   └── mapset_test.go               ★✚ CL 724420 ps 27
│   ├── ordered/                         ?  tree-based ordered.Map/Set
│   │   ├── ordered.go                   ?  (#60630)
│   │   └── ordered_test.go              ?
│   ├── ring/                            =  existing
│   │   ├── example_test.go              =
│   │   ├── ring.go                      =
│   │   └── ring_test.go                 =
│   └── set/                             ✚  the canonical set type, transparently
│                                           map[T]struct{} (#69230)
│       └── set.go                       ★✚ CL 745441 ps 13 — no test file; its only
│                                           test line lives in container_test.go
├── hash/
│   └── maphash/
│       └── maphash.go                   =  Hasher — the custom hash / equivalence
│                                           interface that container/hash builds on;
│                                           #70471, CL 657296, released in go1.27
├── internal/
│   └── fmtsort/
│       └── sort.go                      ★  existing; used by mapset.String to sort
│                                           keys the way fmt does. CL 612217
│                                           modifies it.
└── maps/
    └── maps.go                          ★  existing; CL 760800 adds Identical
                                            (#78456), which mapset uses for aliasing
                                            shortcuts
```

CL 760800 also carries the usual API-promise and release-note files, which the
`set` and `mapset` CLs do not have yet:

```
go/
├── api/
│   └── next/
│       └── 78456.txt                 ✚  CL 760800 — maps.Identical
└── doc/
    └── next/
        └── 6-stdlib/
            └── 99-minor/
                └── maps/
                    └── 78456.md      ✚  CL 760800
```

## Mapping from `vendored/`

```
vendored/                    go/src/
├── set/set.go                ───▶ container/set/set.go
├── mapset/mapset.go          ───▶ container/mapset/mapset.go
├── mapset/mapset_test.go     ───▶ container/mapset/mapset_test.go
├── internal/fmtsort/sort.go  ───▶ internal/fmtsort/sort.go
└── internal/maps/maps.go     ───▶ maps/maps.go (partly: only
                                   Identical, from CL 760800;
                                   Clone and Keys are shims)
```

`vendored/internal/maps` is the one file in `vendored/` that is not upstream
code. Inside GOROOT no shim is needed: `container/mapset` imports `maps`,
`internal/fmtsort` and `container/mapset` directly, so the four rewritten
import lines listed in [`../vendored/README.md`](../vendored/README.md) all
revert to their standard paths.

## Notes

- Every `✚` CL was still open when this was written; nothing here is in a
  released Go, and file names for `?` packages may change.
- `container/ordered` (#60630) and `container/heap/v2` (#77397) are named in
  the umbrella proposal but have no CL on Gerrit as of 2026-08-08, so their
  directory names come from the proposal and their file names are conventional
  guesses.
- The proposal also mentions insertion-ordered hash maps (#80194) and stacks as
  likely future additions; no package path has been proposed for those.
- The working group prefers the word *collection* over *container*, but the new
  packages still go in the existing `container` tree.
