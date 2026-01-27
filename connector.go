package rx

// Connector provides the Connect method for a Connectable[T].
type Connector func(Scheduler, Subscriber)
