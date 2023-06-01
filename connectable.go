package x

// Connectable[T] is an Observable[T] with a Connect method. So it has a both a
// Subscribe method and a Connect method. The Connect method is used to
// instruct the Connectable[T] to subscribe to its source. The Subscribe method
// is used to subscribe to the Connectable[T] itself.
type Connectable[T any] struct {
	Observable[T]
	Connect
}
