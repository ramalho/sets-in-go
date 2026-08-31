// Command unigo-map1 reads lines from a file (or standard input) and writes
// each distinct line exactly once, in order of first appearance.
//
// It is line-for-line the same program as ../unigo-map, except for the
// deduplication rule. Instead of looking the line up and then storing it, this
// version stores it unconditionally and asks whether the map grew: a new line
// makes len increase, a duplicate leaves it alone.
//
// That is exactly what mapset.Insert does internally (see ../unigo-mapset),
// written out by hand — so this version needs nothing outside the standard
// library, and shows what the helper is buying: a name.
package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run is main with the I/O passed in, so tests can capture it.
// It returns the exit code for the process.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	input := stdin
	switch len(args) {
	case 0: // no argument: read standard input
	case 1:
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(stderr, "unigo-map1: %v\n", err)
			return 1
		}
		defer f.Close()
		input = f
	default:
		fmt.Fprintln(stderr, "usage: unigo-map1 [file.txt]")
		return 1
	}

	out := bufio.NewWriter(stdout)

	seen := make(map[string]struct{})
	lines := bufio.NewScanner(input)
	for lines.Scan() {
		line := lines.Text()
		before := len(seen)
		seen[line] = struct{}{}  // just put it there...
		if len(seen) != before { // ...and see whether that added anything
			fmt.Fprintln(out, line)
		}
	}

	// cmp.Or evaluates both arguments, so the buffer is always flushed;
	// a read error wins over a write error when both happen.
	if err := cmp.Or(lines.Err(), out.Flush()); err != nil {
		fmt.Fprintf(stderr, "unigo-map1: %v\n", err)
		return 1
	}
	return 0
}
