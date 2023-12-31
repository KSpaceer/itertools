package erroriter_test

import (
	"errors"
	"github.com/KSpaceer/itertools"
	"github.com/KSpaceer/itertools/erroriter"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	s := []string{"1", "2", "-2", "xnqwe", "5"}
	collected := []itertools.Pair[int, error]{
		{1, nil},
		{2, nil},
		{-2, nil},
		{0, strconv.ErrSyntax},
		{5, nil},
	}

	result := erroriter.Map(
		itertools.NewSliceIterator(s),
		strconv.Atoi,
	).Collect()

	if len(result) != len(collected) {
		t.Errorf("expected %v, got %v", collected, result)
	}

	for i := range collected {
		if !errors.Is(result[i].Second, collected[i].Second) || result[i].First != collected[i].First {
			t.Errorf("expected %v, got %v", collected, result)
		}
	}
}
