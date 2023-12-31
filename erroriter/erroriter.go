package erroriter

import (
	"errors"
	"github.com/KSpaceer/itertools"
)

var ErrIterationStop = errors.New("iteration stop")

// ErrorIterator is an iterator yielding itertools.Pair of value of type T and error.
type ErrorIterator[T any] struct {
	*itertools.Iterator[itertools.Pair[T, error]]
}

// New creates ErrorIterator that yields elements using function f.
// Iterator yields elements until returned error is ErrIterationStop.
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

// Result unpacks itertools.Pair element, returning value and error.
func (i *ErrorIterator[T]) Result() (T, error) {
	v := i.Elem()
	return v.First, v.Second
}

// CollectUntilError collects actual values into slice.
// CollectUntilError returns this slice or first encountered error.
func (i *ErrorIterator[T]) CollectUntilError() ([]T, error) {
	var results []T
	for i.Next() {
		v, err := i.Result()
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}
	return results, nil
}
