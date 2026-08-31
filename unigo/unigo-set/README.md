# unigo — uniq with a set

`unigo` reads lines from a file (or standard input) and writes each distinct
line exactly once, in order of first appearance.

From the repository root:

```sh
go run ./unigo-setunigo/testdata/input.txt

cat unigo/testdata/input.txt | go run ./unigo
```

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
sort file.txt | uniq     # sorted output
go run ./unigo-setfile.txt   # first-occurrence order, one pass
```

`unigo` needs no `sort` because it remembers every line it has already seen.
That memory is a set.

## The whole program

Stripped of I/O, `unigo` is this:

```go
seen := make(set.Set[string])
lines := bufio.NewScanner(input)
for lines.Scan() {
	line := lines.Text()
	if seen.Insert(line) { // Insert reports whether the set changed
		fmt.Fprintln(out, line)
	}
}
```

`Insert` adding the element *and* reporting whether the set changed is what
collapses "look it up, then store it" into one expression. It comes from the
proposed `container/set` package — see [`FUTURE.md`](../FUTURE.md) and the
local copy in [`vendored/set`](../vendored/set). With today's
`map[string]struct{}` the same loop needs the two-step:

```go
if _, dup := seen[line]; !dup {
	seen[line] = struct{}{}
	fmt.Fprintln(out, line)
}
```

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

`main` is a one-line wrapper over `run`, which takes its I/O as arguments:

```go
func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int
```

Returning the exit code instead of calling `os.Exit` inside `run` is what makes
the command testable in-process — `os.Exit` would terminate the test binary.
`main_test.go` drives `run` with `strings.Reader` and `bytes.Buffer`.

## Building and testing

```sh
go test ./unigo-set                # from the repository root
go build -o unigo/unigo-set./unigo-set  # note the -o
```

Two quirks of this workspace are worth knowing:

* **`go build ./unigo-set fails** with `build output "unigo" already exists and is a
  directory` — the binary would have the same name as the directory holding it.
  Pass `-o` with an explicit path, or just use `go run ./unigo-set.

* **`go build ./...` does not work from the repository root**, which is not
  itself a module:
  `pattern ./...: directory prefix . does not contain modules listed in go.work`.
  Name the modules instead — `go build ./unigo-set./examples/... ./vendored/...` —
  or `cd` into one first. This applies to the whole repo, not just `unigo`.

`unigo` is listed in [`go.work`](../go.work), so the import of `vendored/set`
resolves with no further setup. The `replace` directive in
[`go.mod`](go.mod) means it also builds standalone, outside the workspace
(`cd unigo && GOWORK=off go build -o /dev/null .`).
