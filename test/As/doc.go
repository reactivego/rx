/*
Operator As documentation and tests.

	AsAny
	AsBar


AsAny turns a strongly typed ObservableFoo into an any type Observable of
interface{}.

AsBar turns an Observable or ObservableFoo into an ObservableBar. For Observable
a type assertion is used. If the type assertion fails, then the ObservableBar
will emit a type conversion error ErrTypecastToBar at run-time. For
ObservableFoo a type conversion is generated that needs to be possible, because
if it is not, your generated code will not compile.
*/
package As

import _ "github.com/reactivego/rx"

//jig:file {{.package}}.go
