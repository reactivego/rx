// PublishLast returns a Multicaster that shares a single subscription to the
// underlying Observable containing only the last value emitted before it
// completes. When the underlying Obervable terminates with an error, then
// subscribed observers will receive only that error (and no value). After all
// observers have unsubscribed due to an error, the Multicaster does an internal
// reset just before the next observer subscribes.
package PublishLast
