type Set[T comparable] interface {
	Add(val T) bool
	Append(val ...T) int
	AppendFrom(other Set[T]) int
	Cardinality() int
	Clear()
	Clone() Set[T]
	Contains(val ...T) bool
	ContainsOne(val T) bool
	ContainsAny(val ...T) bool
	ContainsAnyElement(other Set[T]) bool
	Difference(other Set[T]) Set[T]
	Equal(other Set[T]) bool
	Intersect(other Set[T]) Set[T]
	IsEmpty() bool
	IsProperSubset(other Set[T]) bool
	IsProperSuperset(other Set[T]) bool
	IsSubset(other Set[T]) bool
	IsSuperset(other Set[T]) bool
	Each(func(T) bool)
	Filter(func(T) bool) Set[T]
	Iter() <-chan T
	Iterator() *Iterator[T]
	Remove(i T)
	RemoveAll(i ...T)
	String() string
	SymmetricDifference(other Set[T]) Set[T]
	Union(other Set[T]) Set[T]
	Pop() (T, bool)
	PopN(n int) ([]T, int)
	ToSlice() []T
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
	MarshalBSONValue() (bsontype.Type, []byte, error)
	UnmarshalBSONValue(bt bsontype.Type, b []byte) error
}