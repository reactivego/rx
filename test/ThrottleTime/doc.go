/*
ThrottleTime emits when the source emits and then starts a timer during which
all emissions from the source are ignored. After the timer expires, ThrottleTime
will again emit the next item the source emits, and so on.
*/
package ThrottleTime
