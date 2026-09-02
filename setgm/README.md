# `setgm` — the proposed `Set` rewritten with generic methods

An alternative implementation of Alan Donovan's proposed `container/set.Set`,
using the **method-level type parameters** added in Go 1.27.

## The gap it closes

`container/mapset`'s binary operations are heterogeneous by design:

```go
func Union[MX ~map[K]VX, MY ~map[K]VY, K comparable, VX, VY bool | struct{}](x MX, y MY) MX
```

A modern `map[K]struct{}` can be unioned with a legacy `map[K]bool`. That is
the package's whole reason for existing. But every `set.Set` method throws the
generality away:

```go
func (x Set[E]) Union(y Set[E]) Set[E]
```

Not because the design wants it — because **before Go 1.27 a method could not
declare type parameters of its own**. The second operand had to be spelled in
terms of the receiver's `E` alone, so it could only ever be another `Set[E]`.
`../examples/superset_test.go` documents the same wall from the caller's side:
"`Equal` and `IntersectionWith` cannot do this: their signatures bind both
operands to a single map type `M`."

`SetGM` restores it:

```go
func (x SetGM[E]) Union[M MapSet[E, V], V bool | struct{}](y M) SetGM[E]
```

```go
perms := setgm.Of("read", "write")
granted := map[string]bool{"read": true, "admin": true} // legacy set

perms.Union(granted)      // {admin, read, write}
perms.Subset(granted)     // false
perms.Intersects(granted) // true
```

No conversion at the call site, and no explicit type arguments — inference
resolves `M` and `V` from the operand.

## Which methods are generic

| Generic in the operand | Non-generic |
|---|---|
| `Equal`, `Superset`, `ProperSuperset`, `Subset`, `ProperSubset`, `Intersects` | `Len`, `Contains`, `All`, `ContainsAll`, `String`, `Clone`, `Clear` |
| `Union`, `Intersection`, `Difference`, `SymmetricDifference` | `Insert`, `InsertAll`, `Delete`, `DeleteAll`, `DeleteFunc` |
| `UnionWith`, `IntersectionWith`, `DifferenceWith`, `SymmetricDifferenceWith` | |

The operand constraint is `MapSet[E, V]`, i.e. `~map[E]V` for `V` in
`bool | struct{}`. It accepts `SetGM[E]`, `set.Set[E]`, `map[E]struct{}`, and a
legacy `map[E]bool`.

Where `mapset` is already general enough (`Union`, `Intersection`, `Difference`,
`SymmetricDifference`, `UnionWith`) the method delegates. Where `mapset` binds
both operands to one type `M` (`Equal`, `Intersects`, and the other three
`*With`), the method is implemented here.

## Deliberate departures from the proposal

- **`Of` returns `SetGM[E]`.** The proposal's `set.Of` is declared as returning
  `map[E]struct{}` although its body returns `Set[E]` — which is why the
  examples in `../examples` pre-declare their variables.
- **`Superset`/`Subset` and proper variants are provided.** The proposal has no
  such predicate; `ContainsAll` *is* the superset test, but it takes an
  `iter.Seq[E]`. These also carry the identity fast path `ContainsAll` lacks.
- **`From`** converts any map-shaped set into a `SetGM[E]`, which a plain
  conversion cannot do — `SetGM[E](m)` is legal only when `m` is already a
  `map[E]struct{}`.

## What generic methods cost

Three limits, all verified against `go1.27.0`:

1. **Interfaces cannot declare them** — `interface method must have no type
   parameters`. Every generic method here is invisible to interface
   satisfaction, so `SetGM` cannot be reached through an interface exposing
   `Union`, `Equal`, `Superset`, … The non-generic methods are unaffected, so
   `SetGM` still satisfies `fmt.Stringer`.
2. **No bare method values** — `cannot use generic function x.Union without
   instantiation`. `x.Union[map[int]bool]` is an ordinary func value; `x.Union`
   is not.
3. **`identical` is duplicated.** The identity fast paths need
   `maps.Identical` (still-open CL 760800). `../vendored/internal/maps` supplies
   it, but an `internal` package cannot be imported across module boundaries,
   so the three lines are repeated in `setgm.go`.

## Requirements

Go 1.27. The `go 1.27` line in `go.mod` is enough: with the default
`GOTOOLCHAIN=auto`, an older `go` downloads and switches to a 1.27 toolchain on
its own.

**Your `gofmt` and `gopls` must also be 1.27+.** Go 1.26's parser rejects
generic methods outright (`method must have no type parameters`), so an editor
running an older toolchain will red-underline the whole file and `gofmt -l`
will report false positives. Use `$(go env GOROOT)/bin/gofmt`, which resolves
to the switched-to toolchain rather than whatever `go` binary is first in PATH.

```sh
cd setgm
go test ./...      # run the tests and verify every Example's output
go doc .           # read the package overview
```
