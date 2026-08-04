# References

* The Go Blog: [Range Over Function Types](https://go.dev/blog/range-functions) includes small example of a `Set` type.

* [`set-of-wart.md`](set-of-wart.md): why `set.Of` in the proposed `container/set`
  returns `map[E]struct{}` instead of `Set[E]`, and what that trade-off costs —
  defined types buy a method set and cost assignability.