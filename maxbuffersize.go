package rx

// MaxBufferSizeOption is a function type used for configuring the maximum buffer size
// of an observable stream.
type MaxBufferSizeOption = func(*int)

// WithMaxBufferSize creates a MaxBufferSizeOption that sets the maximum buffer size to n.
// This option is typically used when creating new observables to control memory usage.
func WithMaxBufferSize(n int) MaxBufferSizeOption {
	return func(MaxBufferSize *int) {
		*MaxBufferSize = n
	}
}
