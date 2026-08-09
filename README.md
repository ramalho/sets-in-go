# Sets in Modern Go

Exploring the design space for set types in modern Go.

Materials for a talk.

## Title and abstract

As approved by the organizers of GopherCon UK 2026

### Sets in Modern Go

Luciano Ramalho, master explainer

Set theory is closely related to Boolean logic and databases.
In practice, set operations can dramatically simplify code,
replacing complex nested loops and ifs with declarative expressions
that are easier to read and may offer better performance.
That's why most popular, modern languages provide sets in
their standard libraries, but not yet Go.

In this talk I show how set operations simplify common programming
tasks, and how to implement a generic set type in modern Go.


## References

### Posts

* [Generics in action](https://bitfieldconsulting.com/posts/generic-set): a set type in Go

* [Writing generic collection types in Go](https://www.dolthub.com/blog/2024-07-01-golang-generic-collections/): the missing documentation

### Packages

* [golang-set](https://github.com/deckarep/golang-set/blob/main/README.md): the missing generic set collection for the Go language.

* [Packages or types named "mapset"](https://pkg.go.dev/search?q=mapset) at `pkg.go.dev`

### Talks

* [Prática de Conjuntos](https://youtu.be/3MqEWOpBKpo) (presented at GopherCon Brasil 2019; in Portuguese)

### Previous studies

* [set-practice](https://github.com/ramalho/set-practice): talk with Go and Python examples

* [runeset](https://github.com/ramalho/runeset): simple set type for Unicode characters in Go (a.k.a. runes)

* [runefinder](https://github.com/ramalho/runefinder): inverted index mapping words to set of runes

* [strset](https://github.com/ramalho/strset): full-featured Go set type for `string` elements

* [intset](https://github.com/ramalho/intset): set type for integer elements in Go, based on the `intset` example from Chapter 6 of
[The Go Programming Language](https://www.gopl.io/), by Alan A. A. Donovan & Brian W. Kernighan.