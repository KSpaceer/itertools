package itertools

// Pair is 2-size tuple of heterogeneous values.
type Pair[T, U any] struct {
	First  T
	Second U
}

func (p Pair[T, U]) Unpack() (T, U) {
	return p.First, p.Second
}

// Enumeration is a specific case of Pair for Enumerate function.
type Enumeration[T any] Pair[T, int]
