package itertools

import (
	"cmp"
	"golang.org/x/exp/constraints"
)

// Chain chains iterators, returning resulting chained iterator.
// The result iterator yields elements of the first iterator,
// then elements of the second one etc.
func Chain[T any](iters ...*Iterator[T]) *Iterator[T] {
	var (
		i    int
		zero T
	)
	return New(func() (T, bool) {
		for i < len(iters) {
			v, ok := iters[i].f()
			if ok {
				return v, ok
			}
			i++
		}
		return zero, false
	})
}

// Zip joins two iterators into a one yielding Pair of the iterators' elements.
// Returned iterator yields Pairs until one of source iterators is empty.
func Zip[T, U any](t *Iterator[T], u *Iterator[U]) *Iterator[Pair[T, U]] {
	return New(func() (Pair[T, U], bool) {
		tElem, ok := t.f()
		if !ok {
			return Pair[T, U]{}, false
		}
		uElem, ok := u.f()
		if !ok {
			return Pair[T, U]{}, false
		}
		return Pair[T, U]{
			First:  tElem,
			Second: uElem,
		}, true
	})
}

// Map returns new iterator that yields elements of type U
// by calling mapper to each element of type T of source iterator.
func Map[T, U any](i *Iterator[T], mapper func(T) U) *Iterator[U] {
	var zero U
	return New(func() (U, bool) {
		v, ok := i.f()
		if !ok {
			return zero, false
		}
		return mapper(v), true
	})
}

// Max return max value of iterator.
func Max[T cmp.Ordered](i *Iterator[T]) T {
	return i.Max(cmp.Compare[T])
}

// Min return min value of iterator.
func Min[T cmp.Ordered](i *Iterator[T]) T {
	return i.Max(func(a T, b T) int {
		return -cmp.Compare(a, b)
	})
}

type Summable interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

// Sum returns sum of iterator elements
func Sum[T Summable](i *Iterator[T]) T {
	var zero T
	return i.Reduce(zero, func(acc T, elem T) T {
		return acc + elem
	})
}

// Find applies function f to elements of iterator, returning
// first element for which the function returned true.
// The returned boolean value shows if the element was found (i.e. is valid).
// If no element was found, Find returns false as second returned value.
func Find[T any](i *Iterator[T], f func(T) bool) (T, bool) {
	var found T
	for i.Next() {
		found = i.Elem()
		if f(found) {
			return found, true
		}
	}
	return found, false
}

// Enumerate creates new iterator that returns Enumeration contating
// current element of source iterator along with current iteration count (starting from 0).
func Enumerate[T any](i *Iterator[T]) *Iterator[Enumeration[T]] {
	var idx int
	return New(func() (Enumeration[T], bool) {
		v, ok := i.f()
		if !ok {
			return Enumeration[T]{}, false
		}
		result := Enumeration[T]{
			First:  v,
			Second: idx,
		}
		idx++
		return result, true
	})
}

// Batched creates new iterator that returns slices of T (aka batch)
// with size up to batchSize, using given source iterator.
func Batched[T any](i *Iterator[T], batchSize int) *Iterator[[]T] {
	if batchSize <= 0 {
		return New(func() ([]T, bool) {
			return nil, false
		})
	}
	var stopped bool
	return New(func() ([]T, bool) {
		if stopped {
			return nil, false
		}
		result := make([]T, 0, batchSize)
		for count := 0; count < batchSize; count++ {
			v, ok := i.f()
			if !ok {
				stopped = true
				if len(result) > 0 {
					break
				}
				return nil, false
			}
			result = append(result, v)
		}
		return result, true
	})
}

// Repeat creates new iterator that endlessly yields elem.
func Repeat[T any](elem T) *Iterator[T] {
	return New(func() (T, bool) {
		return elem, true
	})
}

// Cycle creates new iterator that endlessly repeats elements of source iterator.
// If source iterator is empty, the cycle iterator is also empty.
func Cycle[T any](i *Iterator[T]) *Iterator[T] {
	const (
		original = iota
		cycled
		empty
	)
	var (
		elems []T
		idx   int
	)
	state := original
	return New(func() (T, bool) {
		switch state {
		case original:
			v, ok := i.f()
			if ok {
				elems = append(elems, v)
				return v, true
			}

			if len(elems) == 0 {
				state = empty
				return v, false
			}

			state = cycled
			fallthrough
		case cycled:
			v := elems[idx]
			idx = (idx + 1) % len(elems)
			return v, true
		default:
			var zero T
			return zero, false
		}
	})
}
