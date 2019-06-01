// Package subscriber provides a subscription tree implementation.
//
// The New function creates a subscription and returns a Subscriber interface.
// It allows a publisher and subscribing client to communicate about a
// subscription. The publisher can see whether the client is still subscribed
// by polling the Closed method. The client can control the subscription
// by either calling Unsubscribe itself or when it has defered subscription
// management to other code call Wait to wait for other code to cancel the
// subscription.
//
// Once the client has called Unsubscribe on the subscription, that subscription
// will then call Unsubscribe on all of its children created through the Add
// method. After that it will become lifeless, meaning no events will ever be
// triggered by it anymore.
//
// In case the publisher decides it wants to stop publishing, it MUST NEVER
// call Unsubscribe itself on a subscriber, but should indicate the desire
// through other (out-of-band) means to the subscribing client who must then
// call Unsubscribe itself.
//
// The client that receives a Subcriber interface can call Wait to wait for
// the subscription to be canceled. That may be done by the code that e.g.
// created the subscription in responese to some external event.
//
// The client can also keep a reference to the subscription and call
// Unsubscribe itself to indicate to the publisher it is no longer interested
// in receiving data associated with the subscription.
//
// The implementation is designed to be used from concurrently running
// goroutines. It uses WaitGroups, Mutexes and atomic reference counting.
package subscriber

// Subscription is an interface that allows code to monitor and control a
// subscription it received.
type Subscription interface {
	// Unsubscribe will cancel the subscription (when one is active).
	// Subsequently it will then call Unsubscribe on all child subscriptions
	// added through Add. After a call to Unsubscribe returns, calling
	// Closed on the same interface or any of its child subscriptions
	// will return true. Unsubscribe can be safely called on a closed
	// subscription and performs no operation.
	Unsubscribe()

	// Closed returns true when the subscription has been canceled. There
	// is not need to check if the subscription is canceled before calling
	// Unsubscribe. This is an alias for the Canceled() method.
	Closed() bool

	// Canceled returns true when the subscription has been canceled. There
	// is not need to check if the subscription is canceled before calling
	// Unsubscribe.
	Canceled() bool

	// Wait will block the calling goroutine and wait for the Unsubscribe
	// method to be called on this subscription. Calling Wait on a
	// subscription that has already been closed will return immediately.
	Wait()
}

// Subscriber embeds a Subscription interface. Additionally the Add method
// allows for creating a child subscription. Calling Unsubscribe will close
// the current subscription but will not propagate to parent or children.
// Calling Unsubscribe will traverse recursively all child subcriptions and
// call Unsubscribe on them before settings the subsription state to lifeless.
type Subscriber interface {
	// Subscription is embedded in a Subscriber to make it act like one.
	Subscription

	// Add will create and return a new child Subscriber with the given
	// callback function. The callback will be called when either the
	// Unsubscribe of the parent or of the returned child subscriber is called.
	// Calling the Unsubscribe method on the child will NOT propagate to the
	// parent! The Unsubscribe will start calling callbacks only after it has
	// set the subscription state to canceled. Even if you call Unsubscribe
	// multiple times, callbacks will only be invoked once.
	Add(func()) Subscriber

	// AddChild will create and return a new child Subscriber. Calling the
	// Unsubscribe method on the child will NOT propagate to the parent!
	// The Unsubscribe will start calling callbacks only after it has
	// set the subscription state to canceled. Even if you call Unsubscribe
	// multiple times, callbacks will only be invoked once.
	AddChild() Subscriber

	// OnUnsubscribe will add the given callback function to the subscriber.
	// The callback will be called when either the Unsubscribe of the parent
	// or of the subscriber itself is called.
	OnUnsubscribe(func())

	// OnWait will register a callback to  call when subscription Wait is called.
	OnWait(func())
}

// New will create and return a new Subscriber.
func New() Subscriber {
	return &subscription{}
}
