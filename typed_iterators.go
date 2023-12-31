package itertools

import (
	"reflect"
	"unicode/utf8"
)

// NewSliceIterator creates iterator for given slice,
// meaning iterator will yield all elements of the slice.
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

// NewChanIterator creates iterator yielding values from channel
// until the channel is closed.
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

// NewMapIterator creates iterator yielding key-value pairs from map.
// NewMapIterator uses reflect package to keep iteration state.
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

// NewMapKeysIterator creates iterator yielding keys from map.
// NewMapKeysIterator uses reflect package to keep iteration state.
func NewMapKeysIterator[K comparable, V any](m map[K]V) *Iterator[K] {
	mapIter := reflect.ValueOf(m).MapRange()
	return New(func() (K, bool) {
		if !mapIter.Next() {
			var zero K
			return zero, false
		}
		return mapIter.Key().Interface().(K), true
	})
}

// NewMapValuesIterator creates iterator yielding values from map.
// NewMapValuesIterator uses reflect package to keep iteration state.
func NewMapValuesIterator[K comparable, V any](m map[K]V) *Iterator[V] {
	mapIter := reflect.ValueOf(m).MapRange()
	return New(func() (V, bool) {
		if !mapIter.Next() {
			var zero V
			return zero, false
		}
		return mapIter.Value().Interface().(V), true
	})
}

// NewAsciiIterator creates iterator yielding bytes from string
// (interpreting string as []byte).
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

// NewUTF8Iterator creates iterator yielding runes from string
// (interpreting string as []rune)
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
