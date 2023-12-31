package itertools

type allocOptions struct {
	preallocSize int
}

// AllocationOption allows to manipulate allocations in iteration methods/functions.
type AllocationOption func(options *allocOptions)

// WithPrealloc sets preallocation size (capacity) for allocated buffers/slices.
func WithPrealloc(prealloc int) AllocationOption {
	return func(o *allocOptions) {
		o.preallocSize = prealloc
	}
}
