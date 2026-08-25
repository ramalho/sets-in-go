package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
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

	seen := make([]uint64, 0, count)
	rng := rand.New(rand.NewSource(seed))

	for len(seen) < count {
		n := rng.Uint64()
		dup := false
		for _, x := range seen {
			if x == n {
				dup = true
				break
			}
		}
		if dup {
			continue
		}
		seen = append(seen, n)
		if _, err := w.WriteString(strconv.FormatUint(n, 10)); err != nil {
			log.Fatal(err)
		}
		if err := w.WriteByte('\n'); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("wrote %d distinct uint64 numbers to haystack.txt\n", count)
}
