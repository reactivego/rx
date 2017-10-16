/*
Operator AutoConnect documentation and tests.

	AutoConnect

AutoConnect is an operator on Connectable that returns a normal Observable. The
operator takes an int argument specifying the number of subscribers that need to
subscribe to it before it will call the Connect method on the Connectable.
Calling Connect will subscribe the Connectable to its source Observable.
Unlike RefCount, AutoConnect never automatically unusbscribes.
You can pass a multiple function options; OnSubscribe(), OnUnsubscribe() and
SubscribeOn() to the AutoConnect. The will be passed into the Connect method
when it is called.
*/
package AutoConnect

import _ "github.com/reactivego/rx"
