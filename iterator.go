package itertools

import "slices"

// Iterator is used to process all elements of some collection
// or sequence. Iterator contains methods to access the elements
// and to process or aggregate them in many ways.
type Iterator[T any] struct {
	f          func() (T, bool)
	value      T
	canProceed bool
}

// New creates new Iterator using given iteration function.
// Function returns an element of collection and boolean value indicating
// if the element is valid (i.e. false means that the iteration is over).
func New[T any](f func() (T, bool)) *Iterator[T] {
	return &Iterator[T]{
		f:          f,
		canProceed: true,
	}
}

// Next proceeds iterator to the next element, returning boolean value
// to show that said element exists.
func (i *Iterator[T]) Next() bool {
	if !i.canProceed {
		return false
	}
	i.value, i.canProceed = i.f()
	return i.canProceed
}

// Elem returns the current element of iterator.
// If iterator is empty (Next returns false), result is unspecified.
func (i *Iterator[T]) Elem() T {
	return i.value
}

// Count returns amount of remaining elements in iterator.
// Call of Count consumes all elements.
func (i *Iterator[T]) Count() int {
	var count int
	for i.Next() {
		count++
	}
	return count
}

// Drop skips next n elements in iterator, returning
// amount of skipped elements (if iterator has fewer elements than n, returned value
// is equal to the amount of elements).
func (i *Iterator[T]) Drop(n int) int {
	var droppedCount int
	for ; droppedCount < n && i.Next(); droppedCount++ {
	}
	return droppedCount
}

// Limit produces new iterator that can return at most size elements.
func (i *Iterator[T]) Limit(size int) *Iterator[T] {
	var zero T
	if size <= 0 {
		return New(func() (T, bool) {
			return zero, false
		})
	}

	var count int
	return New(func() (T, bool) {
		if count >= size || !i.Next() {
			return zero, false
		}
		v := i.Elem()
		count++
		return v, true
	})
}

// WithStep produces new iterator that yields every "step"th element of underlying iterator
// If step is non-positive, returns empty iterator.
func (i *Iterator[T]) WithStep(step int) *Iterator[T] {
	var zero T
	if step <= 0 {
		return New(func() (T, bool) {
			return zero, false
		})
	}

	var count = -1
	return New(func() (T, bool) {
		for {
			v, ok := i.f()
			if !ok {
				return zero, false
			}
			count++
			if count%step == 0 {
				return v, ok
			}
		}
	})
}

// Range calls function f for every element of iterator until the function
// returns false
func (i *Iterator[T]) Range(f func(T) bool) {
	for i.Next() {
		if !f(i.Elem()) {
			return
		}
	}
}

// Filter produces new iterator that yields only elements
// for which function f returns true.
func (i *Iterator[T]) Filter(f func(T) bool) *Iterator[T] {
	var zero T
	return New(func() (T, bool) {
		for {
			v, ok := i.f()
			if !ok {
				return zero, false
			}
			if f(v) {
				return v, true
			}
		}
	})
}

// Collect returns all elements of iterator as slice.
func (i *Iterator[T]) Collect(opts ...AllocationOption) []T {
	var options allocOptions
	for _, opt := range opts {
		opt(&options)
	}
	elems := make([]T, 0, options.preallocSize)
	for i.Next() {
		elems = append(elems, i.Elem())
	}
	return elems
}

// Reduce applies given function f to every element of iterator,
// using previous accumulating state and returning updated accumulating state
// on each iteration.
// Reduce also accepts initial value for accumulating state.
// Reduce returns final accumulating state created after applying f
// to all elements of iterator.
func (i *Iterator[T]) Reduce(acc T, f func(acc T, elem T) T) T {
	for i.Next() {
		acc = f(acc, i.Elem())
	}
	return acc
}

// All applies function f to every element of iterator.
// If f returns true for all elements of iterator, All returns true.
// Otherwise, All returns false.
// All is lazy and will stop iterating after first element for which f returns false.
// All returns true for empty iterator.
func (i *Iterator[T]) All(f func(T) bool) bool {
	for i.Next() {
		if !f(i.Elem()) {
			return false
		}
	}
	return true
}

// Any applies function f to every element of iterator.
// If f returns false for all elements of iterator, Any returns false.
// Otherwise, Any return true.
// Any is lazy and will stop iterating after first element for which f returns true.
// Any returns false for empty iterator.
func (i *Iterator[T]) Any(f func(T) bool) bool {
	for i.Next() {
		if f(i.Elem()) {
			return true
		}
	}
	return false
}

// Max returns max element of iterator, using provided comparison function.
// Comparison function returns next results:
//   - -1: if the first argument is less than second one
//   - 0: if two arguments are equal
//   - 1: if the first argument is greater than second one
func (i *Iterator[T]) Max(f func(T, T) int) T {
	if !i.Next() {
		var zero T
		return zero
	}
	maxValue := i.Elem()
	for i.Next() {
		v := i.Elem()
		if f(maxValue, v) < 0 {
			maxValue = v
		}
	}
	return maxValue
}

// SortedBy returns new iterator yielding elements of source iterator in sorted order.
// Sort order is defined by cmp.
// Comparison function cmp returns next results:
//   - -1: if the first argument is less than second one
//   - 0: if two arguments are equal
//   - 1: if the first argument is greater than second one
func (i *Iterator[T]) SortedBy(cmp func(T, T) int, opts ...AllocationOption) *Iterator[T] {
	values := i.Collect(opts...)
	slices.SortFunc(values, cmp)
	return NewSliceIterator(values)
}
