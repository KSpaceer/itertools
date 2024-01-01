package itertools_test

import (
	"fmt"
	"github.com/KSpaceer/itertools"
	"sync"
	"time"
)

func Example_set() {
	values := map[int]struct{}{}
	for i := 0; i < 1_000_000; i++ {
		values[i%40000] = struct{}{}
	}

	const batchSize = 2500

	iter := itertools.
		// batching incoming values into batches with size batchSize (2500)
		Batched(itertools.NewMapKeysIterator(values).
			// keeping only numbers divisible by 4
			Filter(func(n int) bool { return n%4 == 0 }),
			batchSize)

	var wg sync.WaitGroup

	for iter.Next() {
		batch := iter.Elem()
		wg.Add(1)
		go func() {
			defer wg.Done()
			process(batch)
		}()
	}

	wg.Wait()
	// Output:
	// processed 2500 items
	// processed 2500 items
	// processed 2500 items
	// processed 2500 items
}

func process(values []int) {
	// long processing imitation
	time.Sleep(time.Duration(len(values)) * time.Microsecond)
	fmt.Printf("processed %d items\n", len(values))
}
