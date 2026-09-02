# unigo — uniq with a set

`unigo` reads lines from a file (or standard input) and writes each distinct
line exactly once, in order of first appearance.

It exists in five variants. They are the *same program* — byte-for-byte the
same `main_test.go`, the same structure, the same output — differing only in
how they answer "have I seen this line before?". That one question is the whole
point: it is a set membership test, and Go spells it several ways.

| variant | set type | the question | steps |
|---------|----------|--------------|-------|
| [`unigo-set`](unigo-set) | `set.Set[string]` | `if seen.Insert(line)` | one |
| [`unigo-mapset`](unigo-mapset) | `map[string]struct{}` | `if mapset.Insert(seen, line)` | one |
| [`unigo-maplen`](unigo-maplen) | `map[string]struct{}` | store, then `if len(seen) != before` | one |
| [`unigo-map`](unigo-map) | `map[string]struct{}` | `if _, dup := seen[line]; !dup` | two |
| [`unigo-mapbool`](unigo-mapbool) | `map[string]bool` | `if !seen[line]` | two |

The first three are one strategy at three levels of abstraction — behind a
method, behind a generic function, and written out by hand. `unigo-maplen` is
the control that shows what those abstractions cost: nothing measurable. The
last two look the line up before storing it, which is a different *strategy*,
and it wins on duplicate-heavy input. [`perf-lab/dedup`](../perf-lab/dedup)
measures all five.

Each variant's README covers its own spelling. The rest of this page is what
they have in common.

## Why not `uniq`?

The classic `uniq` collapses only **adjacent** runs of equal lines. Given

```
apple
banana
apple
cherry
banana
apple
```

`uniq` prints all six lines, because no two equal lines are neighbors. To make
it work you have to `sort` first — which costs O(n log n) and throws away the
original order:

```sh
sort file.txt | uniq          # sorted output
go run ./unigo/unigo-set file.txt   # first-occurrence order, one pass
```

`unigo` needs no `sort` because it remembers every line it has already seen.
That memory is a set.

Compare with [`examples/dedup_test.go`](../examples/dedup_test.go), which shows
the `slices.Sort` + `slices.Compact` recipe — the in-memory equivalent of
`sort | uniq`, with the same reordering.

## Trade-offs

* `unigo` holds every **distinct** line in memory. `sort | uniq` can spill to
  disk, so it still works on input larger than RAM. The set buys ordering and a
  single O(n) pass; it costs memory proportional to the number of distinct
  lines.
* Lines longer than 64 KiB are rejected: that is the default maximum token size
  of `bufio.Scanner`, reported as an error on exit.

## Structure

Every variant has the same three functions:

```go
func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int

func writeUnique(input io.Reader, output io.Writer) error
```

`main` is a one-line wrapper over `run`, which takes its I/O as arguments.
Returning the exit code instead of calling `os.Exit` inside `run` is what makes
the command testable in-process — `os.Exit` would terminate the test binary.
`main_test.go` drives `run` with `strings.Reader` and `bytes.Buffer`, and is
copied unchanged into every variant: the same tests pass against all five,
which is the point.

`run` handles arguments and reports errors; `writeUnique` does the work. That
split is why the variants are worth comparing at all — `writeUnique` is the
*only* function that differs between them, and its body is the snippet each
variant's README quotes. Everything above it is identical plumbing.

## Building and testing

Each variant is its own module, listed in [`go.work`](../go.work). From the
repository root:

```sh
go test ./unigo/unigo-set
go run ./unigo/unigo-set unigo/unigo-set/testdata/input.txt
```

Two quirks of this workspace are worth knowing:

* **`go build ./unigo/unigo-set` fails** with `build output "unigo-set"
  already exists and is a directory` — the binary would have the same name as
  the directory holding it. Pass `-o` with an explicit path, or just use
  `go run`.

* **`go build ./...` does not work from the repository root**, which is not
  itself a module:
  `pattern ./...: directory prefix . does not contain modules listed in go.work`.
  Name the modules instead, or `cd` into one first. This applies to the whole
  repo, not just `unigo`.

Only `unigo-set` and `unigo-mapset` depend on [`vendored`](../vendored); the
other three need nothing outside the standard library.
