/*
Debounce only emits the last item of a burst from an Observable if a particular
timespan has passed without it emitting another item.

Debounce uses a bare goroutine internally in its implementation.

	Debounce	http://reactivex.io/documentation/operators/debounce.html
*/
package Debounce

