package rx

import "errors"

// Err declares the base error that is joined with every error returned by this package.
// It serves as the foundation for error types in the reactive extensions library,
// allowing for error type checking and custom error creation with consistent taxonomy.
var Err = errors.New("rx")
