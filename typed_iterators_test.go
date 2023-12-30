package itertools_test

import (
	"cmp"
	"fmt"
	"github.com/KSpaceer/itertools"
	"math"
	"slices"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestSliceIterator(t *testing.T) {
	s := []int{89, 716, 122, 151, 475, 415, 512, 631, 258, 635, 451, 832, 476}

	t.Run("collect", func(t *testing.T) {
		result := itertools.NewSliceIterator(s).Collect()

		if cmpResult := sliceCmp(s, result); cmpResult != 0 {
			t.Errorf("expected %v, got %v", s, result)
		}
	})

	t.Run("count", func(t *testing.T) {
		result := itertools.NewSliceIterator(s).Count()

		if result != len(s) {
			t.Errorf("expected Count to return %d, but got %d", len(s), result)
		}
	})

	t.Run("drop", func(t *testing.T) {
		i := itertools.NewSliceIterator(s)
		dropCount := 8

		dropped := i.Drop(dropCount)

		if dropped != dropCount {
			t.Errorf("expected to have %d dropped values, but got %d", dropCount, dropped)
		}

		result := i.Collect()
		expected := s[dropCount:]

		if cmpResult := sliceCmp(expected, result); cmpResult != 0 {
			t.Errorf("expected %v, but got %v", expected, result)
		}

		if dropped = i.Drop(1); dropped != 0 {
			t.Errorf("expected to have nothing remained after iteration, but had %d values", dropped)
		}
	})

	t.Run("drop overflow", func(t *testing.T) {
		i := itertools.NewSliceIterator(s)
		dropCount := 120

		dropped := i.Drop(dropCount)

		if dropped != len(s) {
			t.Errorf("expected to have %d dropped values, but got %d", len(s), dropped)
		}

		result := i.Collect()
		expected := []int{}

		if cmpResult := sliceCmp(expected, result); cmpResult != 0 {
			t.Errorf("expected %v, but got %v", expected, result)
		}

		if dropped = i.Drop(1); dropped != 0 {
			t.Errorf("expected to have nothing remained after drop, but had %d values", dropped)
		}
	})

	t.Run("with step", func(t *testing.T) {
		type tcase struct {
			name     string
			step     int
			expected []int
		}

		tcases := []tcase{
			{
				name:     "singular step",
				step:     1,
				expected: s,
			},
			{
				name:     "triple step",
				step:     3,
				expected: []int{89, 151, 512, 635, 476},
			},
			{
				name: "zero step",
				step: 0,
			},
			{
				name: "negative step",
				step: -1,
			},
			{
				name:     "giant step",
				step:     len(s),
				expected: []int{s[0]},
			},
			{
				name:     "two values",
				step:     len(s) - 1,
				expected: []int{s[0], s[len(s)-1]},
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				i := itertools.NewSliceIterator(s).WithStep(tc.step)
				result := i.Collect()

				if cmpResult := sliceCmp(tc.expected, result); cmpResult != 0 {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			})
		}
	})
	t.Run("range", func(t *testing.T) {
		i := itertools.NewSliceIterator(s)
		var sb strings.Builder
		i.Range(func(n int) bool {
			sb.WriteString(fmt.Sprint(n))
			return true
		})

		var expectedSB strings.Builder
		for _, n := range s {
			expectedSB.WriteString(fmt.Sprint(n))
		}

		result := sb.String()
		expected := expectedSB.String()

		if result != expected {
			t.Errorf("expected %s as concatenated string with range, but got %s", expected, result)
		}
	})
	t.Run("filter", func(t *testing.T) {
		type tcase struct {
			name       string
			filterFunc func(int) bool
			expected   []int
		}

		tcases := []tcase{
			{
				name: "odd only",
				filterFunc: func(n int) bool {
					return n%2 == 1
				},
				expected: []int{89, 151, 475, 415, 631, 635, 451},
			},
			{
				name: "even only",
				filterFunc: func(n int) bool {
					return n%2 == 0
				},
				expected: []int{716, 122, 512, 258, 832, 476},
			},
			{
				name: "divisible by 11 only",
				filterFunc: func(n int) bool {
					return n%11 == 0
				},
				expected: []int{451},
			},
			{
				name: "greater than 500",
				filterFunc: func(n int) bool {
					return n > 500
				},
				expected: []int{716, 512, 631, 635, 832},
			},
			{
				name: "pass all",
				filterFunc: func(int) bool {
					return true
				},
				expected: s,
			},
			{
				name: "pass none",
				filterFunc: func(int) bool {
					return false
				},
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				i := itertools.NewSliceIterator(s).Filter(tc.filterFunc)
				result := i.Collect()

				if cmpResult := sliceCmp(tc.expected, result); cmpResult != 0 {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			})
		}
	})

	t.Run("reduce", func(t *testing.T) {
		type tcase struct {
			name         string
			reducer      func(int, int) int
			initialValue int
			expected     int
		}

		tcases := []tcase{
			{
				name: "sum",
				reducer: func(acc int, n int) int {
					return acc + n
				},
				initialValue: 0,
				expected: func() int {
					var sm int
					for _, n := range s {
						sm += n
					}
					return sm
				}(),
			},
			{
				name: "mul",
				reducer: func(acc int, n int) int {
					return (acc * n) % 1_000_000_007
				},
				initialValue: 1,
				expected: func() int {
					m := 1
					for _, n := range s {
						m *= n
						m %= 1_000_000_007
					}
					return m
				}(),
			},
			{
				name: "conditional div",
				reducer: func(acc int, n int) int {
					if n < 250 {
						acc /= n
					}
					return acc
				},
				initialValue: math.MaxInt,
				expected: func() int {
					res := math.MaxInt
					for _, n := range s {
						if n < 250 {
							res /= n
						}
					}
					return res
				}(),
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewSliceIterator(s).
					Reduce(tc.initialValue, tc.reducer)

				if result != tc.expected {
					t.Errorf("expected %d after reduce, got %d", tc.expected, result)
				}
			})
		}
	})

	t.Run("all", func(t *testing.T) {
		type tcase struct {
			name      string
			condition func(int) bool
			expected  bool
		}

		tcases := []tcase{
			{
				name: "true",
				condition: func(int) bool {
					return true
				},
				expected: true,
			},
			{
				name: "does not end with 7",
				condition: func(n int) bool {
					return n%10 != 7
				},
				expected: true,
			},
			{
				name: "ends with 6",
				condition: func(n int) bool {
					return n%10 == 6
				},
				expected: false,
			},
			{
				name: "greater than 70",
				condition: func(n int) bool {
					return n > 70
				},
				expected: true,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewSliceIterator(s).
					All(tc.condition)
				if result != tc.expected {
					t.Errorf("expected %t in all, but got %t", tc.expected, result)
				}
			})
		}
	})

	t.Run("any", func(t *testing.T) {
		type tcase struct {
			name      string
			condition func(int) bool
			expected  bool
		}

		tcases := []tcase{
			{
				name: "true",
				condition: func(int) bool {
					return true
				},
				expected: true,
			},
			{
				name: "does not ends with 6",
				condition: func(n int) bool {
					return n%10 != 6
				},
				expected: true,
			},
			{
				name: "ends with 7",
				condition: func(n int) bool {
					return n%10 == 7
				},
				expected: false,
			},
			{
				name: "less than 1000",
				condition: func(n int) bool {
					return n < 1000
				},
				expected: true,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewSliceIterator(s).
					Any(tc.condition)
				if result != tc.expected {
					t.Errorf("expected %t in any, but got %t", tc.expected, result)
				}
			})
		}
	})

	t.Run("max", func(t *testing.T) {
		type tcase struct {
			name     string
			f        func(int, int) int
			expected int
		}

		tcases := []tcase{
			{
				name:     "maximum",
				f:        cmp.Compare[int],
				expected: 832,
			},
			{
				name: "minimum",
				f: func(a int, b int) int {
					return -cmp.Compare(a, b)
				},
				expected: 89,
			},
			{
				name: "nearest to 500",
				f: func(a int, b int) int {
					adiff := 500 - a
					if adiff < 0 {
						adiff = -adiff
					}
					bdiff := 500 - b
					if bdiff < 0 {
						bdiff = -bdiff
					}
					return cmp.Compare(bdiff, adiff)
				},
				expected: 512,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewSliceIterator(s).Max(tc.f)
				if tc.expected != result {
					t.Errorf("expected %d as max value, but got %d", tc.expected, result)
				}
			})
		}
	})
}

func TestChanIterator(t *testing.T) {
	s := []complex128{
		1,
		1 + 1i,
		1i,
		-1 + 1i,
		-1,
		-1 - 1i,
		-1i,
		1 - 1i,
	}
	chanConstructor := func() <-chan complex128 {
		ch := make(chan complex128)

		go func() {
			for _, n := range s {
				ch <- n
			}
			close(ch)
		}()

		return ch
	}

	t.Run("collect", func(t *testing.T) {
		result := itertools.NewChanIterator(chanConstructor()).Collect()

		if !sliceEqual(s, result) {
			t.Errorf("expected %v, got %v", s, result)
		}
	})

	t.Run("count", func(t *testing.T) {
		result := itertools.NewChanIterator(chanConstructor()).Count()

		if result != len(s) {
			t.Errorf("expected Count to return %d, but got %d", len(s), result)
		}
	})
	t.Run("drop", func(t *testing.T) {
		i := itertools.NewChanIterator(chanConstructor())
		dropCount := 4

		dropped := i.Drop(dropCount)

		if dropped != dropCount {
			t.Errorf("expected to have %d dropped values, but got %d", dropCount, dropped)
		}

		result := i.Collect()
		expected := s[dropCount:]

		if !sliceEqual(expected, result) {
			t.Errorf("expected %v, but got %v", expected, result)
		}

		if dropped = i.Drop(1); dropped != 0 {
			t.Errorf("expected to have nothing remained after iteration, but had %d values", dropped)
		}
	})

	t.Run("drop overflow", func(t *testing.T) {
		i := itertools.NewChanIterator(chanConstructor())
		dropCount := 120

		dropped := i.Drop(dropCount)

		if dropped != len(s) {
			t.Errorf("expected to have %d dropped values, but got %d", len(s), dropped)
		}

		result := i.Collect()
		expected := []complex128{}

		if !sliceEqual(expected, result) {
			t.Errorf("expected %v, but got %v", expected, result)
		}

		if dropped = i.Drop(1); dropped != 0 {
			t.Errorf("expected to have nothing remained after drop, but had %d values", dropped)
		}
	})

	t.Run("with step", func(t *testing.T) {
		type tcase struct {
			name     string
			step     int
			expected []complex128
		}

		tcases := []tcase{
			{
				name:     "singular step",
				step:     1,
				expected: s,
			},
			{
				name:     "double step",
				step:     2,
				expected: []complex128{1, 1i, -1, -1i},
			},
			{
				name: "zero step",
				step: 0,
			},
			{
				name: "negative step",
				step: -1,
			},
			{
				name:     "giant step",
				step:     len(s),
				expected: []complex128{s[0]},
			},
			{
				name:     "two values",
				step:     len(s) - 1,
				expected: []complex128{s[0], s[len(s)-1]},
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				i := itertools.NewChanIterator(chanConstructor()).WithStep(tc.step)
				result := i.Collect()

				if !sliceEqual(tc.expected, result) {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			})
		}
	})

	t.Run("range", func(t *testing.T) {
		i := itertools.NewChanIterator(chanConstructor())
		var sb strings.Builder
		i.Range(func(n complex128) bool {
			sb.WriteString(fmt.Sprint(n))
			return true
		})

		var expectedSB strings.Builder
		for _, n := range s {
			expectedSB.WriteString(fmt.Sprint(n))
		}

		result := sb.String()
		expected := expectedSB.String()

		if result != expected {
			t.Errorf("expected %s as concatenated string with range, but got %s", expected, result)
		}
	})
	t.Run("filter", func(t *testing.T) {
		type tcase struct {
			name       string
			filterFunc func(complex128) bool
			expected   []complex128
		}

		tcases := []tcase{
			{
				name: "top",
				filterFunc: func(n complex128) bool {
					return imag(n) > 0
				},
				expected: []complex128{1 + 1i, 1i, -1 + 1i},
			},
			{
				name: "left",
				filterFunc: func(n complex128) bool {
					return real(n) < 0
				},
				expected: []complex128{-1 + 1i, -1, -1 - 1i},
			},
			{
				name: "bottom",
				filterFunc: func(n complex128) bool {
					return imag(n) < 0
				},
				expected: []complex128{-1 - 1i, -1i, 1 - 1i},
			},
			{
				name: "right",
				filterFunc: func(n complex128) bool {
					return real(n) > 0
				},
				expected: []complex128{1, 1 + 1i, 1 - 1i},
			},
			{
				name: "pass all",
				filterFunc: func(complex128) bool {
					return true
				},
				expected: s,
			},
			{
				name: "pass none",
				filterFunc: func(complex128) bool {
					return false
				},
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				i := itertools.NewChanIterator(chanConstructor()).Filter(tc.filterFunc)
				result := i.Collect()

				if !sliceEqual(tc.expected, result) {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			})
		}
	})

	t.Run("reduce", func(t *testing.T) {
		type tcase struct {
			name         string
			reducer      func(complex128, complex128) complex128
			initialValue complex128
			expected     complex128
		}

		tcases := []tcase{
			{
				name: "sum",
				reducer: func(acc complex128, n complex128) complex128 {
					return acc + n
				},
				initialValue: 0,
				expected:     0,
			},
			{
				name: "mul",
				reducer: func(acc complex128, n complex128) complex128 {
					return acc * n
				},
				initialValue: 1,
				expected: func() complex128 {
					m := 1 + 0i
					for _, n := range s {
						m *= n
					}
					return m
				}(),
			},
			{
				name: "all complex parts abs sum",
				reducer: func(acc complex128, n complex128) complex128 {
					r := real(n)
					if r < 0 {
						r = -r
					}
					i := imag(n)
					if i < 0 {
						i = -i
					}
					return acc + complex(r+i, 0)
				},
				initialValue: 0,
				expected:     12,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewChanIterator(chanConstructor()).
					Reduce(tc.initialValue, tc.reducer)

				if result != tc.expected {
					t.Errorf("expected %v after reduce, got %v", tc.expected, result)
				}
			})
		}
	})

	t.Run("all", func(t *testing.T) {
		type tcase struct {
			name      string
			condition func(complex128) bool
			expected  bool
		}

		tcases := []tcase{
			{
				name: "true",
				condition: func(complex128) bool {
					return true
				},
				expected: true,
			},
			{
				name: "limited by circle with R=2",
				condition: func(n complex128) bool {
					r := real(n)
					i := imag(n)
					return (r*r + i*i) < 4
				},
				expected: true,
			},
			{
				name: "imaginary part is not negative",
				condition: func(n complex128) bool {
					return imag(n) >= 0
				},
				expected: false,
			},
			{
				name: "non zero",
				condition: func(n complex128) bool {
					return real(n) != 0 || imag(n) != 0
				},
				expected: true,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewChanIterator(chanConstructor()).
					All(tc.condition)
				if result != tc.expected {
					t.Errorf("expected %t in all, but got %t", tc.expected, result)
				}
			})
		}
	})

	t.Run("any", func(t *testing.T) {
		type tcase struct {
			name      string
			condition func(complex128) bool
			expected  bool
		}

		tcases := []tcase{
			{
				name: "true",
				condition: func(complex128) bool {
					return true
				},
				expected: true,
			},
			{
				name: "is in third quadrant",
				condition: func(n complex128) bool {
					return real(n) < 0 && imag(n) < 0
				},
				expected: true,
			},
			{
				name: "not limited by circle with R=3",
				condition: func(n complex128) bool {
					r, i := real(n), imag(n)
					return (r*r + i*i) > 9
				},
				expected: false,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewChanIterator(chanConstructor()).
					Any(tc.condition)
				if result != tc.expected {
					t.Errorf("expected %t in any, but got %t", tc.expected, result)
				}
			})
		}
	})

	t.Run("max", func(t *testing.T) {
		type tcase struct {
			name     string
			f        func(complex128, complex128) int
			expected complex128
		}

		tcases := []tcase{
			{
				name: "max real part",
				f: func(a complex128, b complex128) int {
					return cmp.Compare(real(a), real(b))
				},
				expected: 1,
			},
			{
				name: "min real part",
				f: func(a complex128, b complex128) int {
					return cmp.Compare(real(b), real(a))
				},
				expected: -1 + 1i,
			},
			{
				name: "max imaginary part",
				f: func(a complex128, b complex128) int {
					return cmp.Compare(imag(a), imag(b))
				},
				expected: 1 + 1i,
			},
			{
				name: "min imaginary part",
				f: func(a complex128, b complex128) int {
					return cmp.Compare(imag(b), imag(a))
				},
				expected: -1 - 1i,
			},
			{
				name: "max distance",
				f: func(a complex128, b complex128) int {
					ra, ia := real(a), imag(a)
					rb, ib := real(b), imag(b)
					return cmp.Compare(ra*ra+ia*ia, rb*rb+ib*ib)
				},
				expected: 1 + 1i,
			},
			{
				name: "min distance",
				f: func(a complex128, b complex128) int {
					ra, ia := real(a), imag(a)
					rb, ib := real(b), imag(b)
					return cmp.Compare(rb*rb+ib*ib, ra*ra+ia*ia)
				},
				expected: 1,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewChanIterator(chanConstructor()).Max(tc.f)
				if tc.expected != result {
					t.Errorf("expected %v as max value, but got %v", tc.expected, result)
				}
			})
		}
	})
}

func TestMapIterator(t *testing.T) {
	m := map[string]int{
		"A": 10000,
		"B": 5000,
		"C": 2500,
		"D": 1250,
		"E": 625,
	}

	collectedValues := []itertools.Pair[string, int]{
		{"A", 10000},
		{"B", 5000},
		{"C", 2500},
		{"D", 1250},
		{"E", 625},
	}

	t.Run("nil map", func(t *testing.T) {
		i := itertools.NewMapIterator[int, int](nil)

		if i.Next() {
			t.Errorf("expected to have no elements in nil map iterator")
		}
	})

	t.Run("collect", func(t *testing.T) {
		result := itertools.NewMapIterator(m).Collect()

		slices.SortFunc(result, func(a, b itertools.Pair[string, int]) int {
			return cmp.Compare(a.First, b.First)
		})

		if !sliceEqual(collectedValues, result) {
			t.Errorf("expected %v, got %v", collectedValues, result)
		}
	})

	t.Run("count", func(t *testing.T) {
		result := itertools.NewMapIterator(m).Count()

		if result != len(m) {
			t.Errorf("expected Count to return %d, but got %d", len(m), result)
		}
	})

	t.Run("range", func(t *testing.T) {
		i := itertools.NewMapIterator(m)
		var found bool
		key := "B"
		i.Range(func(p itertools.Pair[string, int]) bool {
			if p.First == key {
				found = true
				return false
			}
			return true
		})

		if !found {
			t.Errorf("expected to found map entry with key %s", key)
		}
	})

	t.Run("filter", func(t *testing.T) {
		type tcase struct {
			name       string
			filterFunc func(itertools.Pair[string, int]) bool
			expected   []itertools.Pair[string, int]
		}

		tcases := []tcase{
			{
				name: "letter position is odd",
				filterFunc: func(p itertools.Pair[string, int]) bool {
					return (p.First[0]-'A')%2 == 1
				},
				expected: []itertools.Pair[string, int]{
					{"B", 5000},
					{"D", 1250},
				},
			},
			{
				name: "letter position is even",
				filterFunc: func(p itertools.Pair[string, int]) bool {
					return (p.First[0]-'A')%2 == 0
				},
				expected: []itertools.Pair[string, int]{
					{"A", 10000},
					{"C", 2500},
					{"E", 625},
				},
			},
			{
				name: "greater than 3000",
				filterFunc: func(n itertools.Pair[string, int]) bool {
					return n.Second > 3000
				},
				expected: []itertools.Pair[string, int]{
					{"A", 10000},
					{"B", 5000},
				},
			},
			{
				name: "pass all",
				filterFunc: func(itertools.Pair[string, int]) bool {
					return true
				},
				expected: collectedValues,
			},
			{
				name: "pass none",
				filterFunc: func(itertools.Pair[string, int]) bool {
					return false
				},
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				i := itertools.NewMapIterator(m).Filter(tc.filterFunc)
				result := i.Collect()

				slices.SortFunc(result, func(a, b itertools.Pair[string, int]) int {
					return cmp.Compare(a.First, b.First)
				})

				if !sliceEqual(tc.expected, result) {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}

			})
		}
	})

	t.Run("reduce", func(t *testing.T) {
		type tcase struct {
			name         string
			reducer      func(itertools.Pair[string, int], itertools.Pair[string, int]) itertools.Pair[string, int]
			initialValue itertools.Pair[string, int]
			expected     itertools.Pair[string, int]
		}

		tcases := []tcase{
			{
				name: "sum values",
				reducer: func(acc itertools.Pair[string, int], n itertools.Pair[string, int]) itertools.Pair[string, int] {
					acc.Second += n.Second
					return acc
				},
				initialValue: itertools.Pair[string, int]{},
				expected:     itertools.Pair[string, int]{Second: 19375},
			},
			{
				name: "join keys ordered",
				reducer: func(acc itertools.Pair[string, int], n itertools.Pair[string, int]) itertools.Pair[string, int] {
					if acc.First == "" {
						acc.First = n.First
					} else {
						elems := strings.Split(acc.First, "-")
						i, _ := slices.BinarySearch(elems, n.First)
						elems = slices.Insert(elems, i, n.First)
						acc.First = strings.Join(elems, "-")
					}
					return acc
				},
				initialValue: itertools.Pair[string, int]{},
				expected:     itertools.Pair[string, int]{First: "A-B-C-D-E"},
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewMapIterator(m).
					Reduce(tc.initialValue, tc.reducer)

				if result != tc.expected {
					t.Errorf("expected %v after reduce, got %v", tc.expected, result)
				}
			})
		}
	})

	t.Run("all", func(t *testing.T) {
		type tcase struct {
			name      string
			condition func(itertools.Pair[string, int]) bool
			expected  bool
		}

		tcases := []tcase{
			{
				name: "true",
				condition: func(itertools.Pair[string, int]) bool {
					return true
				},
				expected: true,
			},
			{
				name: "are keys contains only letters",
				condition: func(n itertools.Pair[string, int]) bool {
					return !strings.ContainsFunc(n.First, func(r rune) bool {
						return !unicode.IsLetter(r)
					})
				},
				expected: true,
			},
			{
				name: "values divisible by 3",
				condition: func(n itertools.Pair[string, int]) bool {
					return n.Second%3 == 0
				},
				expected: false,
			},
			{
				name: "keys have length greater than 2",
				condition: func(n itertools.Pair[string, int]) bool {
					return len(n.First) > 2
				},
				expected: false,
			},
			{
				name: "does not have zero value",
				condition: func(n itertools.Pair[string, int]) bool {
					return n.Second != 0
				},
				expected: true,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewMapIterator(m).
					All(tc.condition)
				if result != tc.expected {
					t.Errorf("expected %t in all, but got %t", tc.expected, result)
				}
			})
		}
	})

	t.Run("any", func(t *testing.T) {
		type tcase struct {
			name      string
			condition func(itertools.Pair[string, int]) bool
			expected  bool
		}

		tcases := []tcase{
			{
				name: "false",
				condition: func(itertools.Pair[string, int]) bool {
					return false
				},
				expected: false,
			},
			{
				name: "has value divisible by 4",
				condition: func(n itertools.Pair[string, int]) bool {
					return n.Second%4 == 0
				},
				expected: true,
			},
			{
				name: "has any key with length greater than 2",
				condition: func(n itertools.Pair[string, int]) bool {
					return len(n.First) > 2
				},
				expected: false,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewMapIterator(m).
					Any(tc.condition)
				if result != tc.expected {
					t.Errorf("expected %t in any, but got %t", tc.expected, result)
				}
			})
		}
	})

	t.Run("max", func(t *testing.T) {
		type tcase struct {
			name     string
			f        func(itertools.Pair[string, int], itertools.Pair[string, int]) int
			expected itertools.Pair[string, int]
		}

		tcases := []tcase{
			{
				name: "max value",
				f: func(a itertools.Pair[string, int], b itertools.Pair[string, int]) int {
					return cmp.Compare(a.Second, b.Second)
				},
				expected: itertools.Pair[string, int]{"A", 10000},
			},
			{
				name: "min value",
				f: func(a itertools.Pair[string, int], b itertools.Pair[string, int]) int {
					return cmp.Compare(b.Second, a.Second)
				},
				expected: itertools.Pair[string, int]{"E", 625},
			},
			{
				name: "max key",
				f: func(a itertools.Pair[string, int], b itertools.Pair[string, int]) int {
					return cmp.Compare(a.First, b.First)
				},
				expected: itertools.Pair[string, int]{"E", 625},
			},
			{
				name: "min key",
				f: func(a itertools.Pair[string, int], b itertools.Pair[string, int]) int {
					return cmp.Compare(b.First, a.First)
				},
				expected: itertools.Pair[string, int]{"A", 10000},
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewMapIterator(m).Max(tc.f)
				if tc.expected != result {
					t.Errorf("expected %v as max value, but got %v", tc.expected, result)
				}
			})
		}
	})

}

func TestAsciiIterator(t *testing.T) {
	text := "Hello World!"

	t.Run("collect", func(t *testing.T) {
		result := itertools.NewAsciiIterator(text).Collect()

		if string(result) != text {
			t.Errorf("expected %s, got %s", text, string(result))
		}
	})

	t.Run("count", func(t *testing.T) {
		result := itertools.NewAsciiIterator(text).Count()

		if result != len(text) {
			t.Errorf("expected Count to return %d, but got %d", len(text), result)
		}
	})

	t.Run("drop", func(t *testing.T) {
		i := itertools.NewAsciiIterator(text)
		dropCount := 6

		dropped := i.Drop(dropCount)

		if dropped != dropCount {
			t.Errorf("expected to have %d dropped values, but got %d", dropCount, dropped)
		}

		result := string(i.Collect())
		expected := text[dropCount:]

		if result != expected {
			t.Errorf("expected %s, but got %s", expected, result)
		}

		if dropped = i.Drop(1); dropped != 0 {
			t.Errorf("expected to have nothing remained after iteration, but had %d values", dropped)
		}
	})

	t.Run("drop overflow", func(t *testing.T) {
		i := itertools.NewAsciiIterator(text)
		dropCount := 120

		dropped := i.Drop(dropCount)

		if dropped != len(text) {
			t.Errorf("expected to have %d dropped values, but got %d", len(text), dropped)
		}

		result := string(i.Collect())
		expected := ""

		if result != expected {
			t.Errorf("expected %s, but got %s", expected, result)
		}

		if dropped = i.Drop(1); dropped != 0 {
			t.Errorf("expected to have nothing remained after drop, but had %d values", dropped)
		}
	})

	t.Run("with step", func(t *testing.T) {
		type tcase struct {
			name     string
			step     int
			expected string
		}

		tcases := []tcase{
			{
				name:     "singular step",
				step:     1,
				expected: text,
			},
			{
				name:     "sextuple step",
				step:     6,
				expected: "HW",
			},
			{
				name: "zero step",
				step: 0,
			},
			{
				name: "negative step",
				step: -1,
			},
			{
				name:     "giant step",
				step:     len(text),
				expected: "H",
			},
			{
				name:     "two values",
				step:     len(text) - 1,
				expected: "H!",
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				i := itertools.NewAsciiIterator(text).WithStep(tc.step)
				result := string(i.Collect())

				if result != tc.expected {
					t.Errorf("expected %s, got %s", tc.expected, result)
				}
			})
		}
	})

	t.Run("range", func(t *testing.T) {
		i := itertools.NewAsciiIterator(text)
		reversed := make([]byte, len(text))
		idx := len(text) - 1
		i.Range(func(b byte) bool {
			reversed[idx] = b
			idx--
			return true
		})

		expected := "!dlroW olleH"
		if string(reversed) != expected {
			t.Errorf("expected %s as reversed string with range, but got %s", expected, string(reversed))
		}
	})

	t.Run("filter", func(t *testing.T) {
		type tcase struct {
			name       string
			filterFunc func(byte) bool
			expected   string
		}

		tcases := []tcase{
			{
				name: "uppercase only",
				filterFunc: func(b byte) bool {
					return b >= 'A' && b <= 'Z'
				},
				expected: "HW",
			},
			{
				name: "lowercase only",
				filterFunc: func(b byte) bool {
					return b >= 'a' && b <= 'z'
				},
				expected: "elloorld",
			},
			{
				name: "non-letters",
				filterFunc: func(b byte) bool {
					return (b < 'a' || b > 'z') &&
						(b < 'A' || b > 'Z')
				},
				expected: " !",
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				i := itertools.NewAsciiIterator(text).Filter(tc.filterFunc)
				result := string(i.Collect())

				if result != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, result)
				}
			})
		}
	})

	t.Run("reduce", func(t *testing.T) {
		type tcase struct {
			name         string
			reducer      func(byte, byte) byte
			initialValue byte
			expected     byte
		}

		tcases := []tcase{
			{
				name: "last whitespace",
				reducer: func(acc byte, b byte) byte {
					if unicode.IsSpace(rune(b)) {
						acc = b
					}
					return acc
				},
				initialValue: 0,
				expected:     ' ',
			},
			{
				name: "latest letter in the alphabet",
				reducer: func(acc byte, b byte) byte {
					if b >= 'A' && b <= 'Z' {
						b = b + 'a' - 'A'
					}
					if b >= 'a' && b <= 'z' && b > acc {
						acc = b
					}
					return acc
				},
				initialValue: 0,
				expected:     'w',
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewAsciiIterator(text).
					Reduce(tc.initialValue, tc.reducer)

				if result != tc.expected {
					t.Errorf("expected %c after reduce, got %c", tc.expected, result)
				}
			})
		}
	})

	t.Run("all", func(t *testing.T) {
		type tcase struct {
			name      string
			condition func(byte) bool
			expected  bool
		}

		tcases := []tcase{
			{
				name: "false",
				condition: func(byte) bool {
					return false
				},
				expected: false,
			},
			{
				name: "non-control ascii characters",
				condition: func(b byte) bool {
					return b > 31 && b != 127
				},
				expected: true,
			},
			{
				name: "non-uppercase characters",
				condition: func(b byte) bool {
					return b < 'A' || b > 'Z'
				},
				expected: false,
			},
			{
				name: "plus or minus",
				condition: func(b byte) bool {
					return b == '+' || b == '-'
				},
				expected: false,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewAsciiIterator(text).
					All(tc.condition)
				if result != tc.expected {
					t.Errorf("expected %t in all, but got %t", tc.expected, result)
				}
			})
		}
	})

	t.Run("any", func(t *testing.T) {
		type tcase struct {
			name      string
			condition func(byte) bool
			expected  bool
		}

		tcases := []tcase{
			{
				name: "false",
				condition: func(byte) bool {
					return false
				},
				expected: false,
			},
			{
				name: "non-control ascii characters",
				condition: func(b byte) bool {
					return b > 31 && b != 127
				},
				expected: true,
			},
			{
				name: "non-uppercase characters",
				condition: func(b byte) bool {
					return b < 'A' || b > 'Z'
				},
				expected: true,
			},
			{
				name: "plus or minus",
				condition: func(b byte) bool {
					return b == '+' || b == '-'
				},
				expected: false,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewAsciiIterator(text).
					Any(tc.condition)
				if result != tc.expected {
					t.Errorf("expected %t in any, but got %t", tc.expected, result)
				}
			})
		}
	})

	t.Run("max", func(t *testing.T) {
		type tcase struct {
			name       string
			comparator func(byte, byte) int
			expected   byte
		}

		tcases := []tcase{
			{
				name: "largest uppercase letter",
				comparator: func(a byte, b byte) int {
					isUppercaseA := a >= 'A' && a <= 'Z'
					isUppercaseB := b >= 'A' && b <= 'Z'
					if isUppercaseA && isUppercaseB {
						return cmp.Compare(a, b)
					}
					if isUppercaseA && !isUppercaseB {
						return 1
					}
					return -1
				},
				expected: 'W',
			},
			{
				name: "largest lowercase letter",
				comparator: func(a byte, b byte) int {
					isLowercaseA := a >= 'a' && a <= 'z'
					isLowercaseB := b >= 'a' && b <= 'z'
					if isLowercaseA && isLowercaseB {
						return cmp.Compare(a, b)
					}
					if isLowercaseA && !isLowercaseB {
						return 1
					}
					return -1
				},
				expected: 'r',
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewAsciiIterator(text).Max(tc.comparator)
				if tc.expected != result {
					t.Errorf("expected %c as max value, but got %c", tc.expected, result)
				}
			})
		}

	})
}

func TestUTF8Iterator(t *testing.T) {
	helloWorlds := []string{
		"Hello Wêreld!",
		"Përshendetje Botë!",
		"ሰላም ልዑል!",
		"مرحبا بالعالم!",
		"Բարեւ աշխարհ!",
		"Kaixo Mundua!",
		"Прывітанне Сусвет!",
		"ওহে বিশ্ব!",
		"Здравей свят!",
		"Hola món!",
		"Moni Dziko Lapansi!",
		"你好世界！",
		"Pozdrav svijete!",
		"Ahoj světe!",
		"Hej Verden!",
		"Hallo Wereld!",
		"Hello World!",
		"Tere maailm!",
		"Hei maailma!",
		"Bonjour monde!",
		"Hallo wrâld!",
		"გამარჯობა მსოფლიო!",
		"Hallo Welt!",
		"Γειά σου Κόσμε!",
		"Sannu Duniya!",
		"שלום עולם!",
		"नमस्ते दुनिया!",
		"Helló Világ!",
		"Halló heimur!",
		"Ndewo Ụwa!",
		"Halo Dunia!",
		"Ciao mondo!",
		"こんにちは世界！",
		"Сәлем Әлем!",
		"Салам дүйнө!",
		"Sveika pasaule!",
		"Labas pasauli!",
		"Moien Welt!",
		"Здраво свету!",
		"Hai dunia!",
		"ഹലോ വേൾഡ്!",
		"Сайн уу дэлхий!",
		"မင်္ဂလာပါကမ္ဘာလောက!",
		"नमस्कार संसार!",
		"Hei Verden!",
		"سلام نړی!",
		"سلام دنیا!",
		"Witaj świecie!",
		"Olá Mundo!",
		"ਸਤਿ ਸ੍ਰੀ ਅਕਾਲ ਦੁਨਿਆ!",
		"Salut Lume!",
		"Привет мир!",
		"Hàlo a Shaoghail!",
		"Здраво Свете!",
		"Lefatše Lumela!",
		"හෙලෝ වර්ල්ඩ්!",
		"Pozdravljen svet!",
		"¡Hola Mundo!",
		"Halo Dunya!",
		"Salamu Dunia!",
		"Hej världen!",
		"Салом Ҷаҳон!",
		"สวัสดีชาวโลก!",
		"Selam Dünya!",
		"Привіт Світ!",
		"Salom Dunyo!",
		"Chào thế giới!",
		"Helo Byd!",
		"Molo Lizwe!",
		"העלא וועלט!",
		"Mo ki O Ile Aiye!",
		"Sawubona Mhlaba!",
	}
	text := strings.Join(helloWorlds, "\n")

	t.Run("collect", func(t *testing.T) {
		result := string(itertools.NewUTF8Iterator(text).Collect())

		if result != text {
			t.Errorf("expected %s, got %s", text, result)
		}
	})

	t.Run("count", func(t *testing.T) {
		result := itertools.NewUTF8Iterator(text).Count()

		expected := utf8.RuneCountInString(text)

		if result != expected {
			t.Errorf("expected Count to return %d, but got %d", expected, result)
		}
	})

	t.Run("drop", func(t *testing.T) {
		i := itertools.NewUTF8Iterator(text)
		dropCount := 40

		dropped := i.Drop(dropCount)

		if dropped != dropCount {
			t.Errorf("expected to have %d dropped values, but got %d", dropCount, dropped)
		}

		result := string(i.Collect())
		expected := string([]rune(text)[dropCount:])

		if result != expected {
			t.Errorf("expected %s, but got %s", expected, result)
		}

		if dropped = i.Drop(1); dropped != 0 {
			t.Errorf("expected to have nothing remained after iteration, but had %d values", dropped)
		}
	})

	t.Run("drop overflow", func(t *testing.T) {
		i := itertools.NewUTF8Iterator(text)
		dropCount := 2000

		dropped := i.Drop(dropCount)

		if dropped != utf8.RuneCountInString(text) {
			t.Errorf("expected to have %d dropped values, but got %d", utf8.RuneCountInString(text), dropped)
		}

		result := string(i.Collect())
		expected := ""

		if result != expected {
			t.Errorf("expected %s, but got %s", expected, result)
		}

		if dropped = i.Drop(1); dropped != 0 {
			t.Errorf("expected to have nothing remained after drop, but had %d values", dropped)
		}
	})

	t.Run("with step", func(t *testing.T) {
		type tcase struct {
			name     string
			step     int
			expected string
		}

		tcases := []tcase{
			{
				name:     "singular step",
				step:     1,
				expected: text,
			},
			{
				name: "decimal step",
				step: 10,
				expected: func() string {
					var result []rune
					for i, r := range []rune(text) {
						if i%10 == 0 {
							result = append(result, r)
						}
					}
					return string(result)
				}(),
			},
			{
				name: "zero step",
				step: 0,
			},
			{
				name: "negative step",
				step: -1,
			},
			{
				name: "giant step",
				step: len(text),
				expected: func() string {
					r, _ := utf8.DecodeRuneInString(text)
					return string(r)
				}(),
			},
			{
				name: "two values",
				step: utf8.RuneCountInString(text) - 1,
				expected: func() string {
					first, _ := utf8.DecodeRuneInString(text)
					last, _ := utf8.DecodeLastRuneInString(text)
					return string([]rune{first, last})
				}(),
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Run(tc.name, func(t *testing.T) {
					i := itertools.NewUTF8Iterator(text).WithStep(tc.step)
					result := string(i.Collect())

					if result != tc.expected {
						t.Errorf("expected %q, got %q", tc.expected, result)
					}
				})
			})
		}
	})

	t.Run("range", func(t *testing.T) {
		i := itertools.NewUTF8Iterator(text)
		var idx int
		targetRune := 'ち'
		i.Range(func(r rune) bool {
			if r == targetRune {
				return false
			}
			idx++
			return true
		})

		expected := slices.Index([]rune(text), targetRune)

		if idx != expected {
			t.Errorf("expected %d as index of rune %c, but got %d", expected, targetRune, idx)
		}
	})

	t.Run("filter", func(t *testing.T) {
		type tcase struct {
			name       string
			filterFunc func(rune) bool
			expected   string
		}

		tcases := []tcase{
			{
				name:       "uppercase only",
				filterFunc: unicode.IsUpper,
				expected: func() string {
					var result strings.Builder
					for _, r := range text {
						if unicode.IsUpper(r) {
							result.WriteRune(r)
						}
					}
					return result.String()
				}(),
			},
			{
				name:       "lowercase only",
				filterFunc: unicode.IsLower,
				expected: func() string {
					var result strings.Builder
					for _, r := range text {
						if unicode.IsLower(r) {
							result.WriteRune(r)
						}
					}
					return result.String()
				}(),
			},
			{
				name: "is cyrillic",
				filterFunc: func(r rune) bool {
					return unicode.Is(unicode.Cyrillic, r)
				},
				expected: func() string {
					var result strings.Builder
					for _, r := range text {
						if unicode.Is(unicode.Cyrillic, r) {
							result.WriteRune(r)
						}
					}
					return result.String()
				}(),
			},
			{
				name: "is japanese alphabet",
				filterFunc: func(r rune) bool {
					return unicode.In(r, unicode.Hiragana, unicode.Katakana)
				},
				expected: func() string {
					var result strings.Builder
					for _, r := range text {
						if unicode.In(r, unicode.Hiragana, unicode.Katakana) {
							result.WriteRune(r)
						}
					}
					return result.String()
				}(),
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				i := itertools.NewUTF8Iterator(text).Filter(tc.filterFunc)
				result := string(i.Collect())

				if result != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, result)
				}
			})
		}
	})

	t.Run("reduce", func(t *testing.T) {
		type tcase struct {
			name         string
			reducer      func(rune, rune) rune
			initialValue rune
			expected     rune
		}

		tcases := []tcase{
			{
				name: "last whitespace",
				reducer: func(acc rune, r rune) rune {
					if unicode.IsSpace(r) {
						acc = r
					}
					return acc
				},
				expected: ' ',
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewUTF8Iterator(text).
					Reduce(tc.initialValue, tc.reducer)

				if result != tc.expected {
					t.Errorf("expected %c after reduce, got %c", tc.expected, result)
				}
			})
		}
	})

	t.Run("all", func(t *testing.T) {
		type tcase struct {
			name      string
			condition func(rune) bool
			expected  bool
		}

		tcases := []tcase{
			{
				name: "true",
				condition: func(rune) bool {
					return true
				},
				expected: true,
			},
			{
				name: "does not have error characters",
				condition: func(r rune) bool {
					return r != utf8.RuneError
				},
				expected: true,
			},
			{
				name: "arabic",
				condition: func(r rune) bool {
					return unicode.Is(unicode.Arabic, r)
				},
				expected: false,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewUTF8Iterator(text).
					All(tc.condition)
				if result != tc.expected {
					t.Errorf("expected %t in all, but got %t", tc.expected, result)
				}
			})
		}
	})

	t.Run("any", func(t *testing.T) {
		type tcase struct {
			name      string
			condition func(rune) bool
			expected  bool
		}

		tcases := []tcase{
			{
				name: "true",
				condition: func(rune) bool {
					return true
				},
				expected: true,
			},
			{
				name: "has any error character",
				condition: func(r rune) bool {
					return r == utf8.RuneError
				},
				expected: false,
			},
			{
				name: "arabic",
				condition: func(r rune) bool {
					return unicode.Is(unicode.Arabic, r)
				},
				expected: true,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewUTF8Iterator(text).
					Any(tc.condition)
				if result != tc.expected {
					t.Errorf("expected %t in all, but got %t", tc.expected, result)
				}
			})
		}
	})

	t.Run("max", func(t *testing.T) {
		type tcase struct {
			name       string
			comparator func(rune, rune) int
			expected   rune
		}

		tcases := []tcase{
			{
				name:       "largest",
				comparator: cmp.Compare[rune],
				expected:   slices.Max([]rune(text)),
			},
			{
				name: "largest hebrew letter",
				comparator: func(a rune, b rune) int {
					isHebrewA := unicode.Is(unicode.Hebrew, a)
					isHebrewB := unicode.Is(unicode.Hebrew, b)
					if isHebrewA && isHebrewB {
						return cmp.Compare(a, b)
					}
					if !isHebrewB {
						return 1
					}
					return -1
				},
				expected: 'ש',
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.NewUTF8Iterator(text).Max(tc.comparator)
				if tc.expected != result {
					t.Errorf("expected %q as max value, but got %q", tc.expected, result)
				}
			})
		}
	})

}
