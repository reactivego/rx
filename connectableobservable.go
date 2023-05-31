package x

// ConnectableObservable[T] is an Observable[T] with a Connect method.
// So it has a Subscribe method and a Connect method.
// The Connect method is used to instruct the ConnectableObservable[T] to
// subscribe to its source. The Subscribe method is used to subscribe to the
// ConnectableObservable[T] itself.
type ConnectableObservable[T any] struct {
	Connectable
	Observable[T]
}
