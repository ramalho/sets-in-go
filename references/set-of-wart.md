# `set.Of` returns the wrong type

Notes on a rough edge in [CL 745441](https://go-review.googlesource.com/c/go/+/745441),
Alan Donovan's proposed `container/set`, as of patchset 13 — still the current
patchset on 2026-08-08. The code discussed here is vendored in `../vendored`;
every compiler claim below was checked against it with Go 1.27rc2.

The constructor is declared to return `map[E]struct{}` rather than `Set[E]`, so
its result carries none of the type's 22 methods. It is a small slip with
surprisingly deep roots, because it sits on a genuine Go design tension: a
defined type buys you a method set and costs you assignability, and no single
signature gets both.

## First: it is almost certainly not deliberate

```go
// Collect creates a new set containing the elements of the sequence.
func Collect[E comparable](seq iter.Seq[E]) Set[E] {
	return Set[E](mapset.Collect(seq))
}

// Of creates a new set containing the elements of the sequence.
func Of[E comparable](elems ...E) map[E]struct{} {
	return Set[E](mapset.Of(elems...))
}
```

Two things in `Of` point the same direction. The body already converts to
`Set[E]`, and the return statement then silently widens it back to the unnamed
type — nobody writes that on purpose, they just return `mapset.Of(elems...)`.
And the doc comment says "the elements of the sequence" for a function taking
variadic elements; that wording belongs to `Collect`, directly above it.

Both look like copy-paste residue. This is code under review.

## The genuine case for returning `map[E]struct{}`

There is a real argument, and it is worth taking seriously.

Go's assignability rule: a value of type V may be assigned to a variable of type
T when they share an underlying type and **at least one of them is unnamed**. So
the return type decides how far the result can travel.

Given a user's own `type MySet map[string]struct{}`:

| Assignment | Compiles? |
| --- | --- |
| `var s set.Set[string] = set.Of("x")` | yes |
| `var m MySet = set.Of("x")` | yes |
| `var m MySet = set.Collect(seq)` — `Collect` returns `Set[E]` | **no** |

The third fails with:

```
cannot use set.Collect(b.All()) (value of map type set.Set[string])
	as MySet value in variable declaration
```

So returning the unnamed type does let the result flow into anyone's
named map type without a conversion, and returning `Set[E]` does not.
It is the same reason `slices.Collect` returns a plain `[]E`.

## Why the case does not hold up here

The `slices` precedent does not transfer: slices have no method set to lose.
Here the defined type *is* the API. Returning the unnamed type discards the
entire reason `Set` is a defined type rather than an alias.

**Chaining dies.** `set.Of(1, 2, 3).Intersection(y)` does not compile:

```
x.Intersection undefined (type map[int]struct{} has no field or method Intersection)
```

Change the signature to return `Set[E]` and it compiles and prints `{3}`. Since
the package's pitch is that set algebra reads better as expressions, losing
composition at the most common construction site cuts against the proposal
itself.

**Inference silently yields the wrong type.** `s := set.Of(1, 2, 3)` infers
`map[int]struct{}`; with the named return it infers `set.Set[int]`. This is the
worst part, because the failure is delayed — the `:=` succeeds, and you learn
about it several lines later when `s.Contains(x)` fails with a message about a
map type having no methods.

The named-return failure mode is the opposite: `var m MySet = set.Of("x")` fails
immediately, at the assignment, and the fix — wrap it in `MySet(...)` — is
obvious and free at runtime. Local errors with obvious fixes beat action at a
distance.

**It contradicts its own sibling.** `Collect` returns `Set[E]` one line above.
Two constructors for the same concept returning different types costs every
reader a trip to the godoc, forever.

**The maximally-assignable constructor already exists.** This is the decisive
point. [`mapset.Of`](https://go-review.googlesource.com/c/go/+/724420) is
declared:

```go
func Of[K comparable](elems ...K) map[K]struct{}
```

It already returns the unnamed type, precisely so it can flow into any named map
type. That is the documented division of labor: `mapset` serves arbitrary and
legacy map-shaped sets, `set` provides the canonical named one.

Under the current signature `set.Of` is a near-duplicate of `mapset.Of` that
adds nothing, while the population it optimizes for — people with their own
`type MySet map[string]struct{}` — is exactly the fragmentation the proposal
exists to end, and is already served one package over.

## The third option, and why it is closed

You might try to have both by making the result type a parameter, mirroring how
`mapset` is written throughout:

```go
func Of[M ~map[E]struct{}, E comparable](elems ...E) M
```

This does not work. There is no argument of type `M` to infer from:

```
in call to OfGeneric, cannot infer M
```

Every call site would have to spell out `set.Of[MySet]("x")`, which is worse
than a conversion at the one place you need it. `maps.Clone[M ~map[K]V](m M) M`
gets away with the pattern only because it takes an `M` as input; a constructor
has nothing to infer from.

So the design space really is just the two signatures.

## Recommendation

Return `Set[E]`.

The cost is an explicit conversion for people assigning into their own named map
types — compile-time only, immediately diagnosed, and already avoidable by
calling `mapset.Of` instead. The benefit is a constructor consistent with
`Collect`, inferring the intended type under `:=`, and producing values you can
actually call methods on.

## Note for the talk

This is a compact illustration of a Go-specific tension that has no clean
resolution: **defined types trade assignability for a method set.** The
standard library answers the question both ways.

Plain `[]byte`, no method set, when the value is inert data: `io.ReadAll`,
`json.Marshal`, `(*bytes.Buffer).Bytes()`. Nothing to attach behavior to, so
the unnamed type keeps the result assignable into anything.

A named byte-slice type when the methods encode semantics a raw byte
comparison would get wrong. `net.IP` is `type IP []byte`, and its `Equal`
method is documented to treat a 4-byte IPv4 address and its 16-byte
IPv4-in-IPv6 form as equal — `bytes.Equal` on the raw slices would say they
differ, since the lengths don't even match. `String` similarly renders that
same dual form and RFC 5952 zero-run compression correctly, which
`fmt.Sprintf("%v", rawBytes)` cannot. `net.HardwareAddr` (`type HardwareAddr
[]byte`) makes the same call for MAC addresses, with `String` formatting as
`xx:xx:xx:xx:xx:xx`.

The mirror is not exact: `[]byte` is itself unnamed, so `net.IP` is still
freely assignable to and from plain `[]byte`. The case that actually matches
`Set[E]`/`MySet` is assigning `net.IP` into a second named byte-slice type —
say `type Payload []byte` — which does need an explicit `Payload(ip)`
conversion, for the same reason `set.Collect`'s result needs `MySet(...)`.

It is also live. The CL is open, so this may be fixed before the talk — check
the current patchset, and re-extract with:

```sh
cd ../vendored && SET_PS=current ./refresh.sh
```

Runnable examples of the affected constructions are in
`../examples/proposed_test.go`.
