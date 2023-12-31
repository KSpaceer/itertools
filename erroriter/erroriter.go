package erroriter

import (
	"errors"
	"github.com/KSpaceer/itertools"
)

var ErrIterationStop = errors.New("iteration stop")

type ErrorIterator[T any] struct {
	*itertools.Iterator[itertools.Pair[T, error]]
}

func New[T any](f func() (T, error)) *ErrorIterator[T] {
	i := itertools.New[itertools.Pair[T, error]](func() (itertools.Pair[T, error], bool) {
		v, err := f()
		if errors.Is(err, ErrIterationStop) {
			return itertools.Pair[T, error]{}, false
		}
		return itertools.Pair[T, error]{
			First:  v,
			Second: err,
		}, true
	})
	return &ErrorIterator[T]{Iterator: i}
}

func (i *ErrorIterator[T]) Result() (T, error) {
	v := i.Elem()
	return v.First, v.Second
}
