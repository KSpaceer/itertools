package itertools_test

import (
	"fmt"
	"github.com/KSpaceer/itertools"
	"math"
)

func Example_mapReduce() {
	data := []int{6, 10, 7, 12, 6, 14, 8, 13, 10, 14}

	sum := itertools.Sum(itertools.NewSliceIterator(data))
	avg := float64(sum) / float64(len(data))

	stddev := itertools.
		Map(
			itertools.NewSliceIterator(data),
			func(n int) float64 { return float64(n) },
		).
		Reduce(0, func(acc float64, elem float64) float64 {
			return acc + (elem-avg)*(elem-avg)
		})

	stddev = math.Sqrt(stddev / float64(len(data)-1))
	fmt.Println("data:", data)
	fmt.Printf("standard deviation: %.2f", stddev)
	// Output:
	// data: [6 10 7 12 6 14 8 13 10 14]
	// standard deviation: 3.16
}
