package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
)

const defaultCount = 150_000
const seed = 42

func main() {
	count := defaultCount
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		count = n
	}

	f, err := os.Create("haystack.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	seen := make([]int, 0, count)
	rng := rand.New(rand.NewSource(seed))

	for len(seen) < count {
		n := rng.Int()
		if slices.Contains(seen, n) {
			continue
		}
		seen = append(seen, n)
		if _, err := w.WriteString(strconv.Itoa(n)); err != nil {
			log.Fatal(err)
		}
		if err := w.WriteByte('\n'); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("wrote %d distinct integers to haystack.txt\n", count)
}
