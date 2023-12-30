package itertools

import (
	"reflect"
	"unicode/utf8"
)

func NewSliceIterator[S ~[]T, T any](s S) *Iterator[T] {
	var (
		idx  int
		zero T
	)
	return New(func() (T, bool) {
		if idx >= len(s) {
			return zero, false
		}
		v := s[idx]
		idx++
		return v, true
	})
}

func NewChanIterator[T any](ch <-chan T) *Iterator[T] {
	var zero T
	return New(func() (T, bool) {
		v, ok := <-ch
		if !ok {
			return zero, false
		}
		return v, true
	})
}

func NewMapIterator[K comparable, V any](m map[K]V) *Iterator[Pair[K, V]] {
	mapIter := reflect.ValueOf(m).MapRange()
	return New(func() (Pair[K, V], bool) {
		if !mapIter.Next() {
			return Pair[K, V]{}, false
		}
		return Pair[K, V]{
			First:  mapIter.Key().Interface().(K),
			Second: mapIter.Value().Interface().(V),
		}, true
	})
}

func NewAsciiIterator(s string) *Iterator[byte] {
	var (
		idx  int
		zero byte
	)
	return New(func() (byte, bool) {
		if idx >= len(s) {
			return zero, false
		}
		v := s[idx]
		idx++
		return v, true
	})
}

func NewUTF8Iterator(s string) *Iterator[rune] {
	return New(func() (rune, bool) {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			return 0, false
		}
		s = s[size:]
		return r, true
	})
}
