/*
	ToChan		http://reactivex.io/documentation/operators/to.html

ToChan returns a channel that emits values. For Observables of type
interface{}, any error is emitted through the channel itself.

Because the channel is fed by subscribing to the observable, ToChan would
block when subscribed on the Trampoline scheduler which is initially synchronous.
That's why the subscribing is done on the Goroutine scheduler. This also means
that the code reading the channel can do so while running on the main goroutine.

To cancel the subscription created internally by ToChan you will need access to
the subscription used internally by ToChan. To get at this subscription, pass
the result of a call to option OnSubscribe(func(Subscription)) as a parameter
to ToChan. On subcription the callback will be called with the subscription
that was created.
*/
package ToChan

import _ "github.com/reactivego/rx"
