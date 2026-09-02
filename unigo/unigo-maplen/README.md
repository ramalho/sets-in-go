# unigo-maplen — one step, by hand

`unigo-maplen` is [`unigo-map`](../unigo-map) with a different deduplication
rule. Instead of looking the line up and then storing it, it stores the line
unconditionally and asks whether the map grew:

```go
seen := make(map[string]struct{})
...
before := len(seen)
seen[line] = struct{}{} // just put it there...
if len(seen) > before { // ...and check map growth
	fmt.Fprintln(buf, line)
}
```

A new line makes `len` increase; a duplicate overwrites the same key and leaves
`len` alone. One hash of the key instead of two.

This module depends on nothing outside the standard library.

From the repository root:

```sh
go run ./unigo/unigo-maplen unigo/unigo-maplen/testdata/input.txt

cat unigo/unigo-maplen/testdata/input.txt | go run ./unigo/unigo-maplen
```

## Why it is worth its own directory

This is precisely what `mapset.Insert` does inside
([`vendored/mapset/mapset.go`](../../vendored/mapset/mapset.go)):

```go
func Insert[M ~map[K]V, K comparable, V bool | struct{}](x M, elem K) bool {
	pre := len(x)
	insert(x, elem) // m[k] = present
	return len(x) != pre
}
```

So `unigo-maplen` is [`unigo-mapset`](../unigo-mapset) with the helper inlined
by hand, and — since `set.Set.Insert` delegates to `mapset.Insert` — it is
[`unigo-set`](../unigo-set) too. Three programs, one strategy, three spellings.

That makes it the control in the experiment. Comparing it against the other two
answers a question the others cannot: **does the generic helper cost anything?**

It does not. Measured over 1,000,000 lines
([`perf-lab/dedup`](../../perf-lab/dedup)), the three are indistinguishable:

| distinct lines | `MapLen` (this) | `mapset.Insert` | `set.Insert` |
|----------------|-----------------|-----------------|--------------|
| 100% (no duplicates) | −5.23% | −5.50% | −4.88% |
| 10%                  | +15.95% | +16.07% | +15.72% |
| 1%                   | +27.38% | +27.20% | +27.53% |

(vs. the two-step `unigo-map` baseline, `benchstat` over 10 runs, all ±1%.)

The spread within each row is under half a percentage point. Generics, the
`~map[K]V` constraint, the method wrapper — all free. What separates the
implementations is the *strategy*, not the packaging.

## Read this next to `gen_haystack.go`

[`../../perf-lab/README.md`](../../perf-lab/README.md) records the same insight,
found the same way:

> I wrote `gen_haystack.go` after noticing that Claude's solution is not
> elegant. Instead of asking whether a number is in the map, just put it there.

`unigo-maplen` is that idea applied to lines instead of integers. It won there by
almost 30% because `gen_haystack` loops until it has *n* distinct numbers, so
nearly every store is a real insertion — the 100% row above, where this version
is 5% ahead.

The rest of the table is the other half of the story, and it is worth knowing
before reaching for the trick: on duplicate-heavy input, storing
unconditionally means writing to the map for nothing, and the two-step version
wins by 16% at 10% distinct and 27% at 1%. "Just put it there" is the right
default when most elements are new, not in general.

## The trade in one line

* **Two steps** (`unigo-map`): a lookup on every line, plus a store only on a
  miss. Cheap when most lines are duplicates.
* **One step** (this, `unigo-mapset`, `unigo-set`): a store on every line, and no
  lookup ever. Cheap when most lines are new.

Allocations are identical either way — the difference is entirely in map
operations.

## Everything else

Behavior, structure, and trade-offs match the other variants, and
[the family README](../README.md) documents them: why plain `uniq` does not do
this, memory proportional to the number of **distinct** lines, the 64 KiB
`bufio.Scanner` line limit, and `main` as a thin wrapper over
`run(cmdName, args, stdin, stdout, stderr) int` so the command is testable
in-process.
`main_test.go` is a byte-for-byte copy of `unigo-set/main_test.go`.

## Building and testing

```sh
go test ./unigo/unigo-maplen                                    # from the repository root
go build -o unigo/unigo-maplen/unigo-maplen ./unigo/unigo-maplen   # note the -o
```

`go build ./unigo/unigo-maplen` fails with `build output "unigo-maplen" already
exists and is a directory` — the binary would take the name of the directory
holding it. And `go build ./...` does not work from the repository root, which is not
itself a module; name the modules, or `cd` into one.
