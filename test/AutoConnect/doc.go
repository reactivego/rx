/*
AutoConnect is a declarative way of calling Connect on a Connectable.

It is an operator on Connectable that returns a normal Observable. The operator
takes an int argument specifying the number of subscribers that need to
subscribe to it before it will call the Connect method on the Connectable. When
AutoConnect calls Connect, this will subscribe the Connectable to its source
Observable.

Unlike RefCount, AutoConnect never automatically Connects again and also will
never call Unsubscribe on the subscription returned by the call to Connect. So
even when the number of subscribers to AutoConnect drops to zero, the
subscription from the Connectable to its source will remain alive.

	AutoConnect
*/
package AutoConnect

import _ "github.com/reactivego/rx"
