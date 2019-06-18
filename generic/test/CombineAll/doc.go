/*
Operator CombineAll documentation and tests.

	CombineAll

CombineAll flattens a higher order observable by collection all emitted
observables until the source completes. Then CombineAll subscribes to all the
collected observables and whenever they all have emitted an item it will emit
a slice of all emitted itenm. When subsequently, one of the observables emits
a new item, CombineAll will emit an updated slice with one of the items
replaced with the newly emitted item.
*/
package CombineAll

import _ "github.com/reactivego/rx"
