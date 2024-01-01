# itertools ![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/KSpaceer/itertools/itertools.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/KSpaceer/itertools)](https://goreportcard.com/report/github.com/KSpaceer/itertools) ![Codecov](https://img.shields.io/codecov/c/github/KSpaceer/itertools) [![Go Reference](https://pkg.go.dev/badge/github.com/KSpaceer/itertools.svg)](https://pkg.go.dev/github.com/KSpaceer/itertools)

```itertools``` is a Go library introducing ```Iterator``` type to iterate over elements of collections/sequences
with many methods similar to Rust or Python iterators.

## Documentation

Documentation for the library packages and examples are available with [GoDoc](https://pkg.go.dev/github.com/KSpaceer/itertools).

## Installation

```go get github.com/KSpaceer/itertools@latest```

## Example

```go
data := []int{8, 1, 2, 3, 8, 7, 4, 5, 1, 6, 7, 8}

// finding max value in data
maxValue := itertools.Max(itertools.NewSliceIterator(data))

// transform number so it is in range [0, 1]
normalizeFunc := func(n int) float64 {
    return float64(n) / float64(maxValue)
}

// iterator for normalized values
normalizedIter := itertools.Map(
    itertools.NewSliceIterator(data),
    normalizeFunc,
)

// iterator for original slice
originalIter := itertools.NewSliceIterator(data)

// zipping normalized iterator with the original one to output values together
iter := itertools.Zip(originalIter, normalizedIter)

for iter.Next() {
    value, normalizedValue := iter.Elem().Unpack()
    fmt.Printf("original value: %d\tnormalized value: %.2f\n", value, normalizedValue)
}
```

[Run with Playground](https://go.dev/play/p/hbnUTW1oZWK).

See more examples in [docs](https://pkg.go.dev/github.com/KSpaceer/itertools).

