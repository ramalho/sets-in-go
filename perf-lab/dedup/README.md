# Deduplication: four spellings, measured

[`dedup.go`](dedup.go) contains the deduplication rule from the `unigo`
commands, written four ways with the I/O stripped out:

| function    | set type              | the question                        | steps | from |
|-------------|-----------------------|-------------------------------------|-------|------|
| `MapStruct` | `map[string]struct{}` | `if _, dup := seen[line]; !dup`     | two | [`../../unigo-map`](../../unigo-map) |
| `MapBool`   | `map[string]bool`     | `if !seen[line]`                    | two | the other legacy idiom |
| `MapLen`    | `map[string]struct{}` | store, then `if len(seen) != before` | one | [`../../unigo-map1`](../../unigo-map1) |
| `Mapset`    | `map[string]struct{}` | `if mapset.Insert(seen, line)`      | one | [`../../unigo-mapset`](../../unigo-mapset) |
| `Set`       | `set.Set[string]`     | `if seen.Insert(line)`              | one | [`../../unigo-set](../../unigo-set |

The last three are the same strategy at three levels of abstraction: written
out by hand, behind a generic function, and behind a method. `MapLen` is the
control that shows what those abstractions cost.

```sh
go test ./perf-lab/dedup                                      # they all agree

go test ./perf-lab/dedup -run='^$' -bench=Dedup -benchmem -count=10 > bench.txt
benchstat -col /impl -row /distinct bench.txt
```

The sub-benchmark names are in `key=value` form so `benchstat` can project on
them. Install it with
`go install golang.org/x/perf/cmd/benchstat@latest`.

## Result: the ranking flips

1,000,000 lines, `-count=10`, Go 1.27, Apple M2 Max. `MapStruct` is the
baseline; `~` means the difference is not statistically significant.

```
        │  MapStruct  │              MapBool               │               MapLen                │               Mapset                │                 Set                 │
        │   sec/op    │   sec/op     vs base               │   sec/op     vs base                │   sec/op     vs base                │   sec/op     vs base                │
100%      58.95m ± 2%   59.52m ± 1%       ~ (p=0.063 n=10)   55.87m ± 1%   -5.23% (p=0.000 n=10)   55.71m ± 1%   -5.50% (p=0.000 n=10)   56.08m ± 1%   -4.88% (p=0.000 n=10)
10%       17.78m ± 0%   19.00m ± 1%  +6.86% (p=0.000 n=10)   20.61m ± 1%  +15.95% (p=0.000 n=10)   20.63m ± 1%  +16.07% (p=0.000 n=10)   20.57m ± 1%  +15.72% (p=0.000 n=10)
1%        11.77m ± 1%   12.72m ± 0%  +8.09% (p=0.000 n=10)   14.99m ± 0%  +27.38% (p=0.000 n=10)   14.97m ± 0%  +27.20% (p=0.000 n=10)   15.01m ± 0%  +27.53% (p=0.000 n=10)
geomean   23.11m        24.32m       +5.26%                  25.85m       +11.86%                  25.82m       +11.74%                  25.87m       +11.97%
```

Every `B/op` and `allocs/op` comparison came back `~`, and the allocation
counts are `all samples are equal`: 4,118 / 531 / 80 allocations per operation,
identical across all five. `set.Set[E]` *is* `map[E]struct{}`; there is no
wrapper and no extra allocation. Whatever the time difference is, it is not
memory.

## The abstractions are free

`MapLen`, `Mapset`, and `Set` are the same strategy written three ways — by
hand, as a generic function, as a method. They land within half a percentage
point of each other in every row:

| distinct lines | `MapLen` | `Mapset` | `Set` | spread |
|----------------|----------|----------|-------|--------|
| 100% | −5.23% | −5.50% | −4.88% | 0.62 pp |
| 10%  | +15.95% | +16.07% | +15.72% | 0.35 pp |
| 1%   | +27.38% | +27.20% | +27.53% | 0.33 pp |

Generics, the `~map[K]V` constraint, and the method wrapper all compile away.
Choose among those three on legibility; the machine cannot tell them apart.

## Why the one-step versions behave differently

They store unconditionally and detect the insertion by watching `len`.
`MapLen` writes it out:

```go
before := len(seen)
seen[line] = struct{}{}
if len(seen) != before { ... }
```

and `mapset.Insert` is the same thing with a name on it:

```go
func Insert[M ~map[K]V, K comparable, V bool | struct{}](x M, elem K) bool {
	pre := len(x)
	insert(x, elem) // m[k] = present — always
	return len(x) != pre
}
```

So each spelling does a different amount of work per line:

* **new line** — the two-step versions hash it twice, once to look it up and
  once to store it. The one-step versions hash it once and store. One step wins.
* **duplicate line** — the two-step versions hash it once and take a read-only
  path. The one-step versions hash it once and *write anyway*, dirtying a
  bucket to no effect. Two steps win.

That is the whole story: one step trades a wasted write on duplicates for a
saved lookup on new elements. The more duplicates in the input, the worse the
trade — the penalty grows from ~16% at 10% distinct to ~27% at 1% distinct, and
inverts to a ~5% *win* when nothing repeats.

`MapBool` is a smaller, separate effect, and note *where* it shows up: with no
duplicates it is not clearly distinguishable from `MapStruct` (`p=0.063`, just
above the 0.05 threshold), but it costs 7–8% once duplicates dominate. That is
the read path, not the write path. A likely explanation — not verified here —
is that `_, dup := seen[line]` on a `struct{}` map only has to answer "is the
key present?", while `!seen[line]` on a `bool` map must also load the value
byte. Confirming that would mean reading the generated code.

## The connection to `gen_haystack.go`

[`../README.md`](../README.md) records that replacing "ask whether the number is
in the map" with "just put it there" made `gen_haystack.go` nearly 30% faster.
That is this same trade — and it won there because `gen_haystack` loops until it
has *n* distinct numbers, so almost every insertion is new. It is the top row of
the table, where one-step is the right call.

The lesson is not "`Insert` is slow". It is that "insert and tell me if it
changed" and "check, then insert" are different operations with different
costs, and the input decides. For `unigo` on ordinary text, where most lines
tend to be distinct, `Insert` is the right default — and it is the one that
reads best.

## Caveats

* One machine, one Go version, in-memory slices. Re-run before quoting. The
  numbers above are `benchstat` medians over 10 runs, with the ± showing the
  95% confidence interval.
* Benchmarking the dedup rule alone exaggerates its share: the real `unigo`
  commands also scan and write, so these differences shrink against actual I/O.
* The four implementations are re-typed here rather than imported, because each
  command is its own `package main`. [`TestImplementationsAgree`](dedup_test.go)
  checks all four still produce identical output; it does not check that they
  still match the code in the command directories.
