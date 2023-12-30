package itertools

type Pair[T, U any] struct {
	First  T
	Second U
}

type Enumeration[T any] Pair[T, int]
