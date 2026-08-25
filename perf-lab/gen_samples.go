package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

const defaultCount = 10_000_000
const needleCount = 500
const seed = 42

// generate returns a slice of count unique random numbers.
func generate(count int) []int {
	seen := make(map[int]struct{}, count)
	sample := make([]int, count)
	rng := rand.New(rand.NewSource(seed))

	for len(seen) < count {
		n := rng.Int()
		sample[len(seen)] = n
		seen[n] = struct{}{}
	}

	return sample
}

// subsample returns count numbers picked at random from numbers,
// with no repeated picks.
func subsample(numbers []int, count int) []int {
	seen := make(map[int]struct{}, count)
	sample := make([]int, count)
	rng := rand.New(rand.NewSource(seed))

	for len(seen) < count {
		i := rng.Intn(len(numbers))
		sample[len(seen)] = numbers[i]
		seen[i] = struct{}{}
	}

	return sample
}

// writeFile writes the numbers to the named file, one per line.
func writeFile(numbers []int, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, n := range numbers {
		if _, err := w.WriteString(strconv.Itoa(n)); err != nil {
			return err
		}
		if err := w.WriteByte('\n'); err != nil {
			return err
		}
	}

	return w.Flush()
}

func main() {
	count := defaultCount
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		count = n
	}

	// a single batch, so the numbers in absent.txt are distinct from those in haystack.txt
	sample := generate(count + needleCount)
	haystack := sample[:count]
	absent := sample[count:]

	files := []struct {
		numbers  []int
		filename string
	}{
		{haystack, "haystack.txt"},
		{absent, "absent.txt"},
		{subsample(haystack, needleCount), "present.txt"},
	}

	for _, f := range files {
		if err := writeFile(f.numbers, f.filename); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("wrote %d distinct integers to %s\n", len(f.numbers), f.filename)
	}
}
