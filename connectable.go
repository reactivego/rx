package x

// Connectable provides the Connect method for a ConnectableObservable[T].
type Connectable func(Scheduler, Subscriber)
