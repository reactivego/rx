/*
PublishReplay returns a Multicaster for a ReplaySubject to an underlying
Observable and turns the subject into a connectable observable. A
ReplaySubject emits to any observer all of the items that were emitted by
the source observable, regardless of when the observer subscribes. When the
underlying Obervable terminates with an error, then subscribed observers
will receive that error. After all observers have unsubscribed due to an
error, the Multicaster does an internal reset just before the next observer
subscribes.
*/
package PublishReplay
