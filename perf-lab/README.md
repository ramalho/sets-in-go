# Performance experiments

## Generate haystack

These scripts generate a file named `haystack.txt` with unique
random `uint64` in random order, from a constant seed.

* `gen_haystack_cl.go` was written by Claude Sonnet 5 with high effort on Aug. 24 2026. 
It uses a map to keep track of previously generated numbers.
The run time is approximately O(n).

* I wrote `gen_haystack.go` after noticing that Claude's solution is not elegant.
Instead of asking whether a number is in the map, just put it there.
The effect is the same, with less code and less work for the program.
The run time is approximately O(n). In practice, it is nearly 30% faster.

* I asked Claude to write `gen_haystack_slice.go` for comparison, to show how slow it is to use
a slice instead of a map to keep track of previous numbers. The run time is approximately O(n²).