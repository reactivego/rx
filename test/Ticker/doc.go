/*
Ticker creates an ObservableTime that emits a sequence of timestamps after
an initialDelay has passed. Subsequent timestamps are emitted using a
schedule of intervals passed in. If only the initialDelay is given, Ticker
will emit only once.

Ticker is like Timer, but will emit the current time instead of a sequence
of ints.

	Ticker	http://reactivex.io/documentation/operators/timer.html
*/
package Ticker
