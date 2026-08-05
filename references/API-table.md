# Set API comparison

Rows for math symbols, Python and `ramalho/genset` come from
[`Python-genset-set-API.csv`](Python-genset-set-API.csv).
The two Go columns document [`vendored/mapset/mapset.go`](../vendored/mapset/mapset.go)
and [`vendored/set/set.go`](../vendored/set/set.go).

In the Go columns, `x` and `y` are sets, `e` an element, and `seq` an `iter.Seq[E]`.
In `mapset` they are unnamed `map[K]struct{}` or `map[K]bool` values; in `set` they are
`Set[E]`, a named `map[E]struct{}` with methods.

| Math symbol | Python set operators | Python set methods | ramalho/genset | vendored/mapset | vendored/set |
| --- | --- | --- | --- | --- | --- |
| | | `set(it)` | | `mapset.Collect(seq)`<br>`mapset.CollectBool(seq)` | `set.Collect(seq)` |
| | `{a, b, c}` | | | `mapset.Of(a, b, c)`<br>`mapset.OfBool(a, b, c)` | `set.Of(a, b, c)` |
| | | | | | |
| S ∪ Z | `s \| z` | `s.union(it, …)` | `s.Union(z)` | `mapset.Union(x, y)` | `x.Union(y)` |
| | `s =\| z` | `s.update(it, …)` | `s.UnionUpdate(z)` | `mapset.UnionWith(x, y)` | `x.UnionWith(y)` |
| | | | `s.AddAll(e…)` | `mapset.InsertAll(x, seq)` | `x.InsertAll(seq)` |
| | | `s.add(e)` | `s.Add(e)` | `mapset.Insert(x, e)` | `x.Insert(e)` |
| | | | | | |
| S ∩ Z | `s & z` | `s.intersection(it, …)` | `s.Intersection(z)` | `mapset.Intersection(x, y)` | `x.Intersection(y)` |
| | `s &= z` | `s.intersection_update(it, …)` | `s.IntersectionUpdate(z)` | `mapset.IntersectionWith(x, y)` | `x.IntersectionWith(y)` |
| | | | | | |
| S \ Z | `s - z` | `s.difference(it, …)` | `s.Difference(z)` | `mapset.Difference(x, y)` | `x.Difference(y)` |
| | `s -= z` | `s.difference_update(it, …)` | `s.DifferenceUpdate(z)` | `mapset.DifferenceWith(x, y)` | `x.DifferenceWith(y)` |
| | | | `s.RemoveAll(e…)` | `mapset.DeleteAll(x, seq)` | `x.DeleteAll(seq)` |
| S ∆ Z | `s ^ z` | `s.symmetric_difference(it)` | `s.SymmetricDifference(z)` | `mapset.SymmetricDifference(x, y)` | `x.SymmetricDifference(y)` |
| | `s ^= z` | `s.symmetric_difference_update(it, …)` | `s.SymmetricDifferenceUpdate(z)` | `mapset.SymmetricDifferenceWith(x, y)` | `x.SymmetricDifferenceWith(y)` |
| | | | | | |
| e ∈ S | `e in s` | `s.__contains__(e)` | `s.Contains(e)` | `mapset.Contains(x, e)` | `x.Contains(e)` |
| | | | `s.ContainsAll(e…)` | `mapset.ContainsAll(x, seq)` | `x.ContainsAll(seq)` |
| | | | | | |
| S = Z | `s == z` | `s.__eq__(z)` | `s.Equal(z)` | `mapset.Equal(x, y)` | `x.Equal(y)` |
| S ⊆ Z | `s <= z` | `s.issubset(it)` | `s.SubsetOf(z)` | | |
| S ⊂ Z | `s < z` | `s.__lt__(z)` | | | |
| S ⊇ Z | `s >= z` | `s.issuperset(it)` | `s.SupersetOf(z)` | `mapset.ContainsAll(x, seq)` | `x.ContainsAll(seq)` |
| S ⊃ Z | `s > z` | `s.__gt__(z)` | | | |
| | | | | | |
| | | `s.isdisjoint(z)` | | `!mapset.Intersects(x, y)` | `!x.Intersects(y)` |
| S ∩ Z ≠ ∅ | | `not s.isdisjoint(z)` | | `mapset.Intersects(x, y)` | `x.Intersects(y)` |
| S × Z | | | | | |
| | | | | | |
| | `len(s)` | `s.__len__()` | `s.Len()` | `len(x)` | `x.Len()` |
| | `not s` | `s.__bool__()` | | `len(x) == 0` | `x.Len() == 0` |
| | | `s.__iter__()` | | `mapset.All(x)` | `x.All()` |
| | `str(s)` | `s.__repr__()` | | `mapset.String(x)` | `x.String()` |
| | | | `s.ToSlice()` | `slices.Collect(mapset.All(x))` | `slices.Collect(x.All())` |
| | | `s.clear()` | `s.Clear()` | `clear(x)` | `x.Clear()` |
| | | `s.copy()` | `s.Copy()` | `maps.Clone(x)` | `x.Clone()` |
| | | | | | |
| | | `s.discard(e)` | `s.Remove(e)` | `mapset.Delete(x, e)` | `x.Delete(e)` |
| | | `s.remove(e)` | | `mapset.Delete(x, e)` | `x.Delete(e)` |
| | | `s.pop()` | `s.Pop()` | | |
| | | | | `mapset.DeleteFunc(x, f)` | `x.DeleteFunc(f)` |

Notes:

- The source file is tab-separated despite the `.csv` extension, and several header cells
  span multiple physical lines.
- Dunder names appear in the source with U+2017 (double low line) characters; they are
  written here with ordinary underscores.
- Blank rows in the source are kept as group separators; runs of consecutive blank rows
  are collapsed to one.
- The source's `.net interface ISet<T>` column is omitted here.
- Where `mapset`/`set` provide no dedicated function, the cell shows the idiomatic Go
  expression instead: builtins (`len`, `clear`), `maps.Clone`, `slices.Collect`.
- `Insert`, `InsertAll`, `Delete`, `DeleteAll` and `DeleteFunc` return a `bool` reporting
  whether the set's length changed; Python's `add`/`discard` return `None`, and
  `s.remove(e)` raises `KeyError` if `e` is absent.
- Neither Go package has subset/proper-subset/proper-superset predicates. `ContainsAll`
  is the superset test, but it takes an `iter.Seq[E]`, not a set.
- `mapset` works on both `map[K]struct{}` and `map[K]bool`; only its constructors
  (`Of`/`OfBool`, `Collect`/`CollectBool`) need to choose between the two.
