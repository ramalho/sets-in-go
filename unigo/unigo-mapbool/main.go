// Command unigo-mapbool reads lines from a file (or standard input) and writes
// each distinct line exactly once, in order of first appearance.
//
// It is line-for-line the same program as ../unigo-map except that the set of
// lines already seen is a map[string]bool rather than a map[string]struct{}.
// A bool's zero value is false, so a missing key already reads as "not seen"
// without a comma-ok lookup.
//
// This module depends on nothing outside the standard library.
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
			fmt.Fprintf(stderr, "unigo-mapbool: %v\n", err)
			return 1
		}
		defer f.Close()
		input = f
	default:
		fmt.Fprintln(stderr, "usage: unigo-mapbool [file.txt]")
		return 1
	}

	out := bufio.NewWriter(stdout)

	seen := make(map[string]bool)
	lines := bufio.NewScanner(input)
	for lines.Scan() {
		line := lines.Text()
		if !seen[line] { // zero value false means "not seen"...
			seen[line] = true // ...then store it
			fmt.Fprintln(out, line)
		}
	}

	// cmp.Or evaluates both arguments, so the buffer is always flushed;
	// a read error wins over a write error when both happen.
	if err := cmp.Or(lines.Err(), out.Flush()); err != nil {
		fmt.Fprintf(stderr, "unigo-mapbool: %v\n", err)
		return 1
	}
	return 0
}
