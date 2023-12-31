package erroriter_test

import (
	"github.com/KSpaceer/itertools/erroriter"
	"strconv"
	"testing"
)

func TestErrorIterator(t *testing.T) {
	t.Run("collect correct", func(t *testing.T) {
		source := []string{"1", "2", "-2", "4", "5"}
		var idx int
		f := func() (int, error) {
			if idx >= len(source) {
				return 0, erroriter.ErrIterationStop
			}
			v, err := strconv.Atoi(source[idx])
			idx++
			return v, err
		}

		collected := []int{1, 2, -2, 4, 5}

		result, err := erroriter.New(f).CollectUntilError()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		} else {
			if len(result) != len(collected) {
				t.Errorf("expected %v, got %v", collected, result)
			}

			for i := range collected {
				if collected[i] != result[i] {
					t.Errorf("expected %v, got %v", collected, result)
				}
			}
		}
	})

	t.Run("collect with error", func(t *testing.T) {
		source := []string{"1", "2", "-dqwqew2", "4", "5"}
		var idx int
		f := func() (int, error) {
			if idx >= len(source) {
				return 0, erroriter.ErrIterationStop
			}
			v, err := strconv.Atoi(source[idx])
			idx++
			return v, err
		}

		result, err := erroriter.New(f).CollectUntilError()
		if err == nil {
			t.Errorf("expected error, got nil. Collected result: %v", result)
		}
	})
}
