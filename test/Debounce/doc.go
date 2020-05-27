/*
Debounce only emits the last item of a burst from an Observable if a particular
timespan has passed without it emitting another item.

Debounce uses a bare goroutine internally in its implementation.

	Debounce    http://reactivex.io/documentation/operators/debounce.html
	            https://rxjs.dev/api/operators/debounceTime
	            https://www.learnrxjs.io/learn-rxjs/operators/filtering/debouncetime
*/
package Debounce

