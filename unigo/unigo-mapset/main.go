// Command unigo-mapset reads lines from a file (or standard input) and writes
// each distinct line exactly once, in order of first appearance.
//
// It is the same program as ../unigo-set and ../unigo-map, in the middle
// position: the set of lines already seen is a plain map[string]struct{}, as
// in existing Go code, but the deduplication rule is one call to
// [mapset.Insert] rather than a lookup followed by a store. Adopting the
// helper package costs no change to the data type.
package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ramalho/sets-in-go/vendored/mapset"
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
	buf := bufio.NewWriter(output)

	seen := make(map[string]struct{})
	lines := bufio.NewScanner(input)
	for lines.Scan() {
		line := lines.Text()
		if mapset.Insert(seen, line) { // true when map changed
			fmt.Fprintln(buf, line)
		}
	}

	// Report read error if it happens; flush always
	return cmp.Or(lines.Err(), buf.Flush())
}
