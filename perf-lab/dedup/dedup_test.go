package dedup

import (
	"math/rand"
	"slices"
	"strconv"
	"testing"
)

const seed = 42 // the same constant the other perf-lab scripts use

// implementations are the four spellings of the deduplication rule.
var implementations = []struct {
	name string
	fn   func([]string) []string
}{
	{"MapStruct", MapStruct},
	{"MapBool", MapBool},
	{"MapLen", MapLen},
	{"Mapset", Mapset},
	{"Set", Set},
}

// makeLines returns count lines drawn from distinct different values, in an
// order fixed by the seed. With distinct == count there are no duplicates;
// with distinct == count/10 each value appears about ten times.
func makeLines(count, distinct int) []string {
	rng := rand.New(rand.NewSource(seed))
	lines := make([]string, count)
	for i := range lines {
		lines[i] = strconv.Itoa(rng.Intn(distinct))
	}
	return lines
}

func TestImplementationsAgree(t *testing.T) {
	lines := makeLines(10_000, 1_000)

	want := MapStruct(lines)
	if len(want) != 1_000 {
		t.Fatalf("got %d distinct lines, want 1000", len(want))
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			if got := impl.fn(lines); !slices.Equal(got, want) {
				t.Errorf("%s disagrees with MapStruct", impl.name)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			if got := impl.fn(nil); len(got) != 0 {
				t.Errorf("got %v, want no lines", got)
			}
		})
	}
}

// BenchmarkDedup runs each implementation over the same input, at two
// duplicate ratios: every line distinct, and each line repeated ten times.
func BenchmarkDedup(b *testing.B) {
	const count = 1_000_000

	// The duplicate ratio is what separates the implementations: see the
	// package doc. All distinct means every Insert really inserts.
	inputs := []struct {
		name  string
		lines []string
	}{
		{"100%", makeLines(count, count)},
		{"10%", makeLines(count, count/10)},
		{"1%", makeLines(count, count/100)},
	}

	// The key=value sub-benchmark names let benchstat project on them:
	//	benchstat -col /impl -row /distinct bench.txt
	for _, input := range inputs {
		for _, impl := range implementations {
			b.Run("distinct="+input.name+"/impl="+impl.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					impl.fn(input.lines)
				}
			})
		}
	}
}
