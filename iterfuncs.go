package itertools

import "cmp"

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

func Max[T cmp.Ordered](i *Iterator[T]) T {
	return i.Max(cmp.Compare[T])
}

func Min[T cmp.Ordered](i *Iterator[T]) T {
	return i.Max(func(a T, b T) int {
		return -cmp.Compare(a, b)
	})
}

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

func Repeat[T any](elem T) *Iterator[T] {
	return New(func() (T, bool) {
		return elem, true
	})
}

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
