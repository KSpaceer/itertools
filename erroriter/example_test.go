package erroriter_test

import (
	"fmt"
	"github.com/KSpaceer/itertools"
	"github.com/KSpaceer/itertools/erroriter"
	"strconv"
)

func ExampleErrorIterator_Result() {
	iter := erroriter.New(func() (int, error) {
		return strconv.Atoi("16")
	})

	iter.Next()
	result, err := iter.Result()
	if err != nil {
		fmt.Println("got error")
	} else {
		fmt.Println(result)
	}

	iter = erroriter.New(func() (int, error) {
		return strconv.Atoi("xadqwe")
	})

	iter.Next()
	result, err = iter.Result()
	if err != nil {
		fmt.Println("got error")
	} else {
		fmt.Println(result)
	}
	// Output:
	// 16
	// got error
}

func ExampleErrorIterator_CollectUntilError() {
	data := []int{1, 2, 3, 4, 5}
	var idx int

	iter := erroriter.New(func() (int, error) {
		if idx >= len(data) {
			return 0, erroriter.ErrIterationStop
		}
		v := data[idx]
		idx++
		return v, nil
	})

	result, err := iter.CollectUntilError()
	if err != nil {
		fmt.Println("got error")
	} else {
		fmt.Println(result)
	}

	sdata := []string{"1", "2", "adqweqw", "4", "5"}
	idx = 0

	iter = erroriter.New(func() (int, error) {
		if idx >= len(sdata) {
			return 0, erroriter.ErrIterationStop
		}
		s := sdata[idx]
		idx++
		return strconv.Atoi(s)
	})

	result, err = iter.CollectUntilError()
	if err != nil {
		fmt.Println("got error")
	} else {
		fmt.Println(result)
	}
	// Output:
	// [1 2 3 4 5]
	// got error
}

func ExampleMap() {
	data := []string{"1", "wdqe", "3", "4", "qwcqwqdq"}

	iter := erroriter.Map(
		itertools.NewSliceIterator(data),
		strconv.Atoi,
	)

	for iter.Next() {
		v, err := iter.Result()
		if err != nil {
			fmt.Println("got error")
		} else {
			fmt.Println(v)
		}
	}

	// Output:
	// 1
	// got error
	// 3
	// 4
	// got error
}
