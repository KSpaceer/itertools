package itertools_test

import (
	"cmp"
	"fmt"
	"github.com/KSpaceer/itertools"
	"slices"
	"strings"
	"testing"
)

func TestBasicIterator(t *testing.T) {
	t.Run("empty iterator", func(t *testing.T) {
		i := itertools.New(func() (int, bool) { return 0, false })
		if i.Next() {
			t.Errorf("expected iterator to be empty, but has element: %d", i.Elem())
		}
	})
	t.Run("calls after stop", func(t *testing.T) {
		i := itertools.New(func() (int, bool) { return 0, false })
		if i.Next() {
			t.Errorf("expected iterator to be empty, but has element: %d", i.Elem())
		}

		_ = i.Elem()
		i.Next()
		_ = i.Elem()
		i.Next()
		_ = i.Elem()
	})
	t.Run("empty max", func(t *testing.T) {
		i := itertools.New(func() (int32, bool) { return 0, false })
		result := i.Max(cmp.Compare[int32])
		var expected int32

		if result != expected {
			t.Errorf("expected zero value as max for empty iterator, but got %d", result)
		}
	})
	t.Run("zero limit", func(t *testing.T) {
		i := itertools.New(func() (int32, bool) { return 1, true }).Limit(0)
		if i.Next() {
			t.Errorf("expected iterator to be empty, but has element: %d", i.Elem())
		}
	})
	t.Run("negative limit", func(t *testing.T) {
		i := itertools.New(func() (int32, bool) { return 1, true }).Limit(-15)
		if i.Next() {
			t.Errorf("expected iterator to be empty, but has element: %d", i.Elem())
		}
	})
}

func TestFibonacciIterator(t *testing.T) {
	const fibonacciLimit = 10000
	collectedValues := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89,
		144, 233, 377, 610, 987, 1597, 2584, 4181, 6765}

	t.Run("iteration", func(t *testing.T) {
		i := itertools.New(fibonacciYielder(fibonacciLimit))
		var result []int
		for i.Next() {
			result = append(result, i.Elem())
		}

		if cmpResult := sliceCmp(collectedValues, result); cmpResult != 0 {
			t.Errorf("expected %v, got %v", collectedValues, result)
		}
	})
	t.Run("collect slice", func(t *testing.T) {
		i := itertools.New(fibonacciYielder(fibonacciLimit))
		result := i.Collect(itertools.WithPrealloc(len(collectedValues)))

		if cmpResult := sliceCmp(collectedValues, result); cmpResult != 0 {
			t.Errorf("expected %v, got %v", collectedValues, result)
		}
	})
	t.Run("count", func(t *testing.T) {
		i := itertools.New(fibonacciYielder(fibonacciLimit))
		if result := i.Count(); result != len(collectedValues) {
			t.Errorf("expected Count to return %d, but got %d", len(collectedValues), result)
		}
	})
	t.Run("drop", func(t *testing.T) {
		i := itertools.New(fibonacciYielder(fibonacciLimit))
		dropCount := 12

		dropped := i.Drop(dropCount)

		if dropped != dropCount {
			t.Errorf("expected to have %d dropped values, but got %d", dropCount, dropped)
		}

		result := i.Collect()
		expected := collectedValues[dropCount:]

		if cmpResult := sliceCmp(expected, result); cmpResult != 0 {
			t.Errorf("expected %v, but got %v", expected, result)
		}

		if dropped = i.Drop(1); dropped != 0 {
			t.Errorf("expected to have nothing remained after iteration, but had %d values", dropped)
		}
	})
	t.Run("drop overflow", func(t *testing.T) {
		i := itertools.New(fibonacciYielder(fibonacciLimit))
		dropCount := 120

		dropped := i.Drop(dropCount)

		if dropped != len(collectedValues) {
			t.Errorf("expected to have %d dropped values, but got %d", len(collectedValues), dropped)
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
	t.Run("limit", func(t *testing.T) {
		const limit = 5
		i := itertools.New(fibonacciYielder(fibonacciLimit)).Limit(limit)

		result := i.Collect()
		expected := collectedValues[:limit]

		if cmpResult := sliceCmp(expected, result); cmpResult != 0 {
			t.Errorf("expected %v, but got %v", expected, result)
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
				expected: collectedValues,
			},
			{
				name: "double step",
				step: 2,
				expected: []int{0, 1, 3, 8, 21, 55,
					144, 377, 987, 2584, 6765},
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
				step:     len(collectedValues),
				expected: []int{0},
			},
			{
				name:     "two values",
				step:     len(collectedValues) - 1,
				expected: []int{0, 6765},
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				i := itertools.New(fibonacciYielder(fibonacciLimit)).WithStep(tc.step)
				result := i.Collect()

				if cmpResult := sliceCmp(tc.expected, result); cmpResult != 0 {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			})
		}
	})
	t.Run("range", func(t *testing.T) {
		i := itertools.New(fibonacciYielder(fibonacciLimit))
		var sb strings.Builder
		i.Range(func(n int) bool {
			sb.WriteString(fmt.Sprint(n))
			return true
		})

		var expectedSB strings.Builder
		for _, n := range collectedValues {
			expectedSB.WriteString(fmt.Sprint(n))
		}

		result := sb.String()
		expected := expectedSB.String()

		if result != expected {
			t.Errorf("expected %s as concatenated string with range, but got %s", expected, result)
		}
	})
	t.Run("limited range", func(t *testing.T) {
		i := itertools.New(fibonacciYielder(fibonacciLimit))
		limit := 5
		var (
			sb  strings.Builder
			idx int
		)
		i.Range(func(n int) bool {
			if idx >= limit {
				return false
			}
			idx++
			sb.WriteString(fmt.Sprint(n))
			return true
		})

		var expectedSB strings.Builder
		for _, n := range collectedValues[:limit] {
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
				expected: []int{1, 1, 3, 5, 13, 21, 55, 89,
					233, 377, 987, 1597, 4181, 6765},
			},
			{
				name: "even only",
				filterFunc: func(n int) bool {
					return n%2 == 0
				},
				expected: []int{0, 2, 8, 34, 144, 610, 2584},
			},
			{
				name: "divisible by 3 only",
				filterFunc: func(n int) bool {
					return n%3 == 0
				},
				expected: []int{0, 3, 21, 144, 987, 6765},
			},
			{
				name: "divisible by 5 only",
				filterFunc: func(n int) bool {
					return n%5 == 0
				},
				expected: []int{0, 5, 55, 610, 6765},
			},
			{
				name: "pass all",
				filterFunc: func(int) bool {
					return true
				},
				expected: collectedValues,
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
				i := itertools.New(fibonacciYielder(fibonacciLimit)).Filter(tc.filterFunc)
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
					var s int
					for _, n := range collectedValues {
						s += n
					}
					return s
				}(),
			},
			{
				name: "mul",
				reducer: func(acc int, n int) int {
					return acc * n
				},
				initialValue: 1,
				expected:     0,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.New(fibonacciYielder(fibonacciLimit)).
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
				name: "does not ends with 6",
				condition: func(n int) bool {
					return n%10 != 6
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
				name: "less than 1000",
				condition: func(n int) bool {
					return n < 1000
				},
				expected: false,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.New(fibonacciYielder(fibonacciLimit)).
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
				name: "ends with 6",
				condition: func(n int) bool {
					return n%10 == 6
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
				result := itertools.New(fibonacciYielder(fibonacciLimit)).
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
				expected: 6765,
			},
			{
				name: "minimum",
				f: func(a int, b int) int {
					return -cmp.Compare(a, b)
				},
				expected: 0,
			},
			{
				name: "nearest to 3000",
				f: func(a int, b int) int {
					adiff := 3000 - a
					if adiff < 0 {
						adiff = -adiff
					}
					bdiff := 3000 - b
					if bdiff < 0 {
						bdiff = -bdiff
					}
					return cmp.Compare(bdiff, adiff)
				},
				expected: 2584,
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				result := itertools.New(fibonacciYielder(fibonacciLimit)).Max(tc.f)
				if tc.expected != result {
					t.Errorf("expected %d as max value, but got %d", tc.expected, result)
				}
			})
		}
	})

	t.Run("sorted by", func(t *testing.T) {
		type tcase struct {
			name     string
			f        func(int, int) int
			expected []int
		}

		tcases := []tcase{
			{
				name:     "asc",
				f:        cmp.Compare[int],
				expected: collectedValues,
			},
			{
				name: "desc",
				f: func(a int, b int) int {
					return cmp.Compare(b, a)
				},
				expected: func() []int {
					s := slices.Clone(collectedValues)
					slices.Reverse(s)
					return s
				}(),
			},
		}

		for _, tc := range tcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				iter := itertools.New(fibonacciYielder(fibonacciLimit)).SortedBy(tc.f)
				result := iter.Collect()
				if !sliceEqual(tc.expected, result) {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			})
		}
	})

}

func fibonacciYielder(limit int) func() (int, bool) {
	var state int
	a, b := 0, 1
	return func() (int, bool) {
		var result int
		switch state {
		case 0:
			result = a
			state = 1
		case 1:
			result = b
			state = 2
		default:
			result = a + b
			a, b = b, result
		}
		if result > limit {
			return 0, false
		}
		return result, true
	}
}

func sliceCmp[S ~[]T, T cmp.Ordered](a, b S) int {
	if len(a) != len(b) {
		return cmp.Compare(len(a), len(b))
	}
	for i := range a {
		if cmpResult := cmp.Compare(a[i], b[i]); cmpResult != 0 {
			return cmpResult
		}
	}
	return 0
}

func sliceEqual[S ~[]T, T comparable](a, b S) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
