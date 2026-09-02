# unigo-mapbool — the same tool, with a `bool` map

`unigo-mapbool` is [`unigo-map`](../unigo-map) with one change: the set of
lines already seen is a `map[string]bool` instead of a `map[string]struct{}`.
Same features, same structure, same tests — the two programs differ in one
type and one `if`.

This module depends on nothing outside the standard library.

From the repository root:

```sh
go run ./unigo/unigo-mapbool unigo/unigo-mapbool/testdata/input.txt

cat unigo/unigo-mapbool/testdata/input.txt | go run ./unigo/unigo-mapbool
```

## The one difference

`unigo-map` stores presence with the empty struct, so seeing a line still
needs a comma-ok lookup — the map only ever holds `true`-shaped values, and
`_, dup := seen[line]` throws the value half away:

```go
seen := make(map[string]struct{})
...
if _, dup := seen[line]; !dup { // look it up...
	seen[line] = struct{}{} // ...then store it
	fmt.Fprintln(out, line)
}
```

A `bool` map's zero value is `false`, so a missing key already reads as "not
seen" — no comma-ok, no throwaway blank identifier:

```go
seen := make(map[string]bool)
...
if !seen[line] { // zero value false means "not seen"...
	seen[line] = true // ...then store it
	fmt.Fprintln(out, line)
}
```

The two maps do **not** allocate identically: `bool` costs a byte per entry
that `struct{}` does not, since `struct{}` occupies zero bytes. For a set
whose values are never inspected — only tested for presence — `struct{}` is
the idiomatic choice for that reason; `unigo-mapbool` trades that byte for a
lookup that reads slightly more plainly.

## Everything else

Behavior, structure, and trade-offs are identical to `unigo-map`, and
[the family README](../README.md) documents them:

* why plain `uniq` does not do this, and what `sort | uniq` costs instead;
* memory proportional to the number of **distinct** lines;
* the 64 KiB `bufio.Scanner` line limit;
* `main` as a one-line wrapper over `run(args, stdin, stdout, stderr) int`, so
  the command is testable in-process without `os.Exit` killing the test binary.

`main_test.go` is a byte-for-byte copy of `unigo-map/main_test.go` — the same
tests pass against both implementations, which is the point.

## Building and testing

```sh
go test ./unigo/unigo-mapbool                              # from the repository root
go build -o unigo/unigo-mapbool/unigo-mapbool ./unigo/unigo-mapbool   # note the -o
```

`go build ./unigo/unigo-mapbool` fails with `build output "unigo-mapbool"
already exists and is a directory`, for the same reason it does in
`unigo-map`: the binary would take the name of the directory holding it.

Because there is no dependency to resolve, this module needs neither the
workspace nor a `replace` directive; `cd unigo/unigo-mapbool && GOWORK=off go
test .` works as-is.
