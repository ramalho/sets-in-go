# unigo-mapset — the helper package, without changing the type

`unigo-mapset` is the third of three identical programs. The set of lines
already seen is a plain `map[string]struct{}`, exactly as in
[`unigo-map`](../unigo-map) — but the deduplication rule is a single call to
`mapset.Insert`, exactly as concise as [`unigo-set`](../unigo-set):

```go
seen := make(map[string]struct{})
...
if mapset.Insert(seen, line) { // Insert reports whether the set changed
	fmt.Fprintln(buf, line)
}
```

From the repository root:

```sh
go run ./unigo/unigo-mapset unigo/unigo-mapset/testdata/input.txt

cat unigo/unigo-mapset/testdata/input.txt | go run ./unigo/unigo-mapset
```

## Why this one matters

`unigo-set` says "adopt the new type". `unigo-map` shows what today's code looks
like. Neither offers anything to the codebase that already has
`map[string]struct{}` in a hundred places and is not going to change them.

That is what [`mapset`](../../vendored/mapset) is for. Its own doc comment puts it
plainly: the package "defines operations on sets represented either as
`map[K]struct{}`, or as `map[K]bool` where the boolean value is ignored" —
"common choices for representing sets in existing Go code". You keep your data
type and get the operations, one function call at a time.

The signatures are what make it work:

```go
func Insert[M ~map[K]V, K comparable, V bool | struct{}](x M, elem K) bool
```

The `~map[K]V` constraint accepts any *named* map type, and `V bool | struct{}`
accepts both spellings of "no value". So the same `mapset.Insert` serves a bare
`map[string]struct{}`, a `map[string]bool`, and your own
`type StringSet map[string]struct{}` — with no conversion.

It is also why `set.Set` can be a named map type rather than a struct wrapping
one. `set.Set.Insert` is a one-line delegation (`vendored/set/set.go:129`):

```go
func (x Set[E]) Insert(elem E) bool { return mapset.Insert(x, elem) }
```

`unigo-set` and `unigo-mapset` therefore run the same code. The difference between
them is method-call syntax and a named type — nothing more.

## Migration path

The three versions form a sequence you can walk one step at a time, with no
step forcing the next:

1. `map[string]struct{}` with hand-written lookups — [`unigo-map`](../unigo-map)
2. the same map, operations from `mapset` — this directory
3. `set.Set[string]`, operations as methods — [`unigo-set`](../unigo-set)

Step 2 is a pure addition: no type changes, no conversions, and `mapset` and
plain map syntax mix freely in the same function. Step 3 is a type change, and
buys method syntax plus the documentation value of a name that says "set".

## About speed

Do not adopt `Insert` for performance — it can cost you. It reports whether the
set changed by comparing `len` around an *unconditional* assignment, so it
writes even when the element is already present, while a hand-written
lookup-then-store skips the write on a duplicate. Over a million lines
([`perf-lab/dedup`](../../perf-lab/dedup)):

| distinct lines | hand-written two-step | `mapset.Insert` | |
|----------------|-----------------------|-----------------|---|
| 100% (no duplicates) | 59.0 ms | **55.7 ms** | −5.5% |
| 10%                  | **17.8 ms** | 20.6 ms | +16.1% |
| 1%                   | **11.8 ms** | 15.0 ms | +27.2% |

One step saves a lookup when the element is new; two steps save a write when it
is a duplicate. Allocations are identical in every case — every `B/op` and
`allocs/op` comparison is statistically indistinguishable. Adopt `mapset`
because the call says what it means — and check the duplicate ratio if the loop
is hot.

The *call* costs nothing, though. [`unigo-maplen`](../unigo-maplen) writes
`Insert`'s len-comparison out by hand, with no generics and no import, and
lands within half a percentage point of this version in every row. The generic
function is free; only the strategy it implements has a price.

(`benchstat` medians over 10 runs of 1,000,000 lines, all `p=0.000`; see
[`perf-lab/dedup`](../../perf-lab/dedup) for the method and the full table.)

## Everything else

Behavior, structure, and trade-offs match the other variants, and
[the family README](../README.md) documents them: why plain `uniq` does
not do this, memory proportional to the number of **distinct** lines, the
64 KiB `bufio.Scanner` line limit, and `main` as a one-line wrapper over
`run(args, stdin, stdout, stderr) int` so the command is testable in-process.
`main_test.go` is a byte-for-byte copy of `unigo-set/main_test.go`.

## Building and testing

```sh
go test ./unigo/unigo-mapset                                    # from the repository root
go build -o unigo/unigo-mapset/unigo-mapset ./unigo/unigo-mapset   # note the -o
```

`go build ./unigo/unigo-mapset` fails with `build output "unigo-mapset" already
exists and is a directory` — the binary would take the name of the directory
holding it. And `go build ./...` does not work from the repository root, which
is not itself a module; name the modules, or `cd` into one.
