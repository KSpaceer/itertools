package itertools_test

import (
	"cmp"
	"fmt"
	"github.com/KSpaceer/itertools"
	"math"
	"slices"
	"strings"
	"sync"
	"unicode"
)

func ExampleNew() {
	arr := [5]int{5, 10, 15, 20, 25}
	var idx int

	// function to manually iterator over slice.
	// It is more comfortable to use itertools.NewSliceIterator in this case
	f := func() (int, bool) {
		if idx >= len(arr) {
			return 0, false
		}
		elem := arr[idx]
		idx++
		return elem, true
	}

	iter := itertools.New(f)

	fmt.Println(iter.Next(), iter.Elem())
	fmt.Println(iter.Next(), iter.Elem())
	fmt.Println(iter.Next(), iter.Elem())
	fmt.Println(iter.Next(), iter.Elem())
	fmt.Println(iter.Next(), iter.Elem())
	fmt.Println(iter.Next())
	// Output:
	// true 5
	// true 10
	// true 15
	// true 20
	// true 25
	// false
}

func ExampleIterator_Count() {
	s := []int{1, 2, 3, 4}
	iter := itertools.NewSliceIterator(s)

	fmt.Println(iter.Count())
	// Output: 4
}

func ExampleIterator_Drop() {
	const skipNextTwo = 0x01

	data := []byte{'H', 'e', 'l', 'l', 'o', 0x01, 'W', 'o', 'r', 'l', 'd'}

	iter := itertools.NewSliceIterator(data)

	for iter.Next() {
		elem := iter.Elem()
		if elem == skipNextTwo {
			dropped := iter.Drop(2)
			fmt.Printf("dropped %d elements\n", dropped)
		} else {
			fmt.Printf("%c\n", elem)
		}
	}
	// Output:
	// H
	// e
	// l
	// l
	// o
	// dropped 2 elements
	// r
	// l
	// d
}

func ExampleIterator_Limit() {
	iter := itertools.NewAsciiIterator("Hello World!").Limit(5)

	for iter.Next() {
		fmt.Printf("%c\n", iter.Elem())
	}
	// Output:
	// H
	// e
	// l
	// l
	// o
}

func ExampleIterator_WithStep() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	iter := itertools.NewSliceIterator(s).WithStep(2)

	for iter.Next() {
		fmt.Println(iter.Elem())
	}
	// Output:
	// 1
	// 3
	// 5
	// 7
	// 9
}

func ExampleIterator_Range() {
	s := []string{"First", "Second", "Third", "Fourth", "Fifth"}

	iter := itertools.NewSliceIterator(s)

	iter.Range(func(s string) bool {
		fmt.Println(s)
		return true
	})

	// Output:
	// First
	// Second
	// Third
	// Fourth
	// Fifth
}

func ExampleIterator_Range_conditional_stop() {
	s := []int{5, 4, 2, 0, -1, -8, -16}

	iter := itertools.NewSliceIterator(s)

	iter.Range(func(n int) bool {
		fmt.Println(n)
		if n < 0 {
			return false
		}
		return true
	})
	// Output:
	// 5
	// 4
	// 2
	// 0
	// -1
}

func ExampleIterator_Filter() {
	s := []int{1, 2, 3, 4, 5, 6, 7}

	oddFilter := func(n int) bool {
		return n%2 == 1
	}

	evenFilter := func(n int) bool {
		return n%2 == 0
	}

	oddIter := itertools.NewSliceIterator(s).Filter(oddFilter)
	evenIter := itertools.NewSliceIterator(s).Filter(evenFilter)

	fmt.Println(oddIter.Collect())
	fmt.Println(evenIter.Collect())

	// Output:
	// [1 3 5 7]
	// [2 4 6]
}

func ExampleIterator_Collect() {
	s := []string{"First", "Second", "Third"}
	idx := len(s) - 1

	// function to iterate over slice in reverse order
	f := func() (string, bool) {
		if idx < 0 {
			return "", false
		}
		elem := s[idx]
		idx--
		return elem, true
	}

	iter := itertools.New(f)

	fmt.Println(iter.Collect())
	// Output:
	// [Third Second First]
}

func ExampleIterator_Reduce() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	sum := func(acc int, n int) int {
		return acc + n
	}

	iter := itertools.NewSliceIterator(s)
	fmt.Println("Sum of slice:", iter.Reduce(0, sum))
	// Output:
	// Sum of slice: 55
}

func ExampleIterator_All_true() {
	s := []int{-4, -3, -2, -1, 0, 1, 2, 3, 4}

	predicate := func(n int) bool {
		if n < 0 {
			n = -n
		}
		return n < 5
	}

	iter := itertools.NewSliceIterator(s)

	fmt.Println(iter.All(predicate))
	// Output:
	// true
}

func ExampleIterator_All_false() {
	s := []int{-4, -3, -2, -1, 0, 1, 2, 3, 4}

	predicate := func(n int) bool {
		return n < 2
	}

	iter := itertools.NewSliceIterator(s)

	fmt.Println(iter.All(predicate))
	// Output:
	// false
}

func ExampleIterator_Any_true() {
	s := []int{-4, -3, -2, -1, 0, 1, 2, 3, 4}

	predicate := func(n int) bool {
		return n < 2
	}

	iter := itertools.NewSliceIterator(s)

	fmt.Println(iter.Any(predicate))
	// Output:
	// true
}

func ExampleIterator_Any_false() {
	s := []int{-4, -3, -2, -1, 0, 1, 2, 3, 4}

	predicate := func(n int) bool {
		if n < 0 {
			n = -n
		}
		return n >= 5
	}

	iter := itertools.NewSliceIterator(s)

	fmt.Println(iter.Any(predicate))
	// Output:
	// false
}

func ExampleIterator_Max() {
	type Person struct {
		Name string
		Age  uint
	}

	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	iter := itertools.NewSliceIterator(people)

	oldest := iter.Max(func(a Person, b Person) int {
		return cmp.Compare(a.Age, b.Age)
	})

	fmt.Printf("Oldest person: %s (age %d)\n", oldest.Name, oldest.Age)
	// Output:
	// Oldest person: John (age 42)
}

func ExampleChain() {
	iter1 := itertools.NewSliceIterator([]byte("Hello"))

	ch := make(chan byte)
	go func() {
		ch <- ' '
		close(ch)
	}()
	iter2 := itertools.NewChanIterator(ch)
	iter3 := itertools.NewAsciiIterator("World!")

	chainedIter := itertools.Chain(iter1, iter2, iter3)

	fmt.Println("chained collected result:", string(chainedIter.Collect()))
	// Output:
	// chained collected result: Hello World!
}

func ExampleZip() {
	names := []string{"Bob", "John", "Michael", "Jenny"}
	ages := []uint{31, 42, 17, 26}

	nameIter := itertools.NewSliceIterator(names)
	ageIter := itertools.NewSliceIterator(ages)

	iter := itertools.Zip(nameIter, ageIter)

	for iter.Next() {
		name, age := iter.Elem().Unpack()
		fmt.Printf("Name: %s ::: Age: %d\n", name, age)
	}
	// Output:
	// Name: Bob ::: Age: 31
	// Name: John ::: Age: 42
	// Name: Michael ::: Age: 17
	// Name: Jenny ::: Age: 26
}

func ExampleMap() {
	const maxPoints = 150
	results := []int{25, 100, 95, 36, 145, 67, 49, 123}

	percentageIter := itertools.Map(
		itertools.NewSliceIterator(results),
		func(result int) int {
			return int(math.Round(float64(result) / maxPoints * 100))
		},
	)

	fmt.Println("Results percentage:", percentageIter.Collect())
	// Output:
	// Results percentage: [17 67 63 24 97 45 33 82]
}

func ExampleMax() {
	s := []int{-80, 23, 0, 54, 13, -39, 45, 33}

	iter := itertools.NewSliceIterator(s)

	fmt.Println("Max value:", itertools.Max(iter))
	// Output:
	// Max value: 54
}

func ExampleMin() {
	s := []int{-80, 23, 0, 54, 13, -39, 45, 33}

	iter := itertools.NewSliceIterator(s)

	fmt.Println("Min value:", itertools.Min(iter))
	// Output:
	// Min value: -80
}

func ExampleFind_found() {
	type Person struct {
		Name string
		Age  uint
	}

	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	target := "Michael"

	iter := itertools.NewSliceIterator(people)

	person, found := itertools.Find(iter, func(person Person) bool {
		return person.Name == target
	})

	if found {
		fmt.Printf("Found person with name %q. Age: %d\n", person.Name, person.Age)
	} else {
		fmt.Printf("Failed to find person with name %q\n", target)
	}
	// Output:
	// Found person with name "Michael". Age: 17
}

func ExampleFind_not_found() {
	type Person struct {
		Name string
		Age  uint
	}

	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	target := "Mike"

	iter := itertools.NewSliceIterator(people)

	person, found := itertools.Find(iter, func(person Person) bool {
		return person.Name == target
	})

	if found {
		fmt.Printf("Found person with name %q. Age: %d\n", person.Name, person.Age)
	} else {
		fmt.Printf("Failed to find person with name %q\n", target)
	}
	// Output:
	// Failed to find person with name "Mike"
}

func ExampleEnumerate() {
	type Person struct {
		Name string
		Age  uint
	}

	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	iter := itertools.Enumerate(itertools.NewSliceIterator(people))

	for iter.Next() {
		person, i := iter.Elem().Unpack()
		fmt.Printf("Index: %d ||| Name: %s ||| Age: %d\n", i, person.Name, person.Age)
	}
	// Output:
	// Index: 0 ||| Name: Bob ||| Age: 31
	// Index: 1 ||| Name: John ||| Age: 42
	// Index: 2 ||| Name: Michael ||| Age: 17
	// Index: 3 ||| Name: Jenny ||| Age: 26
}

func ExampleBatched() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	iter := itertools.Batched(itertools.NewSliceIterator(s), 2)

	for iter.Next() {
		fmt.Println("got batch:", iter.Elem())
	}
	// Output:
	// got batch: [1 2]
	// got batch: [3 4]
	// got batch: [5 6]
	// got batch: [7 8]
	// got batch: [9]
}

func ExampleRepeat() {
	iter := itertools.Repeat("HELLO")

	for i := 0; i < 5; i++ {
		iter.Next()
		fmt.Println(iter.Elem())
	}
	// Output:
	// HELLO
	// HELLO
	// HELLO
	// HELLO
	// HELLO
}

func ExampleCycle() {
	s := []int{1, 2, 3}

	iter := itertools.Cycle(itertools.NewSliceIterator(s))

	for i := 0; i < 10; i++ {
		iter.Next()
		fmt.Println(iter.Elem())
	}
	// Output:
	// 1
	// 2
	// 3
	// 1
	// 2
	// 3
	// 1
	// 2
	// 3
	// 1
}

func ExampleUniq() {
	s := []int{1, 2, 3, 3, 2, 4, 2, 5, 4, 5, 1}

	iter := itertools.Uniq(itertools.NewSliceIterator(s))

	for iter.Next() {
		fmt.Println(iter.Elem())
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleUniqFunc() {
	words := []string{"Hello", "HELLO", "World", "HeLlO", "world", "HEllO", "WoRlD"}

	iter := itertools.UniqFunc(
		itertools.NewSliceIterator(words),
		strings.ToLower,
	)

	for iter.Next() {
		fmt.Println(iter.Elem())
	}
	// Output:
	// Hello
	// World
}

func ExampleSorted() {
	s := []int{5, 8, 6, 4, 3, 7, 2, 1}

	iter := itertools.Sorted(itertools.NewSliceIterator(s))

	for iter.Next() {
		fmt.Println(iter.Elem())
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
}

func ExampleNewSliceIterator() {
	s := []int{1, 2, 3, 4, 5}

	iter := itertools.NewSliceIterator(s)

	for iter.Next() {
		fmt.Println(iter.Elem())
	}
	fmt.Println("is finished:", !iter.Next())
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// is finished: true
}

func ExampleNewChanIterator() {
	ch := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				ch <- j
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	iter := itertools.NewChanIterator(ch)
	result := iter.Collect()
	slices.Sort(result)
	fmt.Println(result)
	// Output:
	// [0 0 0 1 1 1 2 2 2 3 3 3 4 4 4]
}

func ExampleNewMapIterator() {
	months := map[int]string{
		1:  "January",
		2:  "February",
		3:  "March",
		4:  "April",
		5:  "May",
		6:  "June",
		7:  "July",
		8:  "August",
		9:  "September",
		10: "October",
		11: "November",
		12: "December",
	}

	iter := itertools.NewMapIterator(months)

	summerMonths := []string{"June", "July", "August"}

	iter = iter.Filter(func(p itertools.Pair[int, string]) bool {
		return slices.Contains(summerMonths, p.Second)
	})

	// iterating over map keys (months numbers) rather than entire key-value pairs
	numbersIter := itertools.Map(iter, func(p itertools.Pair[int, string]) int {
		return p.First
	})

	summerMonthsNumbers := numbersIter.Collect()
	slices.Sort(summerMonthsNumbers)
	fmt.Println("summer months numbers:", summerMonthsNumbers)
	// Output:
	// summer months numbers: [6 7 8]
}

func ExampleNewAsciiIterator() {
	s := "Hello, World!"

	iter := itertools.NewAsciiIterator(s)

	iter = iter.Filter(func(b byte) bool {
		return unicode.IsPunct(rune(b))
	})

	fmt.Println("phrase punctuation signs:", string(iter.Collect()))
	// Output:
	// phrase punctuation signs: ,!
}

func ExampleNewUTF8Iterator() {
	s := "Hello, 世界"

	iter := itertools.NewUTF8Iterator(s)

	iter = iter.Filter(func(r rune) bool {
		return unicode.Is(unicode.Han, r)
	})

	fmt.Println("chinese hieroglyphs:", string(iter.Collect()))
	// Output:
	// chinese hieroglyphs: 世界
}
