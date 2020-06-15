/*
AutoConnect makes a Multicaster behave like an ordinary Observable that
automatically connects the multicaster to its source when the specified number
of observers have subscribed to it.

If the count is less than 1 it will emit InvalidCount when
subscribed to.

After connecting, when the number of subscribed observers eventually drops to
0, AutoConnect will cancel the source connection if it hasn't terminated yet.

When subsequently the next observer subscribes, AutoConnect will connect to
the source only when it was previously canceled or because the source
terminated with an error. So it will not reconnect when the source completed
succesfully. This specific behavior allows for implementing a caching
observable that can be retried until it succeeds.

Another thing to notice is that AutoConnect will disconnect an active
connection when the number of observers drops to zero. The reason for this is
that not doing so would leak a task and leave it hanging in the scheduler.
*/
package AutoConnect

