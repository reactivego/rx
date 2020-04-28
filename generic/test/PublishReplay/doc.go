/*
	Replay	http://reactivex.io/documentation/operators/replay.html

PublishReplay uses Multicast to control the subscription of a ReplaySubject to a
source observable and turns the subject into a connectable observable.
A ReplaySubject emits to any observer all of the items that were emitted by the
source observable, regardless of when the observer subscribes.

If the source completed and as a result the internal ReplaySubject
terminated, then calling Connect again will replace the old ReplaySubject
with a newly created one.
*/
package PublishReplay

import _ "github.com/reactivego/rx"
