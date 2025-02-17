# subscriber

    import "github.com/reactivego/rx/subscriber"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/rx/subscriber.svg)](https://pkg.go.dev/github.com/reactivego/rx/subscriber#section-documentation)

Package `subscriber` provides a subscription tree implementation.

The `New` function creates a subscription and returns a `Subscriber` interface.
It allows a publisher and subscribing client to communicate about a
subscription. The publisher can see whether the client is still subscribed
by polling the `Closed` method. The client can control the subscription
by either calling `Unsubscribe` itself or when it has defered subscription
management to other code call `Wait` to wait for other code to cancel the
subscription.

Once the client has called `Unsubscribe` on the subscription, that subscription
will then call `Unsubscribe` on all of its children created through the `Add`
method. After that it will become lifeless, meaning no events will ever be
triggered by it anymore.

In case the publisher decides it wants to stop publishing, it MUST NEVER
call `Unsubscribe` itself on a subscriber, but should indicate the desire
through other (out-of-band) means to the subscribing client who must then
call `Unsubscribe` itself.

The client that receives a `Subcriber` interface can call `Wait` to wait for
the subscription to be canceled. That may be done by the code that e.g.
created the subscription in responese to some external event.

The client can also keep a reference to the subscription and call
`Unsubscribe` itself to indicate to the publisher it is no longer interested
in receiving data associated with the subscription.

The implementation is designed to be used from concurrently running
goroutines. It uses WaitGroups, Mutexes and atomic reference counting.
