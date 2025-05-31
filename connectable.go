package rx

// Connectable[T] is an Observable[T] that provides delayed connection to its source.
// It combines both Observable and Connector interfaces:
//   - Subscribe: Allows consumers to register for notifications from this Observable
//   - Connect: Triggers the actual subscription to the underlying source Observable
//
// The key feature of Connectable[T] is that it doesn't subscribe to its source
// until the Connect method is explicitly called, allowing multiple observers to
// subscribe before the source begins emitting items (multicast behavior).
type Connectable[T any] struct {
	Observable[T]
	Connector
}
