package itertools_test

import (
	"github.com/KSpaceer/itertools"
	"math/rand"
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

	t.Run("sum", func(t *testing.T) {
		result := itertools.Sum(itertools.New(fibonacciYielder(fibonacciLimit)))
		var expected int
		for _, n := range collectedValues {
			expected += n
		}
		if result != expected {
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

	t.Run("batched", func(t *testing.T) {
		i := itertools.Batched(
			itertools.New(fibonacciYielder(fibonacciLimit)),
			5,
		)
		result := i.Collect()

		expected := [][]int{
			collectedValues[:5],
			collectedValues[5:10],
			collectedValues[10:],
		}

		if len(result) != len(expected) {
			t.Errorf("expected %v, got %v", expected, result)
		} else {
			for i := range result {
				if !sliceEqual(expected[i], result[i]) {
					t.Errorf("expected %v, got %v", expected, result)
					break
				}
			}
		}
	})

	t.Run("batched: exact size", func(t *testing.T) {
		i := itertools.Batched(
			itertools.New(fibonacciYielder(fibonacciLimit)),
			3,
		)
		result := i.Collect()

		expected := [][]int{
			collectedValues[:3],
			collectedValues[3:6],
			collectedValues[6:9],
			collectedValues[9:],
		}

		if len(result) != len(expected) {
			t.Errorf("expected %v, got %v", expected, result)
		} else {
			for i := range result {
				if !sliceEqual(expected[i], result[i]) {
					t.Errorf("expected %v, got %v", expected, result)
					break
				}
			}
		}
	})

	t.Run("batched: zero size", func(t *testing.T) {
		i := itertools.Batched(
			itertools.New(fibonacciYielder(fibonacciLimit)),
			0,
		)
		if i.Next() {
			t.Errorf("did not expect elements in zero size batch; elem: %v", i.Elem())
		}
	})

	t.Run("repeat", func(t *testing.T) {
		i := itertools.Repeat("hello").Limit(3)
		result := i.Collect()

		expected := []string{"hello", "hello", "hello"}

		if !sliceEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("cycle", func(t *testing.T) {
		i := itertools.Cycle(
			itertools.New(fibonacciYielder(fibonacciLimit)),
			itertools.WithPrealloc(len(collectedValues)),
		).Limit(len(collectedValues) * 3)

		result := i.Collect()

		var expected []int
		expected = append(expected, collectedValues...)
		expected = append(expected, collectedValues...)
		expected = append(expected, collectedValues...)

		if !sliceEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("empty cycle", func(t *testing.T) {
		i := itertools.Cycle(itertools.New(func() (int, bool) {
			return 0, false
		}))

		if i.Next() {
			t.Errorf("did not expect elements in empty cycle iterator; elem: %d", i.Elem())
		}

		if i.Next() {
			t.Errorf("did not expect elements in empty cycle iterator; elem: %d", i.Elem())
		}
	})

	t.Run("uniq", func(t *testing.T) {
		i := itertools.Uniq(
			itertools.New(fibonacciYielder(fibonacciLimit)),
			itertools.WithPrealloc(len(collectedValues)-1),
		)

		result := i.Collect()

		expected := make([]int, len(collectedValues))
		copy(expected, collectedValues)
		expected = slices.Compact(expected)

		if !sliceEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("uniq func", func(t *testing.T) {
		uniq := func(v int) int { return v % 10 }

		i := itertools.UniqFunc(
			itertools.New(fibonacciYielder(fibonacciLimit)),
			uniq,
			itertools.WithPrealloc(len(collectedValues)),
		)

		result := i.Collect()

		expected := []int{0, 1, 2, 3, 5, 8, 34, 89}

		if !sliceEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("sorted for sorted sequence", func(t *testing.T) {
		i := itertools.Sorted(
			itertools.New(fibonacciYielder(fibonacciLimit)),
			itertools.WithPrealloc(len(collectedValues)),
		)

		result := i.Collect()

		if !sliceEqual(collectedValues, result) {
			t.Errorf("expected %v, got %v", collectedValues, result)
		}
	})

	t.Run("sorted", func(t *testing.T) {
		chaoticValues := slices.Clone(collectedValues)
		rand.Shuffle(len(chaoticValues), func(i, j int) {
			chaoticValues[i], chaoticValues[j] = chaoticValues[j], chaoticValues[i]
		})
		chaoticFibonacci := itertools.NewSliceIterator(chaoticValues)

		i := itertools.Sorted(chaoticFibonacci, itertools.WithPrealloc(len(chaoticValues)))

		result := i.Collect()

		if !sliceEqual(collectedValues, result) {
			t.Errorf("expected %v, got %v", collectedValues, result)
		}
	})
}
