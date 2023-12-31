package erroriter

import "github.com/KSpaceer/itertools"

// Map creates new ErrorIterator which contains elements of type U
// produced by applying mapper to elements of source iterator.
func Map[T, U any](i *itertools.Iterator[T], mapper func(T) (U, error)) *ErrorIterator[U] {
	var zero U
	return New(func() (U, error) {
		if !i.Next() {
			return zero, ErrIterationStop
		}
		return mapper(i.Elem())
	})
}
