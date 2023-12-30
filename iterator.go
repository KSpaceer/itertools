package itertools

type Iterator[T any] struct {
	f          func() (T, bool)
	value      T
	canProceed bool
}

func New[T any](f func() (T, bool)) *Iterator[T] {
	return &Iterator[T]{
		f:          f,
		canProceed: true,
	}
}

func (i *Iterator[T]) Next() bool {
	if !i.canProceed {
		return false
	}
	i.value, i.canProceed = i.f()
	return i.canProceed
}

func (i *Iterator[T]) Elem() T {
	return i.value
}

func (i *Iterator[T]) Count() int {
	var count int
	for i.Next() {
		count++
	}
	return count
}

func (i *Iterator[T]) Drop(count int) int {
	var droppedCount int
	for ; droppedCount < count && i.Next(); droppedCount++ {
	}
	return droppedCount
}

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

func (i *Iterator[T]) Range(f func(T) bool) {
	for i.Next() {
		if !f(i.Elem()) {
			return
		}
	}
}

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

func (i *Iterator[T]) Collect() []T {
	var elems []T
	for i.Next() {
		elems = append(elems, i.Elem())
	}
	return elems
}

func (i *Iterator[T]) Reduce(acc T, f func(acc T, elem T) T) T {
	for i.Next() {
		acc = f(acc, i.Elem())
	}
	return acc
}

func (i *Iterator[T]) All(f func(T) bool) bool {
	for i.Next() {
		if !f(i.Elem()) {
			return false
		}
	}
	return true
}

func (i *Iterator[T]) Any(f func(T) bool) bool {
	for i.Next() {
		if f(i.Elem()) {
			return true
		}
	}
	return false
}

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
