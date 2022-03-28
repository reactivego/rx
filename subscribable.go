package observable

// Subscribable is the interface shared by Subscription and Subscriber.
type Subscribable interface {
	// Subscribed returns true when the Subscribable is currently active.
	Subscribed() bool

	// Unsubscribe will do nothing if it not an active. If it is still active
	// however, it will be changed to canceled. Subsequently, it will call
	// Unsubscribe down the subscriber tree on all children, along with all
	// methods added through OnUnsubscribe on those childern. On Unsubscribe
	// any call to the Wait method on a Subsvcription will return the error
	// ErrUnsubscribed.
	Unsubscribe()

	// Canceled returns true when the Subscribable has been canceled.
	Canceled() bool
}
