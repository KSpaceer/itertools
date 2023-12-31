package erroriter

import "github.com/KSpaceer/itertools"

func Map[T, U any](i *itertools.Iterator[T], mapper func(T) (U, error)) *ErrorIterator[U] {
	var zero U
	return New(func() (U, error) {
		if !i.Next() {
			return zero, ErrIterationStop
		}
		return mapper(i.Elem())
	})
}
