# unigo-set — the proposed `container/set`

`unigo-set` is the variant that uses `set.Set[string]` from the proposed
`container/set` package. For what all five variants share — why plain `uniq`
does not do this, the memory and line-length trade-offs, and the `main`/`run`
structure — see [the family README](../README.md).

From the repository root:

```sh
go run ./unigo/unigo-set unigo/unigo-set/testdata/input.txt

cat unigo/unigo-set/testdata/input.txt | go run ./unigo/unigo-set
```

## The whole program

Argument handling lives in `run`; the whole of `unigo-set` is this one
function:

```go
func writeUnique(input io.Reader, output io.Writer) error {
	buf := bufio.NewWriter(output)

	seen := make(set.Set[string])
	lines := bufio.NewScanner(input)
	for lines.Scan() {
		line := lines.Text()
		if seen.Insert(line) { // true when set changed
			fmt.Fprintln(buf, line)
		}
	}

	// Report read error if it happens; flush always
	return cmp.Or(lines.Err(), buf.Flush())
}
```

Every variant has a `writeUnique` with that exact signature, and it is the
only place they differ.

`Insert` adding the element *and* reporting whether the set changed is what
collapses "look it up, then store it" into one expression. It comes from the
proposed `container/set` package — see [`FUTURE.md`](../../FUTURE.md) and the
local copy in [`vendored/set`](../../vendored/set). With today's
`map[string]struct{}` the same loop needs the two-step:

```go
if _, dup := seen[line]; !dup {
	seen[line] = struct{}{}
	fmt.Fprintln(buf, line)
}
```

That is [`unigo-map`](../unigo-map). `set.Set.Insert` delegates to
`mapset.Insert`, so this program is also [`unigo-mapset`](../unigo-mapset) with
the call behind a method, and [`unigo-maplen`](../unigo-maplen) with the helper
written out by hand — three spellings of one strategy, and
[`perf-lab/dedup`](../../perf-lab/dedup) shows the abstractions cost nothing.

## Building and testing

```sh
go test ./unigo/unigo-set                              # from the repository root
go build -o unigo/unigo-set/unigo-set ./unigo/unigo-set   # note the -o
```

`unigo-set` is listed in [`go.work`](../../go.work), so the import of
`vendored/set` resolves with no further setup. The `replace` directive in
[`go.mod`](go.mod) means it also builds standalone, outside the workspace:

```sh
cd unigo/unigo-set && GOWORK=off go build -o /dev/null .
```

Note that the module is still named `github.com/ramalho/sets-in-go/unigo`,
from before the directory was renamed.
