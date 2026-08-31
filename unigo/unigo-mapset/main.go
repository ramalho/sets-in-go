// Command unigo-mapset reads lines from a file (or standard input) and writes
// each distinct line exactly once, in order of first appearance.
//
// It is the same program as ../unigo-setand ../unigo-map, in the middle position:
// the set of lines already seen is a plain map[string]struct{}, as in existing
// Go code, but the deduplication rule is one call to [mapset.Insert] rather
// than a lookup followed by a store. Adopting the helper package costs no
// change to the data type.
package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"os"

	"github.com/ramalho/sets-in-go/vendored/mapset"
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
			fmt.Fprintf(stderr, "unigo-mapset: %v\n", err)
			return 1
		}
		defer f.Close()
		input = f
	default:
		fmt.Fprintln(stderr, "usage: unigo-mapset [file.txt]")
		return 1
	}

	out := bufio.NewWriter(stdout)

	seen := make(map[string]struct{})
	lines := bufio.NewScanner(input)
	for lines.Scan() {
		line := lines.Text()
		if mapset.Insert(seen, line) { // Insert reports whether the set changed
			fmt.Fprintln(out, line)
		}
	}

	// cmp.Or evaluates both arguments, so the buffer is always flushed;
	// a read error wins over a write error when both happen.
	if err := cmp.Or(lines.Err(), out.Flush()); err != nil {
		fmt.Fprintf(stderr, "unigo-mapset: %v\n", err)
		return 1
	}
	return 0
}
