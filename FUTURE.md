# Sets in the Go standard library: where things are heading

As of Go 1.26 the standard library still has **no `Set` type and no set-algebra
functions** — the idiom is `map[E]struct{}`, and operations like intersection,
difference, and subset are written by hand (see [`examples/`](examples/)). 

A set type with full algebra is under active development
for the standard library.

## Current (Go 1.26)

The standard library offers set **primitives** — membership
(`slices.Contains`, map lookup), dedup (`slices.Sort` + `slices.Compact`),
equality (`maps.Equal`), and union (`maps.Copy`) — but no set **algebra**.


## Coming: `container/set` (targeted at Go 1.28)

In late 2025 a **Collections working group** formed (Jonathan Amsterdam, Alan
Donovan, Robert Griesemer, Daniel Martí, Roger Peppe, Keith Randall, Ian Lance
Taylor) to bring generic collection types into the standard library.

- **Umbrella proposal — generic collection types:**
  <https://github.com/golang/go/issues/80590>
  Targets the Go 1.28 milestone and defines the set operations as first-class:
  `Union`/`UnionWith`, `Intersection`/`IntersectionWith`,
  `Difference`/`DifferenceWith`, `SymmetricDifference`/`SymmetricDifferenceWith`,
  `Intersects`, `Equal`, plus `Insert`/`Delete` variants. Each algebra op comes
  in an allocation-efficient mutating `-With` form and a convenience form.

- **The `container/set` package proposal:**
  <https://github.com/golang/go/issues/69230>
  `set.Set[T]` for `comparable` elements, transparently defined as
  `map[T]struct{}` (so it interconverts with the bare map type and supports
  `len` and `range`). Alan Donovan has pinned that this is the API recommended
  by the Collections working group.

- **Original discussion (2021), by Ian Lance Taylor:**
  <https://github.com/golang/go/discussions/47331>
  An early `container/set` sketch that stalled for years.

