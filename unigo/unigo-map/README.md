# unigo-map — the same tool, with a bare map

`unigo-map` is [`unigo-set`](../unigo-set) written against today's Go: the set
of lines already seen is a `map[string]struct{}` instead of the proposed
`container/set`. Same features, same structure, same tests — the two programs
differ in one `if` and one import.

This module depends on nothing outside the standard library.

From the repository root:

```sh
go run ./unigo/unigo-map unigo/unigo-map/testdata/input.txt

cat unigo/unigo-map/testdata/input.txt | go run ./unigo/unigo-map
```

## The one difference

With `set.Set[string]`, `Insert` adds the element *and* reports whether the set
changed, so the whole deduplication rule is one expression:

```go
seen := make(set.Set[string])
...
if seen.Insert(line) { // Insert reports whether the set changed
	fmt.Fprintln(buf, line)
}
```

With a bare map there is no such method. Testing and storing are two steps, and
the "have I seen this?" question is phrased as a comma-ok lookup whose value
half is discarded:

```go
seen := make(map[string]struct{})
...
if _, dup := seen[line]; !dup { // look it up...
	seen[line] = struct{}{} // ...then store it
	fmt.Fprintln(buf, line)
}
```

Three things the map version makes you write that the set version does not: the
`struct{}{}` zero value, the blank-ish `_, dup` lookup, and the negation. None
of it is hard — but none of it is about deduplicating lines either. See
[`FUTURE.md`](../../FUTURE.md) for where `container/set` stands.

The two versions allocate identically — `set.Set[E]` *is* `map[E]struct{}`, a
named map type rather than a wrapper — but they do **not** run at the same
speed, and neither one always wins.

`Insert` reports whether the set changed by comparing `len` around an
unconditional assignment, so it always writes. The two-step version writes only
on a miss. Measured over a million lines
([`perf-lab/dedup`](../../perf-lab/dedup)):

| distinct lines | this version | `set.Insert` | |
|----------------|--------------|--------------|---|
| 100% (no duplicates) | 59.0 ms | **56.1 ms** | −4.9% |
| 10%                  | **17.8 ms** | 20.6 ms | +15.7% |
| 1%                   | **11.8 ms** | 15.0 ms | +27.5% |

One step saves a lookup when the line is new; two steps save a write when it is
a duplicate. On duplicate-heavy input the map version is meaningfully faster —
so the argument for `set.Set` here is legibility, not speed.

Note that this is a difference of *strategy*, not of abstraction. You can write
the one-step version with a bare map too — see
[`unigo-maplen`](../unigo-maplen) — and it performs identically to `set.Insert`.

(`benchstat` medians over 10 runs of 1,000,000 lines, all `p=0.000`; see
[`perf-lab/dedup`](../../perf-lab/dedup) for the method and the full table.)

## Everything else

Behavior, structure, and trade-offs are identical across all five variants, and
[the family README](../README.md) documents them:

* why plain `uniq` does not do this, and what `sort | uniq` costs instead;
* memory proportional to the number of **distinct** lines;
* the 64 KiB `bufio.Scanner` line limit;
* `main` as a one-line wrapper over `run(args, stdin, stdout, stderr) int`, so
  the command is testable in-process without `os.Exit` killing the test binary.

`main_test.go` is a byte-for-byte copy of `unigo-set/main_test.go` — the same
tests pass against both implementations, which is the point.

## Building and testing

```sh
go test ./unigo/unigo-map                              # from the repository root
go build -o unigo/unigo-map/unigo-map ./unigo/unigo-map   # note the -o
```

`go build ./unigo/unigo-map` fails with `build output "unigo-map" already
exists and is a directory`, for the same reason it does in the other variants:
the binary would take the name of the directory holding it. And
`go build ./...` does not work from the repository root, which is not itself a
module — name the modules, or `cd` into one.

Because there is no dependency to resolve, this module needs neither the
workspace nor a `replace` directive; `cd unigo/unigo-map && GOWORK=off go test .`
works as-is.
