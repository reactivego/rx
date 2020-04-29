/*
CombineAll flattens a higher order observable (e.g. an ObservableObservableInt)
by collection all emitted observables (e.g. ObservableInt items) until the
source completes.

Then CombineAll subscribes to all the collected observables and whenever they
all have emitted an item it will emit a slice of all emitted items. When
subsequently, one of the observables emits a new item, CombineAll will emit an
updated slice with one of the items replaced with the newly emitted item.

	CombineAll
*/
package CombineAll
