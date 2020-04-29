/*
SwitchMap transforms the items emitted by an Observable by applying a function
to each item an returning an Observable.

In doing so, it behaves much like MergeMap (previously FlatMap), except that
whenever a new Observable is emitted SwitchMap will unsubscribe from the
previous Observable and begin emitting items from the newly emitted one.

	SwitchMap
		http://reactivex.io/documentation/operators/flatmap.html
		https://www.learnrxjs.io/operators/transformation/switchmap.html
*/
package SwitchMap
