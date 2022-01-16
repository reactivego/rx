/*
Only filters the value stream of an observable and lets only the
values of a specific type pass.

So in case of OnlyInt32 it will only let int32 values pass through.
*/
package Only

// Size is recognized by jig directly because refer type and actual type are both "Size"
type Size struct{ width, height float64 }

// point will be used in a []point
type point struct{ x, y float64 }
