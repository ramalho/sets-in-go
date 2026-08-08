package mapset

import (
    "iter"
)

// Collect returns a new set containing the elements of the sequence.
func Collect[K comparable](seq iter.Seq[K]) map[K]struct{} {
    return collect[K, struct{}](seq)
}

/// ...

func collect[K comparable, V bool | struct{}](seq iter.Seq[K]) map[K]V {
    x := make(map[K]V)
    InsertAll(x, seq)
    return x
}

/// ...

// InsertAll adds each element of the addenda sequence to the set.
// If the set values are boolean, the value 'true' is used.
// It reports whether len(x) changed.
func InsertAll[M ~map[K]V, K comparable, V bool | struct{}](x M, addenda iter.Seq[K]) bool {
    pre := len(x)
    for k := range addenda {
        insert(x, k)
    }
    return len(x) != pre
}

func insert[M ~map[K]V, K comparable, V bool | struct{}](m M, k K) {
    // Choose the distinguished "present" value (true or struct{}{}).
    // This compiles to a load from .rodata.
    var present V
    if _, ok := any(present).(bool); ok {
        present = any(true).(V)
    }

    // This is the canonical insertion operation.
    // All maps created by this API use only the
    // distinguished 'present' value for the result type.
    m[k] = present
}
