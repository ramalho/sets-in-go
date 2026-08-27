# Performance experiments

## Deduplication: four spellings compared

[`dedup/`](dedup/) benchmarks the deduplication rule from the `unigo` commands
written four ways — two hand-written map idioms, `mapset.Insert`, and
`set.Set.Insert`. The ranking flips with the duplicate ratio, for the same
reason `gen_haystack.go` below got faster: `Insert` trades a wasted write on
duplicates for a saved lookup on new elements. See [`dedup/README.md`](dedup/README.md).

## Generate haystack

These scripts generate a file named `haystack.txt` with unique
random `int` in random order, from a constant seed.

* `gen_haystack_cl.go` was written by Claude Sonnet 5 with high effort on Aug. 24 2026. 
It uses a map to keep track of previously generated numbers.
The run time is approximately O(n).

* I wrote `gen_haystack.go` after noticing that Claude's solution is not elegant.
Instead of asking whether a number is in the map, just put it there.
The effect is the same, with less code and less work for the program.
The run time is approximately O(n). In practice, it is nearly 30% faster.
This version does not preserve the order of the generated numbers;
it outputs the numbers as ordered in the map.

* I asked Claude to write `gen_haystack_slice.go` for comparison, to show how slow it is to use
a slice instead of a map to keep track of previous numbers. The run time is approximately O(n²).

* I wrote `gen_haystack_ordered.go` to preserve the ordering of the random numbers.
It pre-allocates a slice, and places each random number at a position
given by the current `len` of the map.
This replaces an explicit `if` with an assignment that may be overwritten later
when the `len` does not change (which means the overwtitten number was not unique).
The run time is approximately O(n). In practice, it is the fastest of these scripts.
