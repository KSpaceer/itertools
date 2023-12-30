package itertools_test

import (
	"github.com/KSpaceer/itertools"
	"slices"
	"strconv"
	"testing"
)

func TestFibonacciIterator_Iterfuncs(t *testing.T) {
	const fibonacciLimit = 100
	collectedValues := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}

	t.Run("chain", func(t *testing.T) {
		i1 := itertools.New(fibonacciYielder(fibonacciLimit))
		i2 := itertools.New(fibonacciYielder(fibonacciLimit))
		i3 := itertools.New(fibonacciYielder(fibonacciLimit))

		result := itertools.Chain(i1, i2, i3).Collect()

		expected := make([]int, 0, len(collectedValues)*3)
		expected = append(expected, collectedValues...)
		expected = append(expected, collectedValues...)
		expected = append(expected, collectedValues...)

		if cmpResult := sliceCmp(expected, result); cmpResult != 0 {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("zip", func(t *testing.T) {
		var count int
		countIter := itertools.New(func() (int, bool) {
			if count >= len(collectedValues) {
				return 0, false
			}
			result := count
			count++
			return result, true
		})
		fibIter := itertools.New(fibonacciYielder(fibonacciLimit))

		result := itertools.Zip(countIter, fibIter).Collect()

		expected := make([]itertools.Pair[int, int], 0, len(collectedValues))
		for i, n := range collectedValues {
			expected = append(expected, itertools.Pair[int, int]{i, n})
		}
		if !sliceEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("non-matching length zip", func(t *testing.T) {
		var count int
		countIter := itertools.New(func() (int, bool) {
			if count >= 1 {
				return 0, false
			}
			result := count
			count++
			return result, true
		})
		fibIter := itertools.New(fibonacciYielder(fibonacciLimit))

		result := itertools.Zip(fibIter, countIter).Collect()

		expected := make([]itertools.Pair[int, int], 0, len(collectedValues))
		for i, n := range collectedValues[:1] {
			expected = append(expected, itertools.Pair[int, int]{n, i})
		}
		if !sliceEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("map", func(t *testing.T) {
		t.Run("halved", func(t *testing.T) {
			i := itertools.New(fibonacciYielder(fibonacciLimit))
			result := itertools.Map(i, func(n int) float64 { return float64(n) / 2 }).Collect()

			expected := make([]float64, 0, len(collectedValues))

			for _, n := range collectedValues {
				expected = append(expected, float64(n)/2)
			}

			if !sliceEqual(expected, result) {
				t.Errorf("expected %v, got %v", expected, result)
			}
		})
		t.Run("stringified", func(t *testing.T) {
			i := itertools.New(fibonacciYielder(fibonacciLimit))
			result := itertools.Map(i, func(n int) string { return strconv.Itoa(n) }).Collect()

			expected := make([]string, 0, len(collectedValues))

			for _, n := range collectedValues {
				expected = append(expected, strconv.Itoa(n))
			}

			if !sliceEqual(expected, result) {
				t.Errorf("expected %v, got %v", expected, result)
			}
		})
	})

	t.Run("max", func(t *testing.T) {
		result := itertools.Max(itertools.New(fibonacciYielder(fibonacciLimit)))

		if expected := slices.Max(collectedValues); result != expected {
			t.Errorf("expected %d, got %d", expected, result)
		}
	})

	t.Run("min", func(t *testing.T) {
		result := itertools.Min(itertools.New(fibonacciYielder(fibonacciLimit)))

		if expected := slices.Min(collectedValues); result != expected {
			t.Errorf("expected %d, got %d", expected, result)
		}
	})

	t.Run("find", func(t *testing.T) {
		t.Run("found", func(t *testing.T) {
			result, ok := itertools.Find(
				itertools.New(fibonacciYielder(fibonacciLimit)),
				func(n int) bool {
					return n > 0 && n%3 == 0 && n%7 == 0
				},
			)
			if !ok {
				t.Errorf("expected to find some value")
			} else if result != 21 {
				t.Errorf("expected %d, got %d", 21, result)
			}
		})
		t.Run("not found", func(t *testing.T) {
			result, ok := itertools.Find(
				itertools.New(fibonacciYielder(fibonacciLimit)),
				func(n int) bool {
					return n < 0
				},
			)
			if ok {
				t.Errorf("did not expect to find value; found %d", result)
			}
		})
	})

	t.Run("enumerate", func(t *testing.T) {
		i := itertools.Enumerate(
			itertools.New(fibonacciYielder(fibonacciLimit)),
		)
		result := i.Collect()

		expected := make([]itertools.Enumeration[int], 0, len(collectedValues))
		for i, n := range collectedValues {
			expected = append(expected, itertools.Enumeration[int]{
				First:  n,
				Second: i,
			})
		}

		if !sliceEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}
