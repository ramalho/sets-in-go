// Command unigo-map reads lines from a file (or standard input) and writes
// each distinct line exactly once, in order of first appearance.
//
// It is line-for-line the same program as ../unigo-set except that the set of
// lines already seen is the idiomatic map[string]struct{} rather than the
// proposed container/set. Only the deduplication rule below differs: with a
// bare map there is no Insert reporting whether the set changed, so the test
// and the store are two separate steps.
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
			fmt.Fprintf(stderr, "unigo-map: %v\n", err)
			return 1
		}
		defer f.Close()
		input = f
	default:
		fmt.Fprintln(stderr, "usage: unigo-map [file.txt]")
		return 1
	}

	if err := writeUnique(input, stdout); err != nil {
		fmt.Fprintf(stderr, "unigo-map: %v\n", err)
		return 1
	}
	return 0
}

// writeUnique copies input to output, keeping only the first occurrence of
// each line. It is the whole program; everything else is plumbing.
func writeUnique(input io.Reader, output io.Writer) error {
	buf := bufio.NewWriter(output)

	seen := make(map[string]struct{})
	lines := bufio.NewScanner(input)
	for lines.Scan() {
		line := lines.Text()
		if _, dup := seen[line]; !dup { // look it up...
			seen[line] = struct{}{} // ...then store it
			fmt.Fprintln(buf, line)
		}
	}

	// Flush always; report read error if it happens
	return cmp.Or(lines.Err(), buf.Flush())
}
