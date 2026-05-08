# Sets in modern Go

* Session (60 min)
* Audience Level: intermediate

## Abstract

Set theory is closely related to Boolean logic and databases. In
practice, set operations can dramatically simplify code, replacing
complex nested loops and ifs with declarative expressions that are
easier to read and may offer better performance. That's why most
popular, modern languages provide sets in their standard libraries,
but not yet Go. In this talk I show how set operations simplify common
programming tasks, and how to implement a generic set type in modern
Go.

## Detailed Description

Before the introduction of generics, I presented a keynote at
Gophercon Brazil titled "Set Practice: Using Sets in Go". The talk was
split into four parts: why sets are useful for daily programming
tasks, an overview of sets in some popular languages, set libraries
available for Go, and building a set type from scratch using
bitfields. The talk was well received, but the last two parts were
overcomplicated. Without generics, code generation was the most common
solution in the popular set libraries back then. The lack of iterators
was a minor annoyance. Now with generics and iterators, Go can finally
support sets as well as languages like Python, C#, Ruby, JavaScript,
Java etc. This talk is about the design of the set type I'd like to
see in the Go standard library.
