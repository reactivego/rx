/*
	Publish	http://reactivex.io/documentation/operators/publish.html

Publish uses Multicast to control the subscription of a Subject to a
source observable and turns the subject it into a connnectable observable.
A Subject emits to an observer only those items that are emitted by
the source Observable subsequent to the time of the subscription.

If the source completed and as a result the internal Subject terminated, then
calling Connect again will replace the old Subject with a newly created one.
So this Publish operator is re-connectable, unlike the RxJS 5 behavior that
isn't. To simulate the RxJS 5 behavior use Publish().AutoConnect(1) this will
connect on the first subscription but will never re-connect.
*/
package Publish

import _ "github.com/reactivego/rx"
