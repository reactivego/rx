/*
	AsObservableBar

AsObservableBar when called on an Observable source will type assert the
interface{} items of the source to bar items. If the assertion is invalid, this
will emit a type assertion error ErrTypecastToBar at run-time. When
AsObservableBar is called on an ObservableFoo, then a type conversion is
generated. This has to be possible, because if it is not, your generated code
will simply not compile. A special case of this is when using AsObservable, as
this will convert to interface{}, so that will always compile and work.
*/
package AsObservable

import _ "github.com/reactivego/rx"
