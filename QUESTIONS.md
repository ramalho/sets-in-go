# Questions

## Why the recursive syntax?

Why does `_AbstractSet` appear in the declaration of `_AbstractSet?` 


```go
type _AbstractSet[E any, S _AbstractSet[E, S]] interface {
```


## Why does the concrete type appear on both sides?

```go
var _ _AbstractSet[int, set.Set[int]] = make(set.Set[int])

var _ _AbstractSet[int, *hash.Set[int]] = new(hash.Set[int])
```
