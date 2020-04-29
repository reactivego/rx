/*
ToChan returns a channel that emits interface{} values. If the source
observable does not emit values but emits an error or complete, then the
returned channel will emit any error and then close without emitting any
values.

This method subscribes to the observable on the Goroutine scheduler because
it needs the concurrency so the returned channel can be used by used
by the calling code directly. To cancel ToChan you will need to supply a
subscriber that you hold on to.

	ToChan		http://reactivex.io/documentation/operators/to.html
*/
package ToChan
