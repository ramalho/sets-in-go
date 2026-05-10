package main

import (
	"fmt"
	"os"
	"strings"

	mapset "github.com/ramalho/sets-tools"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: listinterfaces <file.go>\n")
		os.Exit(1)
	}
	interfaces, err := mapset.ListInterfaces(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(strings.Join(interfaces, "$"))
}
