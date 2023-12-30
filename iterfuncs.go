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
