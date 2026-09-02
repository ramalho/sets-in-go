// Command unigo-maplen reads lines from a file (or standard input) and writes
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
	"path/filepath"
)

func main() {
	cmdName := filepath.Base(os.Args[0])
	os.Exit(run(cmdName, os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run is main with the I/O passed in, so tests can capture it.
// It returns the exit code for the process.
func run(cmdName string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	input := stdin
	switch len(args) {
	case 0: // no argument: read standard input
	case 1:
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(stderr, "%s: %v\n", cmdName, err)
			return 1
		}
		defer f.Close()
		input = f
	default:
		fmt.Fprintf(stderr, "usage: %s [file.txt]\n", cmdName)
		return 1
	}

	if err := writeUnique(input, stdout); err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", cmdName, err)
		return 1
	}
	return 0
}

// writeUnique copies input to output,
// keeping only the first occurrence of each line.
func writeUnique(input io.Reader, output io.Writer) error {
	lines := bufio.NewScanner(input)
	buf := bufio.NewWriter(output)

	seen := make(map[string]struct{})
	for lines.Scan() {
		line := lines.Text()
		before := len(seen)
		seen[line] = struct{}{} // just put it there...
		if len(seen) > before { // ...and check map growth
			fmt.Fprintln(buf, line)
		}
	}

	// Report read error if it happens; flush always
	return cmp.Or(lines.Err(), buf.Flush())
}
